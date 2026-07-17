# Anonime template

`anonime` is a three-page Astro blog template matching the supplied Anonime privacy-blog designs.

## Entry points

- `landing.astro` — blog landing page
- `articles.astro` — all-articles listing
- `article.astro` — article read page

Each entry point accepts the existing CMS `site` and `article` fields. Landing and listing pages also accept `articles`. When those props are omitted, deterministic design fixtures are used for visual preview; an explicitly empty `articles: []` renders the production empty state.

The article template renders a bounded Markdown subset (headings, paragraphs, emphasis, links, lists, and blockquotes) through Astro components. Raw HTML is never injected. Optional newer CMS fields such as category, author, publication date, reading time, and cover image are used when present and receive deterministic display fallbacks otherwise.

Supplied CMS articles require a lowercase hyphenated `slug` and a non-empty `title`; the read-page entry also requires non-empty `contentMarkdown`. Malformed content fails the build with a safe validation message. Listing-only omissions use neutral derived values and never borrow preview-fixture metadata.

The template renders the production Anonime header and footer with absolute product-site navigation. Its accessible theme control switches every page between light and dark themes, uses the operating-system preference on first visit, and persists an explicit choice in browser storage. The CSS-only control path still works when a preview sandbox blocks scripts.

A builder integration should select the three entry points from a fixed template registry and must not wrap them in the generic `BaseLayout` header.

## Verification

From the repository root:

```powershell
npm.cmd --prefix packages/templates/anonime test
```

The test command runs the pure content-model tests and a real Astro fixture build that verifies formatted Markdown output and hostile-HTML escaping.
