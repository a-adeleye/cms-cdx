# Angular Rules

Standards for building scalable, maintainable, accessible, secure, and production-ready Angular applications. Rule levels follow the engineering rules: MUST is required, SHOULD is required by default with justified exceptions, and MAY is optional guidance.

## Framework Standards

- Applications MUST use an Angular version that is actively supported or approved for production use.
- Greenfield applications SHOULD use the latest stable Angular major approved by the organization.
- Applications on LTS or older supported versions MUST have an upgrade plan before support ends.
- Unsupported Angular versions MUST NOT be used for production systems.
- Modern Angular APIs SHOULD be used for new code.
- Deprecated Angular APIs MUST NOT be introduced in new code.
- Legacy or non-preferred APIs such as `*ngIf`, `*ngFor`, constructor injection, `@Input`, and `@Output` SHOULD NOT be introduced in new code unless compatibility, migration, or third-party constraints require them.
- `standalone: true` MUST NOT be written when the project's Angular version treats standalone declarations as the default.
- Built-in Angular features SHOULD be preferred over third-party libraries when the built-in feature meets the need.
- Experimental Angular APIs MUST NOT be used in production-critical paths unless the risk is explicitly accepted.

## Project Configuration

- Strict TypeScript mode MUST be enabled.
- Types SHOULD NOT be explicitly declared when they can be correctly inferred by the TypeScript compiler.
- `any` MUST NOT be used unless the type is genuinely unknown and the usage is explicitly justified.
- `unknown` SHOULD be preferred over `any` at trust boundaries or when type is uncertain.
- Angular template strictness MUST be enabled with `strictTemplates`.
- Angular DI strictness SHOULD be enabled with `strictInjectionParameters`.
- Standalone strictness SHOULD be enabled with `strictStandalone` when the project has completed its standalone migration.
- Angular extended diagnostics SHOULD be enabled and reviewed in CI.
- Formatting, linting, type checking, unit tests, and production builds MUST run in CI.
- Production builds MUST use Angular production optimizations.
- Source maps for production MUST follow the organization's security and observability policy.
- Environment-specific configuration MUST be injected through approved configuration mechanisms, not hardcoded in application code.

## Architecture

- Code MUST be organized by feature or bounded context, not primarily by file type.

Preferred:

```text
/features/orders/
/features/auth/
/features/billing/
```

Avoid as the primary structure:

```text
/components/
/services/
/models/
```

- Each feature MUST separate concerns:
  - UI: components, templates, styles, presentation-only pipes and directives.
  - State: signals, stores, and state services.
  - Business logic: domain services, use cases, validators, and orchestration.
  - Infrastructure: HTTP services, browser APIs, storage, analytics, and third-party integrations.
- Business logic MUST NOT live in components or templates.
- Infrastructure code MUST NOT leak into business logic or components.
- Components MUST depend on feature-level abstractions rather than directly coupling to unrelated features.
- Cross-cutting concerns such as authentication, logging, error handling, tracing, analytics, configuration, and authorization SHOULD be centralized.
- Shared UI and utility code MUST be generic, stable, and free of hidden feature-specific behaviour.

## Components

- Each component MUST have one clear responsibility.
- Components exceeding 300 lines SHOULD be reviewed for splitting.
- Components MUST coordinate presentation, user interaction, and component-local state only.
- Logic that does not relate to presentation MUST be moved to a service, use case, validator, or pure function.
- Direct DOM manipulation MUST NOT be used unless Angular templates, bindings, directives, CDK, or renderer-safe APIs cannot meet the requirement.
- Components that integrate non-Angular DOM libraries MUST isolate that integration and clean up resources on destroy.
- Deeply nested component trees SHOULD be flattened by improving composition, projection, routing, or feature boundaries.
- Components SHOULD use self-closing tags where no projected content exists.
- Components SHOULD use inline templates and styles for components with small templates (e.g., < 15 lines).
- External templates and styles MUST be referenced using paths relative to the component file.
- `host` property MUST be used in the `@Component` decorator instead of `@HostBinding` or `@HostListener`.

## Smart And Presentational Components

- Smart components MAY own data fetching, routing coordination, orchestration, and feature state.
- Presentational components MUST receive data through inputs and communicate changes through outputs.
- Presentational components MUST NOT contain business logic, direct HTTP calls, routing decisions, or global state mutation.
- Presentational components SHOULD be pure: the same inputs should produce the same rendered output.
- Shared presentational components MUST be accessible, themeable, and documented when reused across features.

## Template Syntax

- Modern control flow MUST be used for new templates.
- Legacy structural directives such as `*ngIf`, `*ngFor`, and `*ngSwitch` MUST NOT be introduced in new code unless migration constraints require them.

Modern:

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

Legacy:

```html
<div *ngIf="user"></div>
<div *ngFor="let item of items"></div>
```

- `track` MUST always be used with `@for`.
- `track` expressions MUST uniquely and stably identify items.
- Functions called in templates MUST be pure and cheap.
- Heavy or impure template functions MUST be moved to a `computed()`, service, pipe, or precomputed view model.
- Inline object or array creation in templates SHOULD be avoided because it can cause unnecessary work.
- Template expressions MUST NOT produce side effects.
- Templates SHOULD stay declarative. Complex branching and transformation logic belongs in component code, `computed()`, or services.
- Signals read in templates MUST be invoked.
- `ngClass` and `ngStyle` MUST NOT be used. Use `class` and `style` property bindings instead.
- Templates MUST NOT assume global JavaScript objects (e.g., `Date`, `Math`, `window`) are available; provide them through component properties/methods.

## Signals And State

- Signals are the default primitive for synchronous local and shared reactive state.

```ts
count = signal(0);
double = computed(() => this.count() * 2);
```

- `signal()` MUST be used for local mutable state that affects rendering.
- `computed()` MUST be used for values derived from signals.
- `linkedSignal()` MAY be used for writable state that is intrinsically linked to another signal.
- Derived values SHOULD NOT be copied into separate writable signals.
- Writable signals MUST NOT be mutated from outside the component or service that owns them unless through an explicit public method.
- Shared state MUST have a clear owner.
- Global state SHOULD be owned by dedicated state services or stores, not scattered across components.
- Signal values containing objects or arrays SHOULD be updated immutably.
- Signal `mutate` MUST NOT be used (deprecated/removed in modern Angular); use `set` or `update` instead.
- Read-only signals SHOULD be exposed from services when consumers must observe but not mutate state.

## Signals Vs RxJS

Use signals when:

- The value is synchronous.
- The state is local to a component or owned by a service.
- A value is derived from other reactive state.
- The value is read directly by templates.

Use RxJS when:

- The source is naturally an Observable, such as HTTP, WebSocket, router events, form streams, or external event streams.
- You need cancellation, retry, debouncing, throttling, buffering, or complex async coordination.
- You need stream operators such as `switchMap`, `exhaustMap`, `debounceTime`, or `combineLatest`.

Rules:

- Mixing signals and RxJS SHOULD be done through `toSignal()` and `toObservable()` from `@angular/core/rxjs-interop`.
- Manual subscription management when converting between signals and Observables SHOULD be avoided.
- Observable subscriptions MUST be cleaned up with `takeUntilDestroyed()`, `AsyncPipe`, `toSignal()`, or an equivalent lifecycle-safe mechanism.
- Long-lived subscriptions in services MUST have explicit ownership and cleanup expectations.

## Effects

- `effect()` SHOULD be the last reactive API reached for.
- Prefer `computed()` for derived values.
- Prefer `linkedSignal()` when writable state depends on another signal.
- `effect()` MUST be limited to controlled side effects at the edge of the system.

`effect()` MUST NOT be used for:

- API calls as a substitute for explicit user, route, or lifecycle actions.
- Form patching loops.
- Syncing unrelated signal state.
- Replacing service methods.
- Hiding unclear ownership or bad architecture.

`effect()` MAY be used for:

- Logging and analytics side effects.
- Synchronizing Angular state with browser storage.
- Synchronizing Angular state with a third-party non-Angular library.
- Custom DOM or canvas integration that cannot be expressed declaratively.

Rules:

- Effects that allocate resources MUST register cleanup.
- Effects MUST avoid accidental dependency tracking by using `untracked()` where incidental reads are necessary.
- Effects MUST NOT create circular updates.

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

- `inject()` SHOULD be used as the default injection mechanism for new code.

```ts
private readonly api = inject(ApiService);
```

- Constructor injection SHOULD NOT be used in new code unless compatibility, inheritance, testing, decorators, or framework constraints require it.
- Services MUST NOT be injected unless they are used.
- Injection tokens SHOULD be used for configuration values, abstract dependencies, and environment-specific implementations.
- Providers SHOULD be scoped as narrowly as practical.
- Root-provided services MUST be safe to keep for the lifetime of the application.
- Feature-level providers MUST NOT accidentally create duplicate state unless feature isolation is intentional.

## Inputs And Outputs

- Modern signal-based input and output APIs SHOULD be used in new code.

```ts
user = input.required<User>();
saved = output<User>();
```

Legacy decorator APIs SHOULD NOT be introduced in new code unless required for compatibility:

```ts
@Input() user!: User;
@Output() saved = new EventEmitter<User>();
```

Rules:

- Required inputs MUST use `input.required<T>()`.
- Optional inputs SHOULD use `input<T>()` with a default or explicit `undefined` typing.
- Input transforms MUST be pure and must not change the semantic meaning of the input.
- Output names MUST describe events, not commands.
- Output names SHOULD be past tense or event-like: `saved`, `deleted`, `selected`, `closed`.
- Components MUST NOT mutate input-owned data directly.
- Two-way binding between components SHOULD be used sparingly and only when ownership remains clear.

## Services

- Each service MUST have a single responsibility.
- Services MUST be separated by concern:
  - API services: HTTP calls and response mapping.
  - Business logic services: domain rules and orchestration.
  - State services: shared reactive state.
  - Integration services: browser APIs, analytics, storage, and third-party SDKs.
- Services MUST NOT directly manipulate the UI.
- Services MUST NOT become hidden global state containers by accident.
- God services that accumulate unrelated responsibilities MUST be split.
- Services that expose mutable state MUST provide explicit methods for mutation.

## HTTP And API Layer

- HTTP calls MUST be made from dedicated API or infrastructure services, not directly from components.
- `provideHttpClient()` SHOULD be used for HTTP configuration in modern standalone applications.
- Functional HTTP interceptors SHOULD be preferred over DI-based interceptors for predictable ordering.
- HTTP interceptors SHOULD be used for cross-cutting concerns such as authentication headers, request tracing, deadlines, global error mapping, and safe retry policy.
- HTTP errors MUST have an explicit handling strategy at the callsite, service boundary, resource state, or interceptor.
- `catchError` SHOULD be used when an Observable error must be transformed, recovered, logged, mapped, or converted into UI state.
- Unhandled HTTP errors MUST NOT silently fail.
- Retries MUST be bounded and used only for safe or idempotent operations.
- Retries SHOULD use backoff when repeated failure can cause load amplification.
- API responses MUST be typed. `any` MUST NOT be used as a response type.
- API response DTOs SHOULD be mapped to domain or view models instead of leaking backend shapes throughout the UI.
- Request and response payloads SHOULD be validated at trust boundaries when data integrity matters.
- Long-running API work SHOULD expose loading, success, empty, partial, and error states.

## Forms

- Reactive Forms MUST be used for non-trivial, validated, conditional, or dynamic forms.
- Forms MUST be typed using `FormGroup<T>`, `FormControl<T>`, and `FormArray<T>` where practical.
- Template-driven and reactive forms MUST NOT be mixed in the same feature.
- Validation MUST occur on both client and server.
- Client-side validation is UX, not a security boundary.
- Forms MUST NOT be patched from service responses without explicit user, route, or lifecycle intent.
- Form submission MUST prevent duplicate submissions when duplicate side effects are unsafe.
- Form errors MUST be accessible and associated with the relevant controls.
- Sensitive form fields MUST NOT be logged or stored insecurely.
- Unsaved-change workflows SHOULD protect users from accidental data loss.

## Routing

- Routes SHOULD be lazy-loaded by default.
- Feature route configuration SHOULD be kept in a dedicated routes file per feature.
- Route data SHOULD be typed.
- Route params and query params MUST be validated before use.
- Navigation guards MAY be used for authentication, authorization UX, feature flags, unsaved changes, and conditional navigation.
- Client-side route guards MUST NOT be treated as a security boundary.
- Protected backend data and actions MUST enforce authorization on the server or trusted API layer.
- Guards SHOULD return a `UrlTree` or redirect command instead of returning `false` and navigating imperatively.
- Navigation guards MUST NOT stall indefinitely.
- Long-running guards SHOULD redirect to a loading or access-checking route rather than blocking navigation.
- Route-level providers SHOULD be used intentionally and reviewed for duplicate state.

## Performance And Change Detection

- `ChangeDetectionStrategy.OnPush` SHOULD be used for application components unless a documented reason requires eager/default change detection.
- Components SHOULD be compatible with zoneless change detection patterns where practical.
- Heavy or impure computations MUST NOT run during template evaluation.
- `track` MUST always be used with `@for`.
- Feature routes MUST be lazy-loaded unless eager loading is explicitly justified.
- Large lists MUST use pagination, virtual scrolling, or incremental rendering when the list is unbounded or can exceed expected UI limits.
- `toSignal()` SHOULD be used to consume Observables in templates when signal-based state improves clarity.
- The `async` pipe remains acceptable for simple Observable template consumption.
- Expensive components SHOULD be deferred, lazy-loaded, virtualized, or split.
- Production bundles SHOULD be monitored for size regressions.
- Large third-party libraries MUST be justified and lazy-loaded where practical.
- `NgOptimizedImage` MUST be used for all static images. Note: This does not support base64 inline images.

## SSR, Hydration, And Rendering

- Public, SEO-sensitive, or performance-sensitive applications SHOULD evaluate SSR, prerendering, or hybrid rendering.
- SSR-enabled applications MUST keep server-rendered and client-rendered DOM consistent.
- Browser-only APIs MUST NOT run during server rendering.
- Browser-only initialization SHOULD use render-safe lifecycle APIs or platform-specific providers.
- `isPlatformBrowser` checks in templates SHOULD be avoided because they can cause hydration mismatches.
- Components that opt out of hydration MUST document why.
- Server-side providers MUST avoid leaking request-specific data across users.
- HTTP transfer cache behaviour MUST be reviewed for sensitive authenticated requests.

## Security

- Angular's built-in sanitization MUST NOT be bypassed unless the input source is explicitly verified safe and the bypass is documented.
- `bypassSecurityTrust*` methods MUST NOT be used without documented verification of the value source and security review for high-risk contexts.
- User-controlled values MUST NOT be rendered as HTML without sanitization.
- Direct DOM APIs, `ElementRef.nativeElement`, and third-party DOM manipulation MUST be treated as security-sensitive.
- Content Security Policy and Trusted Types SHOULD be enabled for high-risk or public applications where practical.
- Angular XSRF/CSRF protection SHOULD be enabled for cookie-based authenticated applications.
- XSRF/CSRF protection MUST also be enforced by the backend.
- Tokens and secrets MUST NOT be stored in insecure client-side storage without explicit risk acceptance.
- Angular security advisories MUST be reviewed during framework upgrades.

## Accessibility

- Angular applications MUST meet the organization's accessibility baseline.
- Semantic HTML and native controls MUST be preferred over custom interactive elements.
- Custom interactive components MUST support keyboard navigation, focus states, accessible names, and appropriate ARIA semantics.
- Form controls MUST have accessible labels and error messages.
- Error, loading, empty, and success states SHOULD be announced or exposed appropriately to assistive technologies.
- Route changes SHOULD manage focus so keyboard and screen-reader users can start at the new page content.
- Active navigation links SHOULD expose `aria-current` where appropriate.
- Dialogs, menus, popovers, and overlays MUST handle focus trapping, escape behaviour, and focus restoration.
- ARIA MUST NOT be used to compensate for incorrect HTML when a native semantic element is available.
- Deferred or dynamically inserted important content SHOULD be announced using appropriate live regions.
- Accessibility checks SHOULD be included in component tests, end-to-end tests, or review checklists for user-facing features.
- All interactive features MUST pass automated **AXE** accessibility audits.
- It MUST follow all WCAG AA minimums, including focus management, color contrast, and ARIA attributes.

## Styling

- Styles SHOULD be scoped to components using Angular's default view encapsulation.
- Global styles MUST be intentional and limited to resets, themes, typography, design tokens, and application-level layout.
- Style duplication across components SHOULD be extracted to design tokens, shared mixins, utility classes, or shared styles.
- Inline styles MUST NOT be used unless the value is dynamically bound from component state.
- Styling MUST NOT hide focus indicators without providing an accessible replacement.
- Component styles SHOULD support responsive layouts and user accessibility preferences such as reduced motion and text scaling.
- Design tokens SHOULD be used for shared color, spacing, typography, elevation, and motion values.

## Testing

- Business logic in services, validators, stores, and pure functions MUST be unit tested.
- Critical component behaviour MUST be tested.
- Bug fixes MUST include a regression test unless not practical; exceptions must be documented.
- Tests MUST NOT rely on implementation details such as private methods or internal signal storage.

### Component Tests

- `TestBed` MUST be used for component tests that involve Angular rendering, inputs, outputs, dependency injection, content projection, routing, or forms.
- Component harnesses SHOULD be used for interacting with Angular Material or CDK components.
- Presentational components SHOULD be tested by setting inputs and asserting rendered output and emitted events.
- Smart components SHOULD test orchestration by mocking injected services at the boundary.
- Accessibility-critical components SHOULD include accessibility assertions.

### Service Tests

- Services with no Angular-specific dependencies SHOULD be tested without `TestBed` for speed.
- Services with injected dependencies MUST use `TestBed` or manual `inject()` inside an injection context.
- HTTP calls MUST be tested using `HttpTestingController` or equivalent HTTP test utilities.
- State services MUST test state transitions through public methods.

### Routing Tests

- Guards MUST be tested for allowed, denied, and redirected outcomes.
- Route parameter handling SHOULD be tested for valid, missing, and invalid values.
- Critical navigation flows SHOULD be tested with router-aware test utilities.

### Signal And State Tests

- Signals MUST be tested by invoking public logic and reading observable outcomes.
- `computed()` values MUST be tested through their signal dependencies.
- `effect()` SHOULD NOT be tested directly unless unavoidable; test the external side effect or user-visible outcome instead.

### Test Quality

- Tests MUST be deterministic.
- Flaky tests MUST be fixed or removed from required CI gates until fixed.
- Test names MUST describe the scenario and expected outcome.
- Tests SHOULD use fake clocks or controlled schedulers for timing-sensitive behaviour.
- Tests MUST NOT use production data.

## Code Quality

- `any` MUST NOT be used unless the type is genuinely unknown and the usage is explicitly justified.
- `unknown` SHOULD be preferred over `any` at trust boundaries.
- Type assertions MUST be justified and kept close to validation.
- Unused imports, variables, files, and dependencies MUST be removed.
- Dead code, obsolete feature flags, and abandoned TODOs MUST be removed.
- `console.log` MUST NOT appear in production code.
- Public APIs, components, routes, and events SHOULD use consistent domain language.
- Errors MUST be handled in a way that produces safe user state and useful diagnostic information.

## Definition Of Done

An Angular feature is not complete unless:

- The application uses a supported Angular version.
- Modern Angular APIs are used where practical.
- Business logic is separated from components and templates.
- New component state uses signals where appropriate.
- Loading, success, empty, partial, and error states are handled where relevant.
- Authorization-sensitive flows enforce protection on the backend or trusted API layer.
- Forms are typed, accessible, and validated on both client and server where applicable.
- Routes are lazy-loaded unless eager loading is justified.
- Accessibility implications were considered and handled.
- Security implications were considered and handled.
- Performance impact was considered for routes, lists, images, bundles, and third-party libraries.
- SSR or hydration impact was considered when the application uses server or hybrid rendering.
- Tests cover critical business logic, component behaviour, state transitions, edge cases, and regressions.
- No unused code, debug logs, or unjustified `any` types were introduced.
- Formatting, linting, type checks, production build, and relevant CI checks pass.
