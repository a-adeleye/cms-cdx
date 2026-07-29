# Supromail template

`supromail` is a three-page Astro blog template converted from the standalone Supromail
blog design (`supromail/index.html`, `articles.html`, `article.html`, `blog.css`, `blog.js`).

## Entry points

- `landing.astro` — blog landing page
- `articles.astro` — all-articles listing
- `article.astro` — article read page

Each entry point accepts the existing CMS `site` and `article` fields. Landing and listing
pages also accept `articles`; the listing additionally accepts `page` and `totalPages`. When
those props are omitted, deterministic design fixtures are used for visual preview; an
explicitly empty `articles: []` renders the production empty state.

The article template renders a bounded Markdown subset (headings, paragraphs, emphasis,
inline code, fenced code blocks, links, lists, and blockquotes) through Astro components.
Raw HTML is never injected. Optional newer CMS fields such as category, author, publication
date, reading time, and cover image are used when present and receive deterministic display
fallbacks otherwise — an article without a cover image falls back to the gradient panel and
glyph from the source design.

Supplied CMS articles require a lowercase hyphenated `slug` and a non-empty `title`; the
read-page entry also requires non-empty `contentMarkdown`. Malformed content fails the build
with a safe validation message. Listing-only omissions use neutral derived values and never
borrow preview-fixture metadata.

## Rendering notes

- Every rule in `styles/supromail.css` is scoped to `.supromail-site`, because the source
  design uses generic class names (`.container`, `.section`, `.btn`, `.prose`) that would
  otherwise leak into a host document.
- Product-site navigation (Platform, Documentation, Pricing, FAQ, sign-in, brand and footer
  logos) is absolute against `supromail.com` / `app.supromail.com`. Blog navigation, article
  links, and pagination are root-relative and derived from `site.blogPath`.
- The shell script carries the behaviour from `blog.js`: mobile navigation, staggered scroll
  reveals, the frosted scrolled header, the listing category filter, the article reading
  progress bar, and table-of-contents highlighting. Each block is guarded by the presence of
  its own markup, so the landing page does not run article-only code.
- The production builder selects these entry points from the fixed `supromail` template
  registry. It passes a typed JSON file through `CMS_BUILD_DATA_FILE`, uses
  `CMS_TEMPLATE_ROOT` to locate this package, and must not wrap them in the generic
  `BaseLayout` header. Astro emits its source routes beneath `/articles`; the API moves that
  static subtree to the site's configured `blogPath` after a successful build. The generated
  root `sitemap.xml` uses that configured path, listing the landing page, blog landing,
  article listing, pagination, and published article URLs. A blog-local `sitemap.xml` moves
  with the static subtree and lists only blog URLs.
- This package is trusted repository code, not an arbitrary upload format. Do not add runtime
  dependency installation or server-side template code here.

## Verification

From the repository root:

```powershell
npm.cmd --prefix packages/templates/supromail test
```

The test command runs the pure content-model tests and a real Astro fixture build that
verifies formatted Markdown output, navigation scoping, and hostile-HTML escaping.
