package builder

import (
	"fmt"
	"html"
	"strings"
)

const anonimeArticlesPerPage = 6

func renderAnonimeHomePage(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-template"><div class="anonime-shell">`)
	body.WriteString(renderAnonimeHeader(site))
	body.WriteString(renderAnonimeHero(site))
	body.WriteString(renderAnonimeArticleIndex(site))
	body.WriteString(renderAnonimePricingCallout())
	body.WriteString(renderAnonimeFooter(site))
	body.WriteString(`</div></section>`)
	return renderDocument(site, "Home", body.String(), site.PublicBaseURL+"/")
}

func renderAnonimeArticlesPage(site renderedSite) string {
	return renderAnonimeArticlesPageNumber(site, 1)
}

func renderAnonimeArticlesPageNumber(site renderedSite, requestedPage int) string {
	articles, currentPage, totalPages := anonimeArticlePage(site.Articles, requestedPage)
	var body strings.Builder
	body.WriteString(`<section class="anonime-template"><div class="anonime-shell">`)
	body.WriteString(renderAnonimeHeader(site))
	body.WriteString(renderAnonimeArticlesHero(site))
	body.WriteString(renderAnonimeArticlesIndex(site, articles, currentPage, totalPages))
	body.WriteString(renderAnonimeFooter(site))
	body.WriteString(`</div></section>`)
	return renderDocument(site, "Articles", body.String(), anonimeArticlesPageURL(site, currentPage))
}

func renderAnonimeArticlePage(site renderedSite, article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-template"><div class="anonime-shell">`)
	body.WriteString(renderAnonimeHeader(site))
	body.WriteString(renderAnonimeArticleLayout(site, article))
	body.WriteString(renderAnonimeFooter(site))
	body.WriteString(`</div></section>`)
	return renderDocument(site, article.Title, body.String(), canonicalURL(site, article))
}

func renderAnonimeHeader(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<header class="anonime-header">`)
	body.WriteString(renderAnonimeHeaderBrand())
	body.WriteString(`<nav class="anonime-nav" aria-label="Primary">`)
	body.WriteString(`<a href="https://anonime.io/#top">Home</a>`)
	body.WriteString(`<a href="https://anonime.io/#how-it-works">How it Works</a>`)
	body.WriteString(`<a href="https://anonime.io/pricing">Plans &amp; Pricing</a>`)
	body.WriteString(`<a href="https://anonime.io/blog" aria-current="page">Blog</a>`)
	body.WriteString(`</nav>`)
	body.WriteString(`<div class="anonime-header-actions">`)
	body.WriteString(`<div class="anonime-theme-control"><input class="anonime-theme-checkbox" type="checkbox" id="anonime-theme-toggle"><label class="anonime-theme-toggle" for="anonime-theme-toggle"><span class="anonime-visually-hidden">Toggle dark mode</span><svg class="anonime-theme-moon" viewBox="0 0 24 24" width="19" height="19" fill="none" aria-hidden="true"><path d="M20.5 13.4A8.7 8.7 0 0 1 10.6 3.5a8.6 8.6 0 1 0 9.9 9.9Z" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg><svg class="anonime-theme-sun" viewBox="0 0 24 24" width="19" height="19" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="3.5" stroke="currentColor" stroke-width="1.8"/><path d="M12 2.5v2M12 19.5v2M4.5 12h-2M21.5 12h-2M5.3 5.3l1.4 1.4M17.3 17.3l1.4 1.4M18.7 5.3l-1.4 1.4M6.7 17.3l-1.4 1.4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg></label></div>`)
	body.WriteString(`<a class="anonime-login-link" href="https://app.anonime.io">Log in</a>`)
	body.WriteString(`<a class="anonime-button anonime-button-outline" href="https://app.anonime.io">Get Started</a>`)
	body.WriteString(`<details class="anonime-mobile-nav"><summary aria-label="Open navigation"><svg viewBox="0 0 24 24" width="20" height="20" fill="none" aria-hidden="true"><path d="M4 7h16M4 12h16M4 17h16" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg></summary><nav aria-label="Mobile navigation"><a href="https://anonime.io/#top">Home</a><a href="https://anonime.io/#how-it-works">How it Works</a><a href="https://anonime.io/pricing">Plans &amp; Pricing</a><a href="https://anonime.io/blog">Blog</a><a href="https://app.anonime.io">Log in</a><a href="https://app.anonime.io">Get Started</a></nav></details>`)
	body.WriteString(`</div></header>`)
	return body.String()
}

func renderAnonimeHeaderBrand() string {
	return `<a class="anonime-brand" href="/" aria-label="Anonime home"><img class="brand-logo" src="https://cdn.anonime.io/anonime-logo.svg" alt="Anonime"></a>`
}

func renderAnonimeBrand(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<a class="anonime-brand" href="https://anonime.io/#top">`)
	if logo := strings.TrimSpace(site.Site.LogoURL); logo != "" {
		body.WriteString(`<img class="anonime-brand-mark" src="`)
		body.WriteString(html.EscapeString(logo))
		body.WriteString(`" alt=""><span>`)
		body.WriteString(html.EscapeString(site.Site.Name))
		body.WriteString(`</span></a>`)
	} else {
		body.WriteString(`<span class="anonime-brand-mark" aria-hidden="true"><svg viewBox="0 0 24 24" width="20" height="20" fill="none" aria-hidden="true"><path d="M12 3.2c-4.5 0-7.4 3.3-7.4 7.7 0 4.6 3.1 8.2 7.4 10 4.3-1.8 7.4-5.4 7.4-10 0-4.4-2.9-7.7-7.4-7.7Z" fill="currentColor"/><path d="M9 7.5c1.7 2 2.2 4.5 1.7 7.6M15 7.5c-1.7 2-2.2 4.5-1.7 7.6" stroke="#38c493" stroke-width="2" stroke-linecap="round"/></svg></span><span>`)
		body.WriteString(html.EscapeString(site.Site.Name))
		body.WriteString(`</span></a>`)
	}
	return body.String()
}

func renderAnonimeHero(site renderedSite) string {
	var body strings.Builder
	featured := site.FeaturedArticles
	if len(featured) == 0 {
		featured = site.RecentArticles
	}
	featuredArticle := firstArticle(featured)

	body.WriteString(`<section class="anonime-hero anonime-section">`)
	body.WriteString(`<div class="anonime-hero-copy">`)
	body.WriteString(`<p class="anonime-eyebrow">The privacy blog</p>`)
	body.WriteString(`<h1 class="anonime-title">Thoughts that protect your right to <strong>privacy.</strong></h1>`)
	body.WriteString(`<p class="anonime-subtitle">Deep dives, real-world guides, and the latest thinking on privacy, identity, and taking back control of your data.</p>`)
	body.WriteString(`<div class="anonime-hero-actions">`)
	body.WriteString(`<a class="anonime-button anonime-button-primary" href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">Explore all articles <svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a>`)
	body.WriteString(`<a class="anonime-button anonime-button-secondary" href="https://anonime.io/#how-it-works">See features</a>`)
	body.WriteString(`</div></div>`)
	body.WriteString(renderAnonimeHeroArt())
	body.WriteString(`</section>`)

	if featuredArticle.Title != "" {
		body.WriteString(`<section class="anonime-featured anonime-section" id="features">`)
		body.WriteString(`<div class="anonime-featured-inner">`)
		body.WriteString(`<div class="anonime-featured-copy">`)
		body.WriteString(`<div class="anonime-badge-row"><span class="anonime-badge anonime-badge-solid">Featured</span><span class="anonime-badge anonime-badge-ghost">`)
		body.WriteString(html.EscapeString(topicLabel(featuredArticle)))
		body.WriteString(`</span></div>`)
		body.WriteString(`<h2 class="anonime-featured-title">`)
		body.WriteString(html.EscapeString(featuredArticle.Title))
		body.WriteString(`</h2>`)
		if strings.TrimSpace(featuredArticle.Excerpt) != "" {
			body.WriteString(`<p class="anonime-featured-summary">`)
			body.WriteString(html.EscapeString(featuredArticle.Excerpt))
			body.WriteString(`</p>`)
		}
		body.WriteString(renderAnonimeMeta(site, featuredArticle, true))
		body.WriteString(`</div>`)
		body.WriteString(`<a class="anonime-featured-hero" href="`)
		body.WriteString(html.EscapeString(articlePath(site, featuredArticle)))
		body.WriteString(`" aria-label="`)
		body.WriteString(html.EscapeString(featuredArticle.Title))
		body.WriteString(`"><div class="anonime-featured-lock" aria-hidden="true"></div></a>`)
		body.WriteString(`</div></section>`)
	}

	return body.String()
}

func renderAnonimeArticlesHero(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-hero anonime-section">`)
	body.WriteString(`<div class="anonime-hero-copy">`)
	body.WriteString(`<p class="anonime-eyebrow">The privacy blog</p>`)
	body.WriteString(`<h1 class="anonime-title">All <strong>Articles</strong></h1>`)
	body.WriteString(`<p class="anonime-subtitle">Insights on privacy, identity, and encrypted communication. Real-world guides and ideas to help you take back control.</p>`)
	body.WriteString(`</div>`)
	body.WriteString(renderAnonimeHeroArt())
	body.WriteString(`</section>`)
	return body.String()
}

func renderAnonimeArticleIndex(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-section" id="latest-insights">`)
	body.WriteString(`<div class="anonime-card anonime-panel">`)
	body.WriteString(`<div class="anonime-section-header">
	<h2 class="anonime-panel-title is-tight">Latest insights</h2>
	<a class="anonime-muted anonime-section-link" href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">View all articles</a></div></div>`)
	body.WriteString(`<div class="anonime-content-grid">`)
	body.WriteString(`<div class="anonime-article-stack">`)
	if len(site.RecentArticles) == 0 {
		body.WriteString(`<div class="anonime-panel"><p class="anonime-muted">No published articles yet.</p></div>`)
	} else {
		limit := 3
		if len(site.RecentArticles) < limit {
			limit = len(site.RecentArticles)
		}
		for _, article := range site.RecentArticles[:limit] {
			body.WriteString(renderAnonimeArticleCard(site, article, false))
		}
	}
	body.WriteString(`</div>`)
	body.WriteString(renderAnonimeSidebar(site))
	body.WriteString(`</div>`)
	body.WriteString(`</section>`)
	return body.String()
}

func renderAnonimeArticlesIndex(site renderedSite, articles []ArticleContent, currentPage, totalPages int) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-content-grid anonime-section" id="articles-grid">`)
	body.WriteString(`<div class="anonime-article-stack">`)
	if len(articles) == 0 {
		body.WriteString(`<div class="anonime-panel"><p class="anonime-muted">No articles available.</p></div>`)
	} else {
		for _, article := range articles {
			body.WriteString(renderAnonimeArticleCard(site, article, false))
		}
	}
	body.WriteString(renderAnonimePagination(site, currentPage, totalPages))
	body.WriteString(`</div>`)
	body.WriteString(renderAnonimeSidebar(site))
	body.WriteString(`</section>`)
	return body.String()
}

func renderAnonimeArticleLayout(site renderedSite, article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-article-layout anonime-section">`)
	body.WriteString(`<article class="anonime-article-main">`)
	body.WriteString(`<header class="anonime-article-header">`)
	body.WriteString(`<nav class="anonime-breadcrumbs" aria-label="Breadcrumb"><a href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">Blog</a><span aria-hidden="true">/</span><a href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`?category=`)
	body.WriteString(html.EscapeString(strings.ToLower(topicLabel(article))))
	body.WriteString(`">`)
	body.WriteString(html.EscapeString(topicLabel(article)))
	body.WriteString(`</a><span aria-hidden="true">/</span><span>`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</span></nav>`)
	body.WriteString(`<p class="anonime-eyebrow">Secure communication</p>`)
	body.WriteString(`<h1 class="anonime-article-heading">`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</h1>`)
	body.WriteString(renderAnonimeToolbar(article))
	body.WriteString(`</header>`)
	body.WriteString(renderAnonimeFigure(article))
	body.WriteString(`<div class="anonime-article-body">`)
	body.WriteString(renderMarkdown(article.ContentMarkdown))
	body.WriteString(`</div>`)
	body.WriteString(renderAnonimeFooterCard())
	body.WriteString(`</article>`)
	body.WriteString(renderAnonimeArticleSidebar(site, article))
	body.WriteString(`</section>`)
	return body.String()
}

func renderAnonimeSidebar(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<aside class="anonime-sidebar-panel">`)
	body.WriteString(`<section class="anonime-panel"><h2 class="anonime-panel-title">Explore by topic</h2><div class="anonime-topic-list">`)
	for _, topic := range topicCounts(site.Articles) {
		body.WriteString(`<div class="anonime-topic-entry"><strong>`)
		body.WriteString(html.EscapeString(topic.Label))
		body.WriteString(`</strong><span class="anonime-topic-count">`)
		body.WriteString(fmt.Sprint(topic.Count))
		body.WriteString(`</span></div>`)
	}
	body.WriteString(`</div></section>`)

	body.WriteString(`<section class="anonime-panel"><h2 class="anonime-panel-title">Popular articles</h2><div class="anonime-sidebar-list">`)
	for _, article := range popularArticles(site.Articles) {
		body.WriteString(renderAnonimeCompactArticle(site, article))
	}
	body.WriteString(`</div></section>`)
	body.WriteString(renderAnonimeCallout())
	body.WriteString(`</aside>`)
	return body.String()
}

func renderAnonimeArticleSidebar(site renderedSite, current ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<aside class="anonime-article-sidebar">`)
	body.WriteString(`<section class="anonime-panel"><h2 class="anonime-panel-title">Related Articles</h2><div class="anonime-article-related">`)
	for _, article := range relatedArticles(site.Articles, current.Slug) {
		body.WriteString(renderAnonimeCompactArticle(site, article))
	}
	body.WriteString(`</div></section>`)
	body.WriteString(renderAnonimeCallout())
	body.WriteString(`</aside>`)
	return body.String()
}

func renderAnonimeHeroArt() string {
	return `<div class="anonime-hero-art" aria-hidden="true"><div class="anonime-hero-grid"></div><div class="anonime-hero-small top-left"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M7.5 11.2V8.8a4.5 4.5 0 1 1 9 0v2.4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" /><rect x="6.5" y="11.2" width="11" height="8.5" rx="2.4" stroke="currentColor" stroke-width="1.8" /></svg></div><div class="anonime-hero-small top-right"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M4.5 7.5h15v9h-5.2l-2.3 2.3-2.3-2.3H4.5v-9Z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" /><path d="M8 10.5h8M8 13h5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" /></svg></div><div class="anonime-hero-card"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.2 4.7 6.8v5.2c0 4.9 3.5 8.5 7.3 9.8 3.8-1.3 7.3-4.9 7.3-9.8V6.8L12 3.2Z" stroke="currentColor" stroke-width="1.8" /></svg></div><div class="anonime-hero-small bottom-right"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="9.6" r="3" stroke="currentColor" stroke-width="1.8" /><path d="M5.5 19c1.3-3 4-4.5 6.5-4.5S17.2 16 18.5 19" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" /></svg></div></div>`
}

func renderAnonimePricingCallout() string {
	return `<section class="anonime-section anonime-section-spaced" id="pricing"><div class="anonime-callout dark"><div class="anonime-article-footer-card is-flat"><div><h3>Conversations are yours. Keep them that way.</h3><p>Private communication deserves a design that feels calm, deliberate, and secure.</p></div><div><a class="anonime-button anonime-button-primary" href="#get-started">Explore Anonime <svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a></div><div>Image here</div></div></div></section>`
}

func renderAnonimeCallout() string {
	return `<section class="anonime-callout"><div class="anonime-badge anonime-badge-ghost">Product mission</div><h2 class="anonime-callout-title">Privacy is not just our product. It is our mission.</h2><p class="anonime-callout-copy">Read more insights, guides, and real-world stories on taking back control of your digital life.</p><div class="anonime-callout-actions"><a class="anonime-button anonime-button-primary" href="#get-started">Explore all articles <svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a></div></section>`
}

func renderAnonimeArticleCard(site renderedSite, article ArticleContent, featured bool) string {
	var body strings.Builder
	if featured {
		body.WriteString(`<article class="anonime-card anonime-article-card featured">`)
	} else {
		body.WriteString(`<article class="anonime-card anonime-article-card">`)
	}
	body.WriteString(`<a href="`)
	body.WriteString(html.EscapeString(articlePath(site, article)))
	body.WriteString(`" aria-label="`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`"><div class="anonime-article-art`)
	if featured {
		body.WriteString(` featured`)
	}
	body.WriteString(`" data-art="`)
	body.WriteString(html.EscapeString(anonimeArtVariant(article)))
	body.WriteString(`" aria-hidden="true">`)
	if strings.TrimSpace(article.CoverImageURL) != "" {
		body.WriteString(`<img src="`)
		body.WriteString(html.EscapeString(article.CoverImageURL))
		body.WriteString(`" alt="">`)
	}
	body.WriteString(`</div></a>`)

	body.WriteString(`<div class="anonime-article-copy">`)
	body.WriteString(`<p class="anonime-article-category">`)
	body.WriteString(html.EscapeString(topicLabel(article)))
	body.WriteString(`</p>`)
	body.WriteString(`<h3 class="anonime-article-title"><a href="`)
	body.WriteString(html.EscapeString(articlePath(site, article)))
	body.WriteString(`">`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</a></h3>`)
	if strings.TrimSpace(article.Excerpt) != "" {
		body.WriteString(`<p class="anonime-article-excerpt">`)
		body.WriteString(html.EscapeString(article.Excerpt))
		body.WriteString(`</p>`)
	}
	body.WriteString(renderAnonimeMeta(site, article, false))
	body.WriteString(`</div></article>`)
	return body.String()
}

func renderAnonimeCompactArticle(site renderedSite, article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<article class="anonime-sidebar-item">`)
	body.WriteString(`<a href="`)
	body.WriteString(html.EscapeString(articlePath(site, article)))
	body.WriteString(`" aria-label="`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`"><div class="anonime-sidebar-thumb anonime-article-art" data-art="`)
	body.WriteString(html.EscapeString(anonimeArtVariant(article)))
	body.WriteString(`" aria-hidden="true"></div></a>`)
	body.WriteString(`<div><p class="anonime-article-category is-tight">`)
	body.WriteString(html.EscapeString(topicLabel(article)))
	body.WriteString(`</p><h3 class="anonime-sidebar-title"><a href="`)
	body.WriteString(html.EscapeString(articlePath(site, article)))
	body.WriteString(`">`)
	body.WriteString(html.EscapeString(article.Title))
	body.WriteString(`</a></h3><p class="anonime-sidebar-meta">`)
	body.WriteString(html.EscapeString(dateLabel(article)))
	body.WriteString(` <span class="anonime-meta-separator" aria-hidden="true"></span> `)
	body.WriteString(html.EscapeString(readTimeLabel(article)))
	body.WriteString(`</p></div></article>`)
	return body.String()
}

func renderAnonimeMeta(site renderedSite, article ArticleContent, featured bool) string {
	var body strings.Builder
	body.WriteString(`<div class="`)
	if featured {
		body.WriteString(`anonime-featured-meta`)
	} else {
		body.WriteString(`anonime-article-footer`)
	}
	body.WriteString(`">`)
	body.WriteString(renderAnonimeAvatar(article.AuthorName))
	body.WriteString(`<span>`)
	body.WriteString(html.EscapeString(authorLabel(article)))
	body.WriteString(`</span><span class="anonime-meta-separator" aria-hidden="true"></span><span>`)
	body.WriteString(html.EscapeString(dateLabel(article)))
	body.WriteString(`</span><span class="anonime-meta-separator" aria-hidden="true"></span><span>`)
	body.WriteString(html.EscapeString(readTimeLabel(article)))
	body.WriteString(`</span>`)
	if !featured {
		body.WriteString(`<a class="anonime-article-arrow" href="`)
		body.WriteString(html.EscapeString(articlePath(site, article)))
		body.WriteString(`" aria-label="Open `)
		body.WriteString(html.EscapeString(article.Title))
		body.WriteString(`"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a>`)
	}
	body.WriteString(`</div>`)
	return body.String()
}

func renderAnonimeToolbar(article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<div class="anonime-article-toolbar">`)
	body.WriteString(renderAnonimeAvatar(article.AuthorName))
	body.WriteString(`<span>`)
	body.WriteString(html.EscapeString(authorLabel(article)))
	body.WriteString(`</span><span class="anonime-meta-separator" aria-hidden="true"></span><span>`)
	body.WriteString(html.EscapeString(dateLabel(article)))
	body.WriteString(`</span><span class="anonime-meta-separator" aria-hidden="true"></span><span>`)
	body.WriteString(html.EscapeString(readTimeLabel(article)))
	body.WriteString(`</span>`)
	body.WriteString(`</div>`)
	return body.String()
}

func renderAnonimeFigure(article ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<figure class="anonime-article-figure">`)
	if strings.TrimSpace(article.CoverImageURL) != "" {
		body.WriteString(`<img src="`)
		body.WriteString(html.EscapeString(article.CoverImageURL))
		body.WriteString(`" alt="`)
		body.WriteString(html.EscapeString(article.Title))
		body.WriteString(`">`)
	} else {
		body.WriteString(`<div class="anonime-featured-hero" aria-hidden="true"><div class="anonime-featured-lock"></div></div>`)
	}
	body.WriteString(`</figure>`)
	return body.String()
}

func renderAnonimeFooterCard() string {
	return `<footer class="anonime-article-footer"><div class="anonime-article-footer-card"><div><h3>Your conversations are yours. Keep them that way.</h3><p>Encrypted communication helps protect the people and work that matter most.</p></div><a class="anonime-button" href="#get-started">Explore Anonime <svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a></div></footer>`
}

func anonimeArticlePage(articles []ArticleContent, requestedPage int) ([]ArticleContent, int, int) {
	totalPages := anonimeArticlePageCount(len(articles))
	if requestedPage < 1 {
		requestedPage = 1
	}
	if requestedPage > totalPages {
		requestedPage = totalPages
	}

	start := (requestedPage - 1) * anonimeArticlesPerPage
	if start >= len(articles) {
		return nil, requestedPage, totalPages
	}
	end := start + anonimeArticlesPerPage
	if end > len(articles) {
		end = len(articles)
	}
	return articles[start:end], requestedPage, totalPages
}

func anonimeArticlePageCount(articleCount int) int {
	if articleCount == 0 {
		return 1
	}
	return (articleCount + anonimeArticlesPerPage - 1) / anonimeArticlesPerPage
}

func anonimeArticlesPagePath(site renderedSite, page int) string {
	if page <= 1 {
		return articlesPath(site)
	}
	return strings.TrimRight(articlesPath(site), "/") + "/page/" + fmt.Sprint(page) + "/"
}

func anonimeArticlesPageURL(site renderedSite, page int) string {
	return site.PublicBaseURL + anonimeArticlesPagePath(site, page)
}

func renderAnonimePagination(site renderedSite, currentPage, totalPages int) string {
	if totalPages <= 1 {
		return ""
	}

	var body strings.Builder
	body.WriteString(`<nav class="anonime-pagination" aria-label="Articles pagination">`)
	if currentPage > 1 {
		body.WriteString(renderAnonimePaginationLink(site, currentPage-1, `aria-label="Previous page"`, `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" aria-hidden="true"><path d="M14 7 9 12l5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>`))
	}

	lastRenderedPage := 0
	for page := 1; page <= totalPages; page++ {
		if page != 1 && page != totalPages && (page < currentPage-1 || page > currentPage+1) {
			continue
		}
		if lastRenderedPage > 0 && page-lastRenderedPage > 1 {
			body.WriteString(`<span class="anonime-pagination-ellipsis" aria-hidden="true">&hellip;</span>`)
		}
		if page == currentPage {
			body.WriteString(`<span class="anonime-pagination-button is-active" aria-current="page">`)
			body.WriteString(fmt.Sprint(page))
			body.WriteString(`</span>`)
		} else {
			body.WriteString(renderAnonimePaginationLink(site, page, "", fmt.Sprint(page)))
		}
		lastRenderedPage = page
	}

	if currentPage < totalPages {
		body.WriteString(renderAnonimePaginationLink(site, currentPage+1, `aria-label="Next page"`, `<svg viewBox="0 0 24 24" width="14" height="14" fill="none" aria-hidden="true"><path d="M10 7 15 12l-5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg>`))
	}
	body.WriteString(`</nav>`)
	return body.String()
}

func renderAnonimePaginationLink(site renderedSite, page int, attributes, label string) string {
	return `<a class="anonime-pagination-button" href="` + html.EscapeString(anonimeArticlesPagePath(site, page)) + `" ` + attributes + `>` + label + `</a>`
}

func renderAnonimeFooter(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<footer class="anonime-footer">`)
	body.WriteString(`<div class="anonime-footer-main"><section class="anonime-footer-brand">`)
	body.WriteString(renderAnonimeBrand(site))
	body.WriteString(`<h2>Private communication<br>without <strong>permanent trails.</strong></h2></section>`)
	body.WriteString(`<nav class="anonime-footer-links" aria-label="Footer navigation"><div><h3>Product</h3><ul><li><a href="https://anonime.io/#security">Whispers</a></li><li><a href="https://anonime.io/#private-drops">Private Drops</a></li><li><a href="https://anonime.io/#how-it-works">Nyms</a></li><li><a href="https://anonime.io/pricing">Pricing</a></li></ul></div><div><h3>Company</h3><ul><li><a href="https://anonime.io/#about">About</a></li><li><a href="https://anonime.io/blog">Blog</a></li><li><a href="https://anonime.io/contact">Contact</a></li></ul></div><div><h3>Legal</h3><ul><li><a href="https://anonime.io/privacy">Privacy Policy</a></li><li><a href="https://anonime.io/terms">Terms of Service</a></li><li><a href="https://anonime.io/security">Security</a></li><li><a href="https://anonime.io/acceptable-usage-policy">Acceptable Use Policy</a></li><li><a href="https://anonime.io/data-processing">Data Processing</a></li></ul></div></nav></div>`)
	body.WriteString(`<section class="anonime-privacy-proof" aria-labelledby="anonime-privacy-proof-title"><p id="anonime-privacy-proof-title">Built with <strong>privacy by design.</strong></p><div class="anonime-privacy-features"><div class="anonime-privacy-feature"><span class="anonime-privacy-icon"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M8 10V7a4 4 0 0 1 8 0v3M6.5 10h11v9h-11z" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round"/></svg></span><span>End-to-end encrypted</span></div><div class="anonime-privacy-feature"><span class="anonime-privacy-icon"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M12 3.5 19 7v5.8c0 4-2.9 6.7-7 7.7-4.1-1-7-3.7-7-7.7V7l7-3.5Z" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round"/></svg></span><span>No personal data required</span></div><div class="anonime-privacy-feature"><span class="anonime-privacy-icon"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><circle cx="12" cy="12" r="7" stroke="currentColor" stroke-width="1.7"/><path d="m7 7 10 10" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"/></svg></span><span>No logs. No tracking.</span></div><div class="anonime-privacy-feature"><span class="anonime-privacy-icon"><svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><rect x="5" y="5" width="14" height="5" rx="1.5" stroke="currentColor" stroke-width="1.7"/><rect x="5" y="14" width="14" height="5" rx="1.5" stroke="currentColor" stroke-width="1.7"/><path d="M8 7.5h.01M8 16.5h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg></span><span>Hosted with privacy-first infra</span></div></div></section>`)
	body.WriteString(`<div class="anonime-footer-bottom"><p>&copy; 2026 Anonime, Inc. All rights reserved.</p><p class="anonime-footer-tagline">Your privacy is your right. <strong>We protect it every day.</strong></p><a class="footer-social-link" href="https://x.com/anonimehq" aria-label="X" target="_blank" rel="noopener noreferrer"><svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-4.714-6.231-5.401 6.231H2.74l7.73-8.835L1.254 2.25H8.08l4.258 5.622 5.906-5.622zm-1.161 17.52h1.833L7.084 4.126H5.117z"></path></svg></a></div></footer>`)
	return body.String()
}

func anonimeStyles() string {
	return `
		body.anonime-layout {
			background:
				radial-gradient(circle at 18% 2%, rgba(16, 178, 108, 0.12), transparent 24%),
				radial-gradient(circle at 84% 12%, rgba(16, 178, 108, 0.08), transparent 18%),
				linear-gradient(180deg, #ffffff 0%, #f7fbfa 100%);
			color: #0f1728;
			font-family: Matter, Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
		}
		body.anonime-layout .page-shell {
			width: min(1200px, calc(100vw - 32px));
			padding: 18px 0 32px;
		}
		.anonime-template {
			position: relative;
			isolation: isolate;
		}
		.anonime-template, .anonime-template * { box-sizing: border-box; }
		.anonime-template a { text-decoration: none; color: #057451; font-weight: 500; font-size: .9rem;}
		.anonime-template img { display: block; max-width: 100%; }
		.anonime-shell { width: 100%; }
		.anonime-header {
			display: flex;
			margin-inline: auto;
			align-items: center;
			justify-content: space-between;
			gap: 24px;
			padding: 10px 0 30px;
		}
		.anonime-brand { display: inline-flex; align-items: center; gap: 12px; font-size: 1.4rem; font-weight: 800; letter-spacing: -0.04em; }
		.brand-logo { width: 200px; height: auto; }
		.anonime-brand-mark {
			width: 38px; height: 38px; display: grid; place-items: center; flex: 0 0 auto; border-radius: 999px; object-fit: contain;
			background: #38c493; color: #071620; box-shadow: 0 10px 22px rgba(16, 178, 108, 0.18);
		}
		.anonime-nav { display: inline-flex; align-items: center; justify-content: center; gap: clamp(24px, 4vw, 52px); color: #172238; font-size: 0.9rem; }
		.anonime-nav a { position: relative; padding: 8px 2px; }
		.anonime-nav a[aria-current="page"] { color: #10b26c;}
		.anonime-nav a[aria-current="page"]::after {
			content: ""; position: absolute; left: 50%; bottom: 0; width: 100%; height: 2px; border-radius: 999px;
			transform: translateX(-50%); background: #10b26c;
		}
		.anonime-header-actions { display: inline-flex; align-items: center; justify-content: flex-end; gap: 18px; font-size: .9rem;}
		.anonime-login-link { color: #177a59; font-size: 0.9rem; font-weight: 700; white-space: nowrap;  padding: 0 22px;}
		.anonime-button-outline { color: #177a59; border: 1px solid #169568 !important; background: transparent; box-shadow: none; }
		.anonime-theme-control { position: relative; display: inline-flex; }
		.anonime-theme-checkbox { position: absolute; width: 1px; height: 1px; overflow: hidden; opacity: 0; }
		.anonime-theme-toggle { width: 42px; height: 42px; display: grid; place-items: center; border-radius: 999px; color: #536079; cursor: pointer; transition: color 160ms ease, background-color 160ms ease, transform 160ms ease; }
		.anonime-theme-toggle:hover { color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-theme-checkbox:focus-visible + .anonime-theme-toggle { outline: 3px solid rgba(16, 178, 108, 0.35); outline-offset: 3px; }
		.anonime-theme-sun { display: none; }
		.anonime-theme-checkbox:checked + .anonime-theme-toggle .anonime-theme-moon { display: none; }
		.anonime-theme-checkbox:checked + .anonime-theme-toggle .anonime-theme-sun { display: block; }
		.anonime-visually-hidden { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0 0 0 0); white-space: nowrap; border: 0; }
		.anonime-mobile-nav { position: relative; display: none; }
		.anonime-mobile-nav summary { width: 42px; height: 42px; display: grid; place-items: center; border-radius: 999px; color: #536079; cursor: pointer; list-style: none; }
		.anonime-mobile-nav summary::-webkit-details-marker { display: none; }
		.anonime-mobile-nav nav { position: absolute; z-index: 20; top: 50px; right: 0; width: 218px; display: grid; gap: 4px; padding: 12px; border: 1px solid rgba(18, 28, 43, 0.1); border-radius: 16px; background: rgba(255, 255, 255, 0.98); box-shadow: 0 18px 45px rgba(15, 23, 42, 0.14); }
		.anonime-mobile-nav nav a { padding: 10px 12px; border-radius: 10px; color: #172238; font-size: 0.9rem; }
		.anonime-mobile-nav nav a:hover { color: #087d58; background: rgba(16, 178, 108, 0.08); }
		.anonime-icon-button, .anonime-button, .anonime-pill, .anonime-pagination-button, .anonime-help-button {
			border: 0; cursor: pointer; transition: transform 160ms ease, box-shadow 160ms ease, border-color 160ms ease, background-color 160ms ease;
		}
		.anonime-icon-button {
			width: 42px; height: 42px; display: grid; place-items: center; border-radius: 999px; color: #5b6474; background: transparent;
		}
		.anonime-button {
			display: inline-flex; align-items: center; justify-content: center; gap: 10px; min-height: 42px; padding: 0 18px;
			border-radius: 6px; line-height: 42px; font-size: .9rem;
		}
		.anonime-button-primary { color: #fff !important; background: linear-gradient(180deg, #19c27b 0%, #0b9461 100%);}
		.anonime-button-secondary { color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-hero {
			display: grid; grid-template-columns: 58% 42%; align-items: center; gap: 2rem; padding: 26px 0 30px;
		}
		.anonime-eyebrow { margin: 0 0 14px; font-size: 0.78rem; letter-spacing: 0.18em; text-transform: uppercase; color: #10b26c; }
		.anonime-title, .anonime-featured-title, .anonime-article-heading { margin: 0; letter-spacing: -0.06em; }
		.anonime-title { font-size: clamp(2rem, 2vw, 3.75rem); line-height: 1.04; font-weight: 500; max-width: 20ch;}
		.anonime-title strong { color: #10b26c; font-style: normal; }
		.anonime-subtitle { margin: 18px 0 0; max-width: 46ch; font-size: 1.03rem; line-height: 1.8; color: #5b6474; }
		.anonime-hero-actions { display: flex; gap: 12px; flex-wrap: wrap; margin-top: 26px; }
		.anonime-hero-art {
			position: relative; min-height: 360px; border-radius: 34px;
			background:
				radial-gradient(circle at 16% 22%, rgba(16, 178, 108, 0.14), transparent 14%),
				radial-gradient(circle at 72% 72%, rgba(16, 178, 108, 0.1), transparent 16%),
				linear-gradient(180deg, rgba(16, 178, 108, 0.06), rgba(16, 178, 108, 0.02));
		}
		.anonime-hero-art::before, .anonime-hero-art::after, .anonime-featured-hero::before, .anonime-featured-hero::after {
			content: ""; position: absolute; border-radius: 24px; background: rgba(255, 255, 255, 0.68); box-shadow: 0 24px 60px rgba(14, 23, 38, 0.08); backdrop-filter: blur(12px);
		}
		.anonime-hero-art::before { top: 18px; right: 6px; width: 200px; height: 110px; transform: rotate(-6deg); }
		.anonime-hero-art::after { left: 24px; bottom: 22px; width: 220px; height: 124px; transform: rotate(10deg); }
		.anonime-hero-grid {
			position: absolute; inset: 0; border-radius: 34px;
			background:
				linear-gradient(rgba(16, 178, 108, 0.08) 1px, transparent 1px) 0 0 / 28px 28px,
				linear-gradient(90deg, rgba(16, 178, 108, 0.08) 1px, transparent 1px) 0 0 / 28px 28px;
			mask-image: radial-gradient(circle at center, black, transparent 70%); opacity: 0.24;
		}
		.anonime-hero-card, .anonime-hero-small {
			position: absolute; border-radius: 24px; background: rgba(255, 255, 255, 0.7); box-shadow: 0 24px 60px rgba(14, 23, 38, 0.08); backdrop-filter: blur(12px);
		}
		.anonime-hero-card { left: 28%; top: 28%; width: 240px; height: 190px; display: grid; place-items: center; }
		.anonime-hero-small { width: 72px; height: 72px; display: grid; place-items: center; }
		.anonime-hero-small.top-left { top: 48px; left: 38px; }
		.anonime-hero-small.top-right { top: 84px; right: 44px; }
		.anonime-hero-small.bottom-right { right: 62px; bottom: 48px; }
		.anonime-hero-card svg, .anonime-hero-small svg, .anonime-featured-hero svg { width: 42px; height: 42px; color: #10b26c; }
		.anonime-featured {
			margin-top: 8px; padding: 24px; border-radius: 10px; color: #fff;
			background: radial-gradient(circle at 76% 34%, rgba(16, 178, 108, 0.28), transparent 18%), linear-gradient(135deg, #0a1410 0%, #041f18 52%, #0d1514 100%);
		}
		.anonime-featured-inner { display: grid; grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.08fr); gap: 24px; align-items: center; }
		.anonime-featured-summary, .anonime-callout-copy, .anonime-article-excerpt, .anonime-subtitle { color: #5b6474; }
		.anonime-featured-copy { max-width: 470px; }
		.anonime-badge-row { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; margin-bottom: 18px; }
		.anonime-badge { display: inline-flex; align-items: center; min-height: 28px; padding: 0 12px; border-radius: 4px; font-size: 0.82rem; line-height: 28px; }
		.anonime-badge-solid { color: #fff; background: rgb(56 179 116); }
		.anonime-badge-ghost {background: transparent; color: #fff;}
		.anonime-featured-title { font-size: 1.8rem; line-height: 1.05; letter-spacing: -0.05em; font-weight: 500; max-width: 25ch; }
		.anonime-featured-summary { margin: 16px 0 0; max-width: 48ch; line-height: 1.8; color: rgba(255, 255, 255, 0.76); display: -webkit-box;
			-webkit-box-orient: vertical;
			-webkit-line-clamp: 3;
			line-clamp: 3;
			overflow: hidden; 
		}
		.anonime-featured-meta, .anonime-article-footer, .anonime-meta-inline, .anonime-article-feedback { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; margin-top: 18px; font-size: 0.88rem; }
		.anonime-meta-separator { width: 4px; height: 4px; border-radius: 999px; background: currentColor; opacity: 0.4; }
		.anonime-avatar { width: 32px; height: 32px; border-radius: 999px; object-fit: cover; border: 2px solid rgba(255, 255, 255, 0.24); }
		.anonime-avatar-fallback { display: grid; place-items: center; border: 0; background: rgba(16, 178, 108, 0.1); color: #059669; font-size: 0.72rem; }
		.anonime-featured-hero {
			position: relative; min-height: 314px; border-radius: 28px; overflow: hidden;
			background: radial-gradient(circle at 52% 44%, rgba(24, 243, 154, 0.18), transparent 14%), radial-gradient(circle at 50% 50%, rgba(16, 178, 108, 0.32), transparent 23%), linear-gradient(135deg, rgba(4, 17, 25, 0.9), rgba(4, 45, 33, 0.94));
		}
		.anonime-featured-hero::before { inset: 10% 18%; border: 1px solid rgba(24, 243, 154, 0.28); box-shadow: 0 0 0 20px rgba(24, 243, 154, 0.04), 0 0 0 44px rgba(24, 243, 154, 0.03), 0 0 0 76px rgba(24, 243, 154, 0.02); }
		.anonime-featured-hero::after { inset: 18% 36%; border: 2px solid rgba(24, 243, 154, 0.85); box-shadow: 0 0 0 8px rgba(24, 243, 154, 0.16), inset 0 0 26px rgba(24, 243, 154, 0.18); }
		.anonime-featured-lock { position: absolute; left: 50%; top: 50%; width: 138px; height: 138px; transform: translate(-50%, -50%); border-radius: 32px; border: 2px solid rgba(24, 243, 154, 0.86); background: linear-gradient(180deg, rgba(12, 57, 43, 0.92), rgba(4, 17, 25, 0.72)); box-shadow: 0 0 0 18px rgba(24, 243, 154, 0.04), 0 0 56px rgba(24, 243, 154, 0.3); }
		.anonime-featured-lock::before { content: ""; position: absolute; left: 50%; top: 40%; width: 58px; height: 42px; transform: translateX(-50%); border: 7px solid rgba(24, 243, 154, 0.86); border-bottom-width: 0; border-radius: 32px 32px 0 0; }
		.anonime-featured-lock::after { content: ""; position: absolute; left: 50%; top: 58%; width: 20px; height: 28px; transform: translateX(-50%); border-radius: 12px; background: rgba(24, 243, 154, 0.82); box-shadow: 0 0 24px rgba(24, 243, 154, 0.52); }
		.anonime-topic-row { display: flex; flex-wrap: wrap; gap: 10px; margin-top: 22px; }
		.anonime-pill {
			display: inline-flex; align-items: center; gap: 10px; min-height: 56px; padding: 0 18px; border-radius: 16px; border: 1px solid rgba(18, 28, 43, 0.1); color: #334155; background: rgba(255, 255, 255, 0.72); box-shadow: 0 4px 14px rgba(15, 23, 42, 0.03);
		}
		.anonime-pill.is-active { border-color: rgba(16, 178, 108, 0.18); color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-content-grid { display: grid; grid-template-columns: minmax(0, 1fr) 280px; gap: 24px; align-items: start; margin-top: 24px; }
		.anonime-panel { padding: 20px; }
		.anonime-panel-title { margin: 0 0 16px; font-size: 1rem; letter-spacing: -0.03em; font-weight: 500; }
		.anonime-panel-title.is-tight { margin-bottom: 0; }
		.anonime-section-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 16px; }
		.anonime-section-link { font-weight: 500; }
		.anonime-article-stack { display: grid; gap: 16px; }
		.anonime-article-card { display: grid; grid-template-areas: "media copy"; grid-template-columns: minmax(260px, 0.95fr) minmax(0, 1.05fr); gap: 18px; padding: 14px; overflow: hidden; }
		.anonime-article-card.is-compact { grid-template-columns: 72px minmax(0, 1fr); align-items: center; padding: 12px; }
		.anonime-article-art, .anonime-sidebar-thumb { position: relative; min-height: 176px; border-radius: 10px; overflow: hidden; background: linear-gradient(135deg, rgba(10, 20, 17, 0.94), rgba(7, 68, 52, 0.9));}
		.anonime-article-card > a:first-child { grid-area: media; min-height: 100%; }
		.anonime-article-card .anonime-article-art { height: 100%; min-height: 100%; }
		.anonime-article-card .anonime-article-art img { width: 100%; height: 100%; object-fit: cover; }
		.anonime-sidebar-thumb { min-height: 64px; border-radius: 14px; }
		.anonime-article-art::before, .anonime-article-art::after, .anonime-sidebar-thumb::before, .anonime-sidebar-thumb::after { content: ""; position: absolute; border-radius: 999px; background: rgba(16, 178, 108, 0.22); }
		.anonime-article-art::before { width: 110px; height: 110px; left: 50%; top: 16%; transform: translateX(-50%); box-shadow: 0 0 0 16px rgba(16, 178, 108, 0.12), 0 0 0 40px rgba(16, 178, 108, 0.08); }
		.anonime-article-art::after { width: 56px; height: 68px; left: 50%; top: 44%; transform: translateX(-50%); border-radius: 16px; background: rgba(255, 255, 255, 0.12); }
		.anonime-article-copy { grid-area: copy; display: flex; min-width: 0; padding: 4px 0 2px 6px; flex-direction: column;  gap: .5rem;}
		.anonime-article-category { margin: 0 0 8px; font-size: 0.72rem; letter-spacing: 0.08em; text-transform: uppercase; color: #10b26c; }
		.anonime-article-category.is-tight { margin-bottom: 4px; }
		.anonime-article-title { margin: 0; font-size: 1.34rem; line-height: 1.2; letter-spacing: -0.04em; }
		.anonime-article-title a {font-size: 1rem; font-weight: 500;}
		.anonime-article-excerpt { font-size: .9rem; margin: 12px 0; line-height: 1.72; display: -webkit-box; -webkit-box-orient: vertical; -webkit-line-clamp: 3; line-clamp: 3; overflow: hidden; }
		.anonime-article-footer { margin-top: auto; color: #5b6474; }
		.anonime-article-arrow { display: inline-grid; place-items: center; width: 34px; height: 34px; margin-left: auto; border-radius: 999px; color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-sidebar-item { display: grid; grid-template-columns: 84px minmax(0, 1fr); gap: 12px; align-items: stretch; }
		.anonime-sidebar-item > a { display: block; min-width: 0; }
		.anonime-sidebar-item .anonime-sidebar-thumb { width: 100%; height: 100%; min-height: 80px; border-radius: 10px; }
		.anonime-sidebar-item > div { display: flex; min-width: 0; flex-direction: column; }
		.anonime-sidebar-item .anonime-sidebar-title { display: -webkit-box; margin: 0; overflow: hidden; color: #0f1728; font-size: .92rem; font-weight: 600; line-height: 1.3; -webkit-box-orient: vertical; -webkit-line-clamp: 2; line-clamp: 2; }
		.anonime-sidebar-item .anonime-sidebar-title a { color: inherit; font: inherit; }
		.anonime-sidebar-item .anonime-sidebar-meta { display: flex; align-items: center; gap: 8px; margin: 8px 0 0; color: #5b6474; font-size: .78rem; white-space: nowrap; }
		.anonime-sidebar-panel, .anonime-article-sidebar { display: grid; gap: 16px; position: sticky; top: 20px; }
		.anonime-topic-list { display: grid; gap: 12px; }
		.anonime-topic-entry { display: flex; align-items: center; justify-content: space-between; gap: 12px; color: #334155; font-size: 0.95rem; }
		.anonime-topic-entry strong {font-weight: 500;}
		.anonime-topic-count { color: #10b26c;}
		.anonime-callout { position: relative; overflow: hidden; padding: 36px 22px; border-radius: 10px; background: linear-gradient(180deg, #ffffff 0%, #edf9f2 100%); }
		.anonime-callout.dark { margin-top: 1.5rem; color: #fff; background: radial-gradient(circle at 82% 24%, rgba(16, 178, 108, 0.24), transparent 22%), linear-gradient(135deg, #08111a 0%, #0a2a22 54%, #091518 100%); }
		.anonime-callout-title { margin: 0; max-width: 19ch; font-size: 1.3rem; line-height: 1.1; letter-spacing: -0.05em; font-weight: 500;}
		.anonime-callout-copy { margin: 12px 0 0; line-height: 1.7; }
		.anonime-callout.dark .anonime-callout-copy { color: rgba(255, 255, 255, 0.72); }
		.anonime-callout-actions { margin-top: 18px; }
		.anonime-pagination { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px; margin-top: 18px; }
		.anonime-pagination-button {display: grid; place-content: center;
			min-width: 36px; min-height: 36px; padding: 0 12px; border-radius: 6px; background: rgba(255, 255, 255, 0.9); color: #334155;
		}
		.anonime-pagination-button.is-active { border-color: rgba(16, 178, 108, 0.18); color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-sidebar-list {display: flex; flex-direction: column; gap: 1rem;}
		.anonime-footer { width: min(1530px, calc(100vw - 32px)); margin: 64px 0 0 50%; padding: 0; transform: translateX(-50%); }
		.anonime-footer-main { display: flex; align-items: flex-start; gap: 56px; padding-top: 52px; padding-bottom: 48px; border-bottom: 1px solid rgba(29, 77, 65, 0.12); }
		.anonime-footer-brand { max-width: 760px; }
		.anonime-footer-brand .anonime-brand { color: #35bb8d; font-size: 1.9rem; }
		.anonime-footer-brand h2 { margin: 64px 0 0; font-size: clamp(2.1rem, 3vw, 3.25rem); line-height: 1.08; letter-spacing: -0.05em; font-weight: 600;}
		.anonime-footer-brand h2 strong { color: #2c9b73; font-weight: 600; }
		.anonime-footer-links { flex: 1; display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); align-items: start; margin-left: auto; max-width: 820px;}
		.anonime-footer-links > div { min-height: 276px; padding: 2px 34px; border-left: 1px solid rgba(29, 77, 65, 0.1); }
		.anonime-footer-links > div:first-child { border-left: 0; }
		.anonime-footer-links h3 { margin: 0 0 26px; color: #26986f; font-size: 0.79rem; letter-spacing: 0.02em; text-transform: uppercase; }
		.anonime-footer-links h3::after { content: ""; display: block; width: 63px; height: 1px; margin-top: 12px; background: #26986f; }
		.anonime-footer-links ul { margin: 0; padding: 0; list-style: none; display: grid; gap: 22px; color: #43536b; font-size: 0.94rem; }
		.anonime-privacy-proof { min-height: 196px; padding: 29px 36px 30px; border-bottom: 1px solid rgba(29, 77, 65, 0.12); }
		.anonime-footer-links ul li a {font-weight: 400;}
		.anonime-privacy-proof > p { margin: 0 0 28px; text-align: center; color: #24324a; }
		.anonime-privacy-proof > p strong, .anonime-footer-tagline strong { color: #078f5d; font-weight: 500; }
		.anonime-privacy-features { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); }
		.anonime-privacy-feature { display: flex; align-items: center; gap: 16px; min-height: 84px; padding: 10px 22px; border-left: 1px solid rgba(29, 77, 65, 0.12); color: #43536b; font-size: 0.92rem; }
		.anonime-privacy-feature:first-child { border-left: 0; }
		.anonime-privacy-icon { flex: 0 0 auto;
			display: inline-grid;
			place-items: center;
			width: 52px;
			height: 52px;
			color: var(--green);
			background: radial-gradient(circle at 30% 28%, rgba(255, 255, 255, 0.96) 0 18%, transparent 19%), linear-gradient(180deg, #effaf5 0%, #d8f2e4 58%, #bce6cf 100%);
			border: 1px solid rgba(7, 152, 104, 0.12);
			border-radius: 12px;
			box-shadow: 4 10px 20px rgba(7, 152, 104, 0.12), 0 3px 8px rgba(16, 32, 51, 0.08), inset 0 2px 0 rgba(255, 255, 255, 0.88), inset 0 -6px 10px rgba(7, 152, 104, 0.06);
			transform: translateY(-1px);
			transition: transform 0.18s ease, box-shadow 0.18s ease, border-color 0.18s ease;
		}
		.anonime-privacy-icon svg { width: 23px; height: 23px; }
		.anonime-footer-bottom { display: grid; min-height: 156px; padding: 21px 55px 32px; grid-template-columns: 1fr auto 1fr; align-items: start; color: #74829a; font-size: 0.84rem; }
		.anonime-footer-bottom p { margin: 0; }
		.footer-social-link { display: grid; place-items: center; width: 36px; height: 36px; justify-self: end; color: #3d4b61; }
		.footer-social-link svg { width: 20px; height: 20px; }
		.anonime-footer-tagline { align-self: end; color: #63718a; text-align: center; font-size: 0.84rem; }
		.anonime-featured-summary, .anonime-muted { color: #5b6474; }
		.anonime-article-layout { display: grid; grid-template-columns: minmax(0, 1fr) 280px; gap: 24px; align-items: start; }
		.anonime-article-main { min-width: 0; }
		.anonime-article-header { padding: 24px 0 0; }
		.anonime-breadcrumbs { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-bottom: 16px; color: #5b6474; font-size: 0.86rem; }
		.anonime-breadcrumbs a { color: #059669; }
		.anonime-article-heading { max-width: 27ch; font-size: clamp(2rem, 2vw, 2rem); line-height: 1.5; font-weight: 500;}
		.anonime-article-intro { margin: 18px 0 0; max-width: 58ch; font-size: 1.03rem; line-height: 1.8; color: #5b6474; }
		.anonime-article-toolbar { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; margin-top: 20px; }
		.anonime-article-toolbar .anonime-button, .anonime-article-toolbar .anonime-icon-button, .anonime-help-button { border: 1px solid rgba(18, 28, 43, 0.1); background: rgba(255, 255, 255, 0.88); }
		.anonime-article-figure { margin: 22px 0 0; padding: 0; overflow: hidden; border-radius: 22px; box-shadow: 0 24px 60px rgba(14, 23, 38, 0.08); }
		.anonime-article-figure img { width: 100%; height: auto; }
		.anonime-article-body { padding: 22px 0 0; font-size: 1.02rem; line-height: 1.82; color: #334155; }
		.anonime-article-body h2, .anonime-article-body h3 { margin: 28px 0 10px; color: #0f1728; letter-spacing: -0.04em; }
		.anonime-article-body p { margin: 0 0 18px; }
		.anonime-article-body blockquote { margin: 22px 0; padding: 18px 20px; border-left: 3px solid #10b26c; border-radius: 16px; background: rgba(16, 178, 108, 0.06); color: #0f1728; }
		.anonime-inline-callout { display: grid; grid-template-columns: 56px minmax(0, 1fr); gap: 14px; align-items: center; margin: 22px 0; padding: 16px 18px; border-radius: 18px; background: rgba(16, 178, 108, 0.06); }
		.anonime-callout-icon { width: 56px; height: 56px; display: grid; place-items: center; border-radius: 16px; color: #10b26c; background: rgba(255, 255, 255, 0.9); }
		.anonime-article-footer-card { display: grid; grid-template-columns: minmax(0, 1fr) 1fr 1fr; gap: 20px; align-items: center; padding: 22px 24px; border-radius: 10px; color: #fff; background: linear-gradient(135deg, #081115 0%, #04231c 55%, #081513 100%); box-shadow: 0 30px 90px rgba(14, 23, 38, 0.16); }
		.anonime-article-footer-card.is-flat { padding: 0; background: transparent; box-shadow: none; }
		.anonime-article-footer-card p { margin: 0; max-width: 40ch; color: rgba(255, 255, 255, 0.8); line-height: 1.7; }
		.anonime-article-footer-card h3 { margin: 0 0 8px; font-size: 1.15rem; line-height: 1.12; letter-spacing: -0.05em; font-weight: 500;}
		.anonime-article-footer-card .anonime-button { background: linear-gradient(180deg, #19c27b 0%, #0b9461 100%); color: #fff; }
		.anonime-article-feedback { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; margin-top: 22px; }
		.anonime-help-button { display: inline-flex; align-items: center; gap: 8px; min-height: 40px; padding: 0 14px; border-radius: 12px; color: #334155; }
		.anonime-article-toc-list { display: grid; gap: 10px; margin: 0; padding-left: 0; list-style: none; color: #334155; }
		.anonime-article-toc-list li { position: relative; padding-left: 18px; line-height: 1.55; }
		.anonime-article-toc-list li::before { content: ""; position: absolute; left: 0; top: 0.7em; width: 6px; height: 6px; border-radius: 999px; background: rgba(148, 163, 184, 0.7); }
		.anonime-article-toc-list li.is-active { color: #10b26c; font-weight: 700; }
		.anonime-article-related { display: grid; gap: 14px; }
		.anonime-article-related-card { display: grid; grid-template-columns: 72px minmax(0, 1fr); gap: 12px; align-items: center; padding: 0 0 14px; border-bottom: 1px solid rgba(15, 23, 42, 0.08); }
		.anonime-article-related-card:last-child { padding-bottom: 0; border-bottom: 0; }
		.anonime-article-related-card .anonime-sidebar-title { font-size: 0.92rem; }
		.anonime-article-related-card .anonime-sidebar-meta { margin-top: 4px; color: #5b6474; }
		.anonime-muted { color: #5b6474; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) {
			color-scheme: dark;
			color: #edf6f2;
			background:
				radial-gradient(circle at 18% 2%, rgba(16, 178, 108, 0.12), transparent 24%),
				radial-gradient(circle at 84% 12%, rgba(16, 178, 108, 0.08), transparent 18%),
				linear-gradient(180deg, #08100e 0%, #0b1512 100%);
		}
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-brand,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-nav,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-title,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-title,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-panel-title,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-heading,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-callout-title,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-brand h2,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-privacy-proof > p { color: #edf6f2; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-subtitle,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-excerpt,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-footer,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-intro,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-body,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-topic-entry,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-toc-list,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-sidebar-meta,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-related-card .anonime-sidebar-meta,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-muted,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-links ul,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-privacy-feature,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-bottom,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-tagline { color: #a9b8b1; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-card,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-panel,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-callout:not(.dark) { border-color: rgba(191, 239, 215, 0.13); background: rgba(16, 28, 24, 0.94); box-shadow: 0 16px 40px rgba(0, 0, 0, 0.18); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-pill,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-pagination-button,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-help-button,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-toolbar .anonime-icon-button { border-color: rgba(191, 239, 215, 0.14); color: #c4d2cd; background: rgba(17, 31, 26, 0.9); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-button-outline { border-color: #36bd8d; color: #64d4ab; background: transparent; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-login-link,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-theme-toggle,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-mobile-nav summary { color: #64d4ab; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-mobile-nav nav { border-color: rgba(191, 239, 215, 0.14); background: rgba(17, 31, 26, 0.98); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-mobile-nav nav a { color: #edf6f2; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-hero-art::before,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-hero-art::after,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-hero-card,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-hero-small { background: rgba(18, 36, 29, 0.78); box-shadow: 0 24px 60px rgba(0, 0, 0, 0.22); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-body h2,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-body h3,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-body blockquote { color: #edf6f2; }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-article-body blockquote,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-inline-callout { background: rgba(16, 178, 108, 0.1); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-callout-icon,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-privacy-icon { border-color: rgba(73, 211, 159, 0.22); background: linear-gradient(180deg, rgba(29, 64, 51, 0.96), rgba(13, 42, 32, 0.96)); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-links > div,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-privacy-proof,
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-privacy-feature { border-color: rgba(191, 239, 215, 0.13); }
		body.anonime-layout:has(#anonime-theme-toggle:checked) .anonime-footer-x { color: #c5d2cd; }
		@media (max-width: 1240px) {
			.anonime-footer-main { padding-inline: 24px; grid-template-columns: minmax(390px, 1fr) minmax(520px, 1.25fr); gap: 38px; }
			.anonime-footer-links > div { padding-inline: 22px; }
			.anonime-privacy-feature { padding-inline: 20px; }
		}
		@media (max-width: 1020px) {
			.anonime-footer-main { min-height: 0; grid-template-columns: 1fr; gap: 42px; }
			.anonime-footer-brand h2 { margin-top: 38px; }
		}
		@media (max-width: 920px) {
			.anonime-content-grid, .anonime-article-layout { grid-template-columns: 1fr; }
			.anonime-sidebar-panel, .anonime-article-sidebar { position: static; }
			.anonime-featured-inner { grid-template-columns: 1fr; }
			.anonime-article-card { grid-template-areas: "copy" "media"; grid-template-columns: 1fr; }
			.anonime-hero { grid-template-columns: 1fr; }
		}
		@media (max-width: 820px) {
			.anonime-nav, .anonime-login-link, .anonime-button-outline { display: none; }
			.anonime-mobile-nav { display: block; }
			.anonime-featured, .anonime-panel, .anonime-card, .anonime-callout, .anonime-article-footer-card { border-radius: 10px; }
			.anonime-privacy-proof { padding-inline: 16px; }
			.anonime-privacy-features { grid-template-columns: repeat(2, minmax(0, 1fr)); }
			.anonime-privacy-feature:nth-child(3) { border-left: 0; }
		}
		@media (max-width: 620px) {
			.anonime-shell { width: 100%; }
			.anonime-header { gap: 12px; padding-bottom: 20px; }
			.anonime-header-actions { gap: 8px; }
			.anonime-brand { gap: 8px; font-size: 1.2rem; }
			.anonime-brand-mark { width: 34px; height: 34px; }
			.anonime-header-actions .anonime-button { min-height: 38px; padding-inline: 11px; font-size: 0.82rem; }
			.anonime-theme-toggle { width: 36px; height: 36px; }
			.anonime-login-link { font-size: 0.82rem; }
			.anonime-featured-copy, .anonime-panel, .anonime-card, .anonime-callout, .anonime-article-footer-card { padding-left: 16px; padding-right: 16px; }
			.anonime-article-card { padding: 12px; }
			.anonime-footer { margin-top: 42px; }
			.anonime-footer-main { padding: 0 14px 34px; }
			.anonime-footer-brand h2 { font-size: 2.1rem; }
			.anonime-footer-links { grid-template-columns: 1fr; }
			.anonime-footer-links > div { min-height: 0; padding: 24px 0; border-left: 0; border-top: 1px solid rgba(29, 77, 65, 0.12); }
			.anonime-footer-links > div:first-child { border-top: 0; }
			.anonime-footer-links ul { gap: 14px; }
			.anonime-footer-bottom { min-height: 176px; padding-inline: 14px; grid-template-columns: 1fr auto; }
			.anonime-footer-tagline { grid-column: 1 / -1; grid-row: 2; text-align: left; }
		}
		@media (max-width: 420px) {
			.anonime-privacy-features { grid-template-columns: 1fr; }
			.anonime-privacy-feature { border-left: 0; border-top: 1px solid rgba(29, 77, 65, 0.12); }
			.anonime-privacy-feature:first-child { border-top: 0; }
		}
	`
}

func renderAnonimeAvatar(name string) string {
	initials := initialsForName(name)
	if strings.TrimSpace(initials) == "" {
		initials = "A"
	}
	var body strings.Builder
	body.WriteString(`<span class="anonime-avatar anonime-avatar-fallback" aria-hidden="true">`)
	body.WriteString(html.EscapeString(initials))
	body.WriteString(`</span>`)
	return body.String()
}

func anonimeArtVariant(article ArticleContent) string {
	switch strings.ToLower(strings.TrimSpace(article.CategoryName)) {
	case "guides":
		return "key"
	case "security":
		return "shield"
	case "identity":
		return "fingerprint"
	case "product":
		return "cube"
	case "welcome":
		return "fingerprint"
	default:
		return "lock"
	}
}

func topicLabel(article ArticleContent) string {
	if strings.TrimSpace(article.CategoryName) != "" {
		return article.CategoryName
	}
	if strings.TrimSpace(article.CategorySlug) != "" {
		return humanizeSlug(article.CategorySlug)
	}
	return "Privacy Basics"
}

func authorLabel(article ArticleContent) string {
	if strings.TrimSpace(article.AuthorName) != "" {
		return article.AuthorName
	}
	return "Rene Carter"
}

func dateLabel(article ArticleContent) string {
	date := strings.TrimSpace(article.PublishedAt)
	if hasISODatePrefix(date) {
		return date[:10]
	}
	if date != "" {
		return date
	}
	return "Feb 7, 2026"
}

func hasISODatePrefix(value string) bool {
	if len(value) < len("2006-01-02") || value[4] != '-' || value[7] != '-' {
		return false
	}
	for index, character := range value[:10] {
		if index == 4 || index == 7 {
			continue
		}
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func readTimeLabel(article ArticleContent) string {
	words := len(strings.Fields(article.ContentMarkdown))
	if words == 0 {
		words = len(strings.Fields(article.Excerpt))
	}
	if words == 0 {
		return "5 min read"
	}
	minutes := words / 220
	if words%220 != 0 {
		minutes++
	}
	if minutes < 1 {
		minutes = 1
	}
	return fmt.Sprintf("%d min read", minutes)
}

type anonimeTopic struct {
	Label string
	Count int
}

type anonimeTocEntry struct {
	ID     string
	Label  string
	Active bool
}

func topicCounts(articles []ArticleContent) []anonimeTopic {
	order := []string{"Privacy Basics", "Security", "Anonymity", "Guides", "Product", "Identity", "Welcome"}
	counts := make(map[string]int, len(order))
	for _, article := range articles {
		counts[topicLabel(article)]++
	}
	topics := make([]anonimeTopic, 0, len(order))
	for _, label := range order {
		if count := counts[label]; count > 0 {
			topics = append(topics, anonimeTopic{Label: label, Count: count})
		}
	}
	if len(topics) == 0 {
		topics = []anonimeTopic{
			{Label: "Privacy Basics", Count: 12},
			{Label: "Security", Count: 10},
			{Label: "Anonymity", Count: 8},
			{Label: "Guides", Count: 11},
			{Label: "Product", Count: 6},
			{Label: "Identity", Count: 7},
			{Label: "Welcome", Count: 2},
		}
	}
	return topics
}

func anonimeTocEntries(article ArticleContent) []anonimeTocEntry {
	return []anonimeTocEntry{
		{ID: "individuals", Label: "Individuals Protecting Their Personal Lives", Active: true},
		{ID: "businesses", Label: "Businesses Handling Sensitive Information"},
		{ID: "journalists", Label: "Journalists, Activists, and Communities at Risk"},
		{ID: "bottom-line", Label: "The Bottom Line"},
	}
}

func popularArticles(articles []ArticleContent) []ArticleContent {
	if len(articles) == 0 {
		return nil
	}
	items := make([]ArticleContent, 0, 5)
	for _, article := range articles {
		items = append(items, article)
		if len(items) == 5 {
			break
		}
	}
	return items
}

func relatedArticles(articles []ArticleContent, currentSlug string) []ArticleContent {
	items := make([]ArticleContent, 0, 4)
	for _, article := range articles {
		if strings.EqualFold(article.Slug, currentSlug) {
			continue
		}
		items = append(items, article)
		if len(items) == 4 {
			break
		}
	}
	return items
}

func firstArticle(articles []ArticleContent) ArticleContent {
	if len(articles) == 0 {
		return ArticleContent{}
	}
	return articles[0]
}

func initialsForName(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "A"
	}
	if len(parts) == 1 {
		return strings.ToUpper(string(parts[0][0]))
	}
	return strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
}

func humanizeSlug(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	return strings.Join(strings.FieldsFunc(value, func(r rune) bool { return r == '-' || r == '_' }), " ")
}
