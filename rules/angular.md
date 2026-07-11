# Angular Rules

Read this file when writing any Angular component, service, directive, pipe, or template. It contains only Angular-specific rules: the general files still apply in full — architecture.md (boundaries, folder structure), reliability.md (errors, timeouts, retries), security.md, testing.md, maintainability.md. Where this file tightens a general rule for Angular code, the tighter rule wins within Angular scope (GEN-03).

## Framework Standards

- **NG-01.** Applications MUST use an Angular major that is actively supported (current or LTS per the official support policy); unsupported versions MUST NOT run production systems, and applications on LTS MUST have an upgrade plan before support ends.
- **NG-02.** Greenfield applications SHOULD start on the latest stable Angular major.
- **NG-03.** Deprecated Angular APIs MUST NOT be introduced in new code. This includes the structural directives `*ngIf`, `*ngFor`, and `*ngSwitch` (deprecated since v20) — use built-in control flow (`@if`, `@for`, `@switch`).
- **NG-04.** Supported-but-legacy APIs — constructor injection, `@Input`, `@Output`, `@HostBinding`/`@HostListener` — SHOULD NOT be introduced in new code unless compatibility, inheritance, or third-party constraints require them; note the constraint in the PR.
- **NG-05.** `standalone: true` MUST NOT be written where the project's Angular version treats standalone as the default. Gate: CI (lint).
- **NG-06.** Built-in Angular features SHOULD be preferred over third-party libraries when the built-in meets the need; experimental APIs MUST NOT be used in critical code without an exception per GEN-05.

## Project Configuration

- **NG-07.** Strict TypeScript mode and `strictTemplates` MUST be enabled; `strictInjectionParameters` SHOULD be enabled, and `strictStandalone` SHOULD be enabled once the standalone migration is complete. Gate: CI.
- **NG-08.** `any` MUST NOT be used unless the type is genuinely unknown and the usage is justified at the use site; `unknown` SHOULD be preferred at trust boundaries. Type assertions MUST be justified and kept close to validation. Types SHOULD NOT be declared where the compiler infers them correctly.
- **NG-09.** Angular extended diagnostics SHOULD be enabled and reviewed in CI.
- **NG-10.** Production builds MUST use Angular production optimizations, and bundle budgets MUST be configured in `angular.json` (initial bundle and component styles) and MUST fail the CI production build when exceeded. Gate: CI.
- **NG-11.** Production source maps MUST be hidden (not publicly served) and uploaded only to the error-tracking system, unless a documented org policy states otherwise.
- **NG-12.** Environment-specific configuration MUST be injected through approved configuration mechanisms, not hardcoded (MNT-35).

## Architecture

Folder structure and boundary rules are owned by architecture.md (ARCH-25/26 examples apply). For Angular code the feature-first organization is MUST, not SHOULD:

- **NG-13.** Angular code MUST be organized by feature or bounded context (`/features/orders/`, `/features/auth/`), not primarily by file type (`/components/`, `/services/`).
- **NG-14.** Each feature MUST separate concerns: UI (components, templates, presentation-only pipes/directives), state (signals, stores, state services), business logic (domain services, use cases, validators), and infrastructure (HTTP, browser APIs, storage, third-party integrations).
- **NG-15.** Business logic MUST NOT live in components or templates, and infrastructure code MUST NOT leak into business logic or components.
- **NG-16.** Cross-cutting concerns (authentication, logging, error handling, tracing, analytics, configuration, authorization) SHOULD be centralized; shared UI and utility code follows ARCH-07.

## Components

- **NG-17.** Each component MUST have one responsibility (MNT-01) and MUST coordinate presentation, user interaction, and component-local state only; non-presentation logic MUST move to a service, use case, validator, or pure function.
- **NG-18.** Components SHOULD NOT exceed 300 source lines (`.ts` file including inline template; project-configurable). Keeping a larger component requires a one-line justification in the PR accepted by the reviewer.
- **NG-19.** Direct DOM manipulation MUST NOT be used unless Angular templates, bindings, directives, CDK, or renderer-safe APIs cannot meet the requirement; components integrating non-Angular DOM libraries MUST isolate that integration and clean up resources on destroy.
- **NG-20.** The `host` property in the decorator SHOULD be used instead of `@HostBinding`/`@HostListener` in new code; the decorators MAY be used in shared base classes where `host` metadata is unavailable.
- **NG-21.** Components SHOULD use self-closing tags where no projected content exists, and SHOULD use inline templates and styles when the template is under 15 lines. External templates and styles MUST use paths relative to the component file.
- **NG-22.** Deeply nested component trees SHOULD be flattened through composition, projection, routing, or feature boundaries when a presentational chain passes inputs through more than two intermediate layers.

## Smart and Presentational Components

- **NG-23.** Smart components MAY own data fetching, routing coordination, orchestration, and feature state. Presentational components MUST receive data through inputs, communicate through outputs, and MUST NOT contain business logic, HTTP calls, routing decisions, or global state mutation.
- **NG-24.** Presentational components SHOULD be pure: the same inputs produce the same rendered output.
- **NG-25.** Shared presentational components MUST be accessible, themeable, and documented when reused across features.

## Templates

- **NG-26.** New templates MUST use built-in control flow (`@if`, `@for`, `@switch`) — see NG-03 for the deprecated directives.

```html
@if (user) {
  <div>{{ user.name }}</div>
} @else {
  <div>No user</div>
}

@for (item of items; track item.id) {
  <div>{{ item.name }}</div>
} @empty {
  <div>No items</div>
}
```

- **NG-27.** `@for` track expressions MUST uniquely and stably identify items (the compiler already requires `track` to be present).
- **NG-28.** Functions called in templates MUST perform no I/O, no mutation, and no per-call object or array allocation; heavier or impure work MUST move to `computed()`, a service, a pipe, or a precomputed view model. Template expressions MUST NOT produce side effects.
- **NG-29.** Templates SHOULD stay declarative: complex branching and transformation belongs in component code, `computed()`, or services. Signals read in templates MUST be invoked.
- **NG-30.** `ngClass` and `ngStyle` MUST NOT be used; use `class` and `style` property bindings.

Non-normative note: template expressions can only reference component members and template context (enforced by `strictTemplates`); expose values like `Math` or `Date` computations via component properties or `computed()`.

## Signals and State

Signals are the default primitive for synchronous reactive state.

```ts
count = signal(0);
double = computed(() => this.count() * 2);
```

- **NG-31.** `signal()` MUST be used for local mutable state that affects rendering; `computed()` MUST be used for values derived from signals; derived values MUST NOT be copied into separate writable signals. `linkedSignal()` MAY be used for writable state intrinsically linked to another signal.
- **NG-32.** Writable signals MUST NOT be mutated from outside the owning component or service except through an explicit public method, and services exposing state SHOULD expose read-only signals.
- **NG-33.** Shared state MUST have one owner (REL-09); global state SHOULD live in dedicated state services or stores.
- **NG-34.** Signal values containing objects or arrays SHOULD be updated immutably with `set` or `update`; the removed `mutate` API MUST NOT appear in code targeting modern Angular.

## Signals vs RxJS

- **NG-35.** Signals MUST be used for synchronous, template-facing state and values derived from other reactive state. RxJS SHOULD be used when the source is naturally an Observable (HTTP, WebSocket, router events, form streams) or when stream coordination is needed (cancellation, retry, debounce, `switchMap`, `combineLatest`).
- **NG-36.** Conversion between signals and Observables SHOULD go through `toSignal()`/`toObservable()` from `@angular/core/rxjs-interop`, avoiding manual subscription management.
- **NG-37.** Observable subscriptions MUST be cleaned up with `takeUntilDestroyed()`, `AsyncPipe`, `toSignal()`, or an equivalent lifecycle-safe mechanism; long-lived subscriptions in services MUST have explicit ownership and cleanup.

## Effects

- **NG-38.** `effect()` MUST be limited to controlled side effects at the edge of the system; `computed()` SHOULD be tried first for derivation and `linkedSignal()` for dependent writable state.

`effect()` MUST NOT be used for: API calls as a substitute for explicit user, route, or lifecycle actions; form patching loops; syncing unrelated signal state; replacing service methods; hiding unclear ownership.

`effect()` MAY be used for: logging and analytics; synchronizing state with browser storage; synchronizing with a non-Angular library; custom DOM or canvas integration that cannot be expressed declaratively.

- **NG-39.** Effects that allocate resources MUST register cleanup. All effects MUST use `untracked()` for incidental reads and MUST NOT create circular updates.

Bad:

```ts
effect(() => {
  this.api.getData(this.id()).subscribe();
});
```

Good:

```ts
loadData() {
  return this.api.getData(this.id());
}
```

## Dependency Injection

- **NG-40.** `inject()` SHOULD be the default injection mechanism in new code (constructor injection per NG-04).

```ts
private readonly api = inject(ApiService);
```

- **NG-41.** Injected services MUST be used; injection tokens SHOULD be used for configuration values and abstract dependencies; providers SHOULD be scoped as narrowly as practical for their lifetime.
- **NG-42.** Root-provided services MUST be safe to keep for the application's lifetime, and feature-level providers MUST NOT create duplicate state unless isolation is intentional and documented.

## Inputs and Outputs

- **NG-43.** New code SHOULD use signal-based `input()`/`output()`; required inputs MUST use `input.required<T>()`.

```ts
user = input.required<User>();
saved = output<User>();
```

- **NG-44.** Input transforms MUST be pure and MUST NOT change the input's semantic meaning; components MUST NOT mutate input-owned data.
- **NG-45.** Output names MUST describe events, not commands, and SHOULD be past tense (`saved`, `deleted`, `selected`). Two-way binding SHOULD be used sparingly and only with clear ownership.

## Services

- **NG-46.** Services MUST be separated by concern — API services (HTTP and response mapping), business logic services, state services, integration services — and each MUST have one responsibility; god services MUST be split.
- **NG-47.** Services MUST NOT manipulate the UI directly, and services exposing mutable state MUST provide explicit mutation methods.

## HTTP and API Layer

Timeouts, bounded retries, and backoff are owned by reliability.md §External Services (REL-14 to REL-16) and apply in full to Angular HTTP code.

- **NG-48.** HTTP calls MUST be made from dedicated API or infrastructure services, never directly from components.
- **NG-49.** `provideHttpClient()` SHOULD configure HTTP, with functional interceptors preferred for predictable ordering; interceptors SHOULD carry cross-cutting concerns (auth headers, tracing, deadlines, global error mapping, safe retry policy).
- **NG-50.** Every HTTP error MUST have an explicit handling strategy (callsite, service boundary, resource state, or interceptor) mapped to a defined outcome per REL-01; unhandled HTTP errors MUST NOT silently fail.
- **NG-51.** API responses MUST be typed (`any` prohibited per NG-08), and response DTOs SHOULD be mapped to domain or view models instead of leaking backend shapes through the UI.
- **NG-52.** Long-running API work SHOULD expose loading, success, empty, partial, and error states (REL-11).
- **NG-53.** Applications MUST register a global `ErrorHandler` (plus browser global error listeners where supported) that reports unhandled errors to the telemetry sink and leaves the user in a safe, recoverable state.

## Forms

- **NG-54.** Reactive Forms MUST be used when a form has any of: cross-field or async validation, dynamically added or removed controls, conditional structure, or a multi-step flow. Template-driven forms MAY be used below all of those triggers.
- **NG-55.** New forms MUST be typed (`FormGroup<T>`, `FormControl<T>`, `FormArray<T>`); `UntypedFormGroup`/`UntypedFormControl` MUST NOT be introduced outside migration code, justified with a comment. Template-driven and reactive forms MUST NOT be mixed in one feature.
- **NG-56.** Validation MUST occur on both client and server; client-side validation is UX, not a boundary (REL-07).
- **NG-57.** Forms MUST NOT be patched from service responses without explicit user, route, or lifecycle intent, and submission MUST prevent duplicates where duplicate side effects are unsafe (REL-12).
- **NG-58.** Form errors MUST be accessible and associated with their controls; sensitive fields MUST NOT be logged or stored insecurely; unsaved-change workflows SHOULD protect users from data loss.

## Routing

- **NG-59.** Feature routes MUST be lazy-loaded; eager loading requires a documented justification in the PR. (Single rule; no duplicate in Performance.)
- **NG-60.** Feature route configuration SHOULD live in a dedicated routes file per feature, route data SHOULD be typed, and route params and query params MUST be validated before use.
- **NG-61.** Client-side route guards MUST NOT be treated as a security boundary; authorization is enforced per SEC-10. Guards MAY handle authentication UX, feature flags, unsaved changes, and conditional navigation.
- **NG-62.** Guards SHOULD return a `UrlTree` or redirect command instead of returning `false` and navigating imperatively, MUST NOT stall indefinitely, and guards awaiting slow checks SHOULD redirect to a loading route rather than blocking navigation.
- **NG-63.** Route-level providers SHOULD be intentional and reviewed for duplicate state (NG-42).

## Performance and Change Detection

- **NG-64.** `ChangeDetectionStrategy.OnPush` SHOULD be used for application components; eager/default change detection requires a documented reason.
- **NG-65.** Components MUST NOT rely on Zone.js to trigger change detection: state changes MUST flow through signals, `AsyncPipe`, or `markForCheck()`. Projects targeting zoneless MUST run their test suite with `provideZonelessChangeDetection()`.
- **NG-66.** Lists that are unbounded or can exceed the project's rendering threshold (default 100 DOM rows, project-configurable) MUST use pagination, virtual scrolling, or incremental rendering.
- **NG-67.** `toSignal()` SHOULD be used to consume Observables in templates where signal-based state improves clarity; `AsyncPipe` remains acceptable for simple consumption.
- **NG-68.** Expensive components SHOULD be deferred (`@defer`), lazy-loaded, virtualized, or split.
- **NG-69.** A third-party library adding more than the project's budget to the initial bundle (default 50 KB gzipped) MUST be justified in the PR, and SHOULD be lazy-loaded unless it is needed on initial render; NG-10's budgets enforce the aggregate.
- **NG-70.** `NgOptimizedImage` MUST be used for `<img>` elements loading URL-referenced raster images. Base64/data URIs, CSS background images, and inline SVG are out of scope; those images MUST still declare explicit dimensions to prevent layout shift.

## SSR, Hydration, and Rendering

- **NG-71.** Public, SEO-sensitive, or performance-sensitive applications SHOULD evaluate SSR, prerendering, or hybrid rendering; the decision is recorded per GEN-11.
- **NG-72.** SSR applications MUST keep server-rendered and client-rendered DOM consistent, and browser-only APIs MUST NOT run during server rendering — use render-safe lifecycle APIs or platform-specific providers, not `isPlatformBrowser` checks in templates (hydration mismatch risk).
- **NG-73.** Components that opt out of hydration MUST document why.
- **NG-74.** Server-side providers MUST NOT leak request-specific data across users, and HTTP transfer-cache behaviour MUST be reviewed for sensitive authenticated requests.

## Security (Angular-specific; security.md applies in full)

- **NG-75.** `bypassSecurityTrust*` MUST NOT be used unless the value source is verified safe, the justification is documented at the call site, and the usage passes security review when the value can be influenced by users or external systems. User-controlled values MUST NOT be rendered as HTML without sanitization (SEC-19).
- **NG-76.** Direct DOM APIs, `ElementRef.nativeElement`, and third-party DOM manipulation MUST be treated as security-sensitive and reviewed as such.
- **NG-77.** Content Security Policy and Trusted Types SHOULD be enabled for public or high-risk applications (SEC-29).
- **NG-78.** Angular's XSRF helper SHOULD be enabled for cookie-authenticated applications; CSRF protection MUST be enforced by the backend regardless (SEC-27).
- **NG-79.** Tokens and secrets MUST NOT be stored in insecure client-side storage without an exception per GEN-05, and Angular security advisories MUST be reviewed during framework upgrades.

## Accessibility

- **NG-80.** Applications MUST meet WCAG 2.1 AA (or a stricter documented org baseline). User-facing features MUST produce zero axe-core critical/serious violations, asserted in component or e2e tests in CI where the project has axe tooling; where it has none, the author MUST flag the gap in the PR rather than skip the check. Gate: CI/Evidence.
- **NG-81.** Semantic HTML and native controls MUST be preferred over custom interactive elements; ARIA MUST NOT compensate for incorrect HTML where a native element exists.
- **NG-82.** Custom interactive components MUST support keyboard navigation, focus states, accessible names, and appropriate ARIA semantics; dialogs, menus, and overlays MUST handle focus trapping, escape behaviour, and focus restoration.
- **NG-83.** Form controls MUST have accessible labels and error messages, and status states (error, loading, empty, success) SHOULD be exposed to assistive technologies; dynamically inserted important content SHOULD use live regions.
- **NG-84.** Route changes SHOULD manage focus to the new page content, and active navigation links SHOULD expose `aria-current`.

## Styling

- **NG-85.** Styles SHOULD be scoped via Angular's default view encapsulation; global styles MUST be limited to resets, themes, typography, design tokens, and application-level layout.
- **NG-86.** Design tokens SHOULD be used for shared color, spacing, typography, elevation, and motion; style duplication SHOULD be extracted to tokens, mixins, or shared styles.
- **NG-87.** Inline styles MUST NOT be used unless the value is dynamically bound from component state, and styling MUST NOT hide focus indicators without an accessible replacement.
- **NG-88.** Component styles SHOULD support responsive layouts and user preferences (reduced motion, text scaling).

## Internationalization

Applies when the application supports or plans to support multiple locales; otherwise inert.

- **NG-89.** User-facing strings MUST go through the i18n mechanism, never hardcoded; translation keys MUST be stable and reviewed like contracts.
- **NG-90.** Dates, numbers, and currency MUST use locale-aware formatting APIs, and layouts SHOULD tolerate text expansion and RTL where required.

## Testing (Angular-specific; testing.md applies in full)

Generic rules — regression tests, determinism, flakiness, naming, production data — are owned by testing.md; do not restate them here.

- **NG-91.** New projects MUST use a supported, non-deprecated test runner (e.g. the Angular CLI's supported runner or Jest/Vitest per org choice); Karma MUST NOT be adopted for new projects.
- **NG-92.** `TestBed` MUST be used for component tests involving rendering, inputs/outputs, DI, content projection, routing, or forms; component harnesses SHOULD be used for Material/CDK components.
- **NG-93.** Presentational components SHOULD be tested by setting inputs and asserting rendered output and emitted events; smart components SHOULD be tested by mocking injected services at the boundary.
- **NG-94.** Services without Angular-specific dependencies SHOULD be tested without `TestBed`; services with injected dependencies MUST use `TestBed` or `inject()` in an injection context. HTTP calls MUST be tested with `HttpTestingController` or equivalent.
- **NG-95.** Guards MUST be tested for allowed, denied, and redirected outcomes, and route parameter handling SHOULD be tested for valid, missing, and invalid values.
- **NG-96.** Signals MUST be tested through public logic and observable outcomes; `computed()` through its dependencies; `effect()` SHOULD be tested via the external side effect, not directly.
- **NG-97.** Accessibility-critical components SHOULD include accessibility assertions (NG-80).

## Definition of Done (Angular delta)

All items in engineering-rules.md §Definition of Done apply. Additionally, an Angular change is not complete unless: modern APIs are used per NG-03/NG-04; state uses signals where NG-31 applies; UI states are handled per NG-52; routes are lazy-loaded per NG-59; accessibility (NG-80) and SSR impact (NG-71 to NG-74, when applicable) were addressed; and bundle budgets (NG-10) pass.
