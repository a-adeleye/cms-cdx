package builder

import (
	"fmt"
	"html"
	"strings"
)

func renderAnonimeHomePage(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-template"><div class="anonime-shell">`)
	body.WriteString(renderAnonimeHeader(site))
	body.WriteString(renderAnonimeHero(site))
	body.WriteString(renderAnonimeArticleIndex(site))
	body.WriteString(renderAnonimeFooter(site))
	body.WriteString(`</div></section>`)
	return renderDocument(site, "Home", body.String(), site.PublicBaseURL+"/")
}

func renderAnonimeArticlesPage(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-template"><div class="anonime-shell">`)
	body.WriteString(renderAnonimeHeader(site))
	body.WriteString(renderAnonimeArticlesHero(site))
	body.WriteString(renderAnonimeArticlesIndex(site))
	body.WriteString(renderAnonimeFooter(site))
	body.WriteString(`</div></section>`)
	return renderDocument(site, "Articles", body.String(), articlesURL(site))
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
	body.WriteString(`<a class="anonime-brand" href="/">`)
	body.WriteString(`<span class="anonime-brand-mark" aria-hidden="true">`)
	body.WriteString(`<svg viewBox="0 0 24 24" width="18" height="18" fill="none" aria-hidden="true"><path d="M12 2.5 20 6.8v6.8c0 4.6-3.4 7.7-8 8.9-4.6-1.2-8-4.3-8-8.9V6.8L12 2.5Z" fill="currentColor" /><path d="M9.2 12.2a2.8 2.8 0 1 1 5.6 0c0 1.1-.6 2.1-1.5 2.6l.3 2.3h-3.2l.3-2.3a2.8 2.8 0 0 1-1.5-2.6Z" fill="#fff" opacity=".92" /></svg>`)
	body.WriteString(`</span><span>anonime</span></a>`)
	body.WriteString(`<nav class="anonime-nav" aria-label="Primary">`)
	body.WriteString(`<a href="#features">Features</a>`)
	body.WriteString(`<a href="#pricing">Pricing</a>`)
	body.WriteString(`<a href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`" aria-current="page">Blog</a>`)
	body.WriteString(`</nav>`)
	body.WriteString(`<div class="anonime-header-actions">`)
	body.WriteString(`<button class="anonime-icon-button" type="button" aria-label="Toggle theme"><svg viewBox="0 0 24 24" width="19" height="19" fill="none" aria-hidden="true"><path d="M20.5 13.4A8.7 8.7 0 0 1 10.6 3.5a8.6 8.6 0 1 0 9.9 9.9Z" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg></button>`)
	body.WriteString(`<a class="anonime-button anonime-button-primary" href="#get-started">Get Started</a>`)
	body.WriteString(`</div></header>`)
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
	body.WriteString(`<a class="anonime-button anonime-button-secondary" href="#features">See features</a>`)
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

	body.WriteString(renderAnonimeTopicRow(site.Articles))
	body.WriteString(renderAnonimeArticleIndex(site))
	body.WriteString(renderAnonimePricingCallout())
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
	body.WriteString(renderAnonimeTopicPills())
	return body.String()
}

func renderAnonimeArticleIndex(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-content-grid anonime-section" id="latest-insights">`)
	body.WriteString(`<div class="anonime-article-stack">`)
	body.WriteString(`<div class="anonime-card anonime-panel">`)
	body.WriteString(`<div class="anonime-section-header"><h2 class="anonime-panel-title is-tight">Latest insights</h2><a class="anonime-muted anonime-section-link" href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">View all articles</a></div></div>`)
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
	body.WriteString(`</section>`)
	return body.String()
}

func renderAnonimeArticlesIndex(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-content-grid anonime-section" id="articles-grid">`)
	body.WriteString(`<div class="anonime-article-stack">`)
	if len(site.Articles) == 0 {
		body.WriteString(`<div class="anonime-panel"><p class="anonime-muted">No articles available.</p></div>`)
	} else {
		for _, article := range site.Articles {
			body.WriteString(renderAnonimeArticleCard(site, article, false))
		}
	}
	body.WriteString(renderAnonimePagination())
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
	if strings.TrimSpace(article.Excerpt) != "" {
		body.WriteString(`<p class="anonime-article-intro">`)
		body.WriteString(html.EscapeString(article.Excerpt))
		body.WriteString(`</p>`)
	}
	body.WriteString(renderAnonimeToolbar(article))
	body.WriteString(`</header>`)
	body.WriteString(renderAnonimeFigure(article))
	body.WriteString(`<div class="anonime-article-body">`)
	body.WriteString(renderMarkdown(article.ContentMarkdown))
	body.WriteString(`</div>`)
	body.WriteString(renderAnonimeFooterCard())
	body.WriteString(renderAnonimeFeedback())
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
	body.WriteString(`<section class="anonime-panel"><h2 class="anonime-panel-title">On this page</h2><ol class="anonime-article-toc-list">`)
	for _, entry := range anonimeTocEntries(current) {
		if entry.Active {
			body.WriteString(`<li class="is-active"><a href="#`)
		} else {
			body.WriteString(`<li><a href="#`)
		}
		body.WriteString(html.EscapeString(entry.ID))
		body.WriteString(`">`)
		body.WriteString(html.EscapeString(entry.Label))
		body.WriteString(`</a></li>`)
	}
	body.WriteString(`</ol></section>`)
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

func renderAnonimeTopicRow(articles []ArticleContent) string {
	var body strings.Builder
	body.WriteString(`<section class="anonime-section" aria-label="Article categories"><div class="anonime-topic-row">`)
	body.WriteString(`<button class="anonime-pill is-active" type="button">All Articles</button>`)
	for _, topic := range topicCounts(articles) {
		body.WriteString(`<button class="anonime-pill" type="button">`)
		body.WriteString(html.EscapeString(topic.Label))
		body.WriteString(`</button>`)
	}
	body.WriteString(`</div></section>`)
	return body.String()
}

func renderAnonimeTopicPills() string {
	return `<section class="anonime-section" aria-label="Article categories"><div class="anonime-topic-row"><button class="anonime-pill is-active" type="button">All Articles</button><button class="anonime-pill" type="button">Privacy Basics</button><button class="anonime-pill" type="button">Security</button><button class="anonime-pill" type="button">Anonymity</button><button class="anonime-pill" type="button">Guides</button><button class="anonime-pill" type="button">Product</button></div></section>`
}

func renderAnonimePricingCallout() string {
	return `<section class="anonime-section anonime-section-spaced" id="pricing"><div class="anonime-callout dark"><div class="anonime-badge anonime-badge-solid">Take back control</div><div class="anonime-article-footer-card is-flat"><div><h3>Conversations are yours. Keep them that way.</h3><p>Private communication deserves a design that feels calm, deliberate, and secure.</p></div><a class="anonime-button anonime-button-primary" href="#get-started">Explore Anonime <svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 12h10m0 0-4-4m4 4-4 4" stroke="currentColor" stroke-width="1.9" stroke-linecap="round" stroke-linejoin="round" /></svg></a></div></div></section>`
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
	body.WriteString(html.EscapeString(article.PublishedAt))
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
	body.WriteString(`<button class="anonime-icon-button" type="button" aria-label="Bookmark article"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M7 4.8h10v14.7l-5-3.2-5 3.2V4.8Z" stroke="currentColor" stroke-width="1.8" stroke-linejoin="round" /></svg></button>`)
	body.WriteString(`<button class="anonime-icon-button" type="button" aria-label="Share article"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M16.5 8.5a2.5 2.5 0 1 0-2.3-3.5 2.5 2.5 0 0 0 2.3 3.5ZM7.5 14.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5Zm9-2a2.5 2.5 0 1 0 2.2 3.5 2.5 2.5 0 0 0-2.2-3.5ZM9.5 15.6l5.3-2.8m-5.3-2.8 5.3 2.8" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg></button>`)
	body.WriteString(`<button class="anonime-icon-button" type="button" aria-label="Copy link"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M9.2 14.8 14.8 9.2m-5.4 6.8a3 3 0 0 1-4.2 0 3 3 0 0 1 0-4.2l2.1-2.1m8.7-.1 2.1-2.1a3 3 0 0 1 4.2 0 3 3 0 0 1 0 4.2l-2.1 2.1" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg></button>`)
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

func renderAnonimeFeedback() string {
	return `<div class="anonime-article-feedback"><span class="anonime-muted">Was this article helpful?</span><button class="anonime-help-button" type="button"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M9.5 11.5V20H6V11.5h3.5Zm1.8-.8 1.6-4.1A2 2 0 0 1 14.8 5h1.4a1.8 1.8 0 0 1 1.7 2.2l-1 4.2H21v2.4a4 4 0 0 1-4 4h-5.5a2 2 0 0 1-1.9-1.4l-1.8-5a1.8 1.8 0 0 1 1.5-2.5Z" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>Yes</button><button class="anonime-help-button" type="button"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M14.5 12.5V4H18v8.5h-3.5Zm-1.8.8-1.6 4.1A2 2 0 0 1 9.2 19H7.8a1.8 1.8 0 0 1-1.7-2.2l1-4.2H3v-2.4a4 4 0 0 1 4-4h5.5a2 2 0 0 1 1.9 1.4l1.8 5a1.8 1.8 0 0 1-1.5 2.5Z" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" /></svg>No</button></div>`
}

func renderAnonimePagination() string {
	return `<div class="anonime-pagination" aria-label="Pagination"><button class="anonime-pagination-button" type="button" aria-label="Previous page"><svg viewBox="0 0 24 24" width="14" height="14" fill="none" aria-hidden="true"><path d="M14 7 9 12l5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg></button><button class="anonime-pagination-button is-active" type="button">1</button><button class="anonime-pagination-button" type="button">2</button><button class="anonime-pagination-button" type="button">3</button><button class="anonime-pagination-button" type="button">4</button><button class="anonime-pagination-button" type="button">5</button><button class="anonime-pagination-button" type="button">...</button><button class="anonime-pagination-button" type="button">12</button><button class="anonime-pagination-button" type="button" aria-label="Next page"><svg viewBox="0 0 24 24" width="14" height="14" fill="none" aria-hidden="true"><path d="M10 7 15 12l-5 5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" /></svg></button></div>`
}

func renderAnonimeFooter(site renderedSite) string {
	var body strings.Builder
	body.WriteString(`<footer class="anonime-footer">`)
	body.WriteString(`<section class="anonime-footer-brand"><a class="anonime-brand" href="/"><span class="anonime-brand-mark" aria-hidden="true"><svg viewBox="0 0 24 24" width="18" height="18" fill="none" aria-hidden="true"><path d="M12 2.5 20 6.8v6.8c0 4.6-3.4 7.7-8 8.9-4.6-1.2-8-4.3-8-8.9V6.8L12 2.5Z" fill="currentColor" /><path d="M9.2 12.2a2.8 2.8 0 1 1 5.6 0c0 1.1-.6 2.1-1.5 2.6l.3 2.3h-3.2l.3-2.3a2.8 2.8 0 0 1-1.5-2.6Z" fill="#fff" opacity=".92" /></svg></span><span>anonime</span></a><p>Private communication without permanent trails.</p><div class="anonime-social-row" aria-label="Social links"><a class="anonime-social-link" href="#" aria-label="X"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M17.5 4.5 13.4 9l5.9 10.5h-4.2L10.9 12.8l-4.7 6.7H2.4l4.4-6.1L1 4.5h4.3l4.7 8.6 5.9-8.6h1.6Z" fill="currentColor" /></svg></a><a class="anonime-social-link" href="#" aria-label="LinkedIn"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M6.8 8.6H3.7V20h3.1V8.6ZM5.2 3.8A1.8 1.8 0 1 0 5.2 7a1.8 1.8 0 0 0 0-3.2ZM20.3 20h-3.1v-5.8c0-1.4 0-3.2-2-3.2s-2.3 1.6-2.3 3.1V20h-3.1V8.6h3v1.6h.1c.4-.8 1.4-1.8 3-1.8 3.2 0 3.8 2.1 3.8 4.8V20Z" fill="currentColor" /></svg></a><a class="anonime-social-link" href="#" aria-label="GitHub"><svg viewBox="0 0 24 24" width="16" height="16" fill="none" aria-hidden="true"><path d="M12 2.3a9.7 9.7 0 0 0-3 .5 9.8 9.8 0 0 0-1 .4 10 10 0 0 0-5.8 9.1c0 4.2 2.6 7.9 6.5 9.3.5.1.7-.2.7-.5v-1.8c-2.7.6-3.3-1.1-3.3-1.1-.4-1.1-1.1-1.4-1.1-1.4-.9-.6.1-.6.1-.6 1 0 1.6 1 1.6 1 .9 1.5 2.4 1.1 3 .8.1-.7.4-1.1.7-1.3-2.1-.2-4.3-1-4.3-4.5 0-1 .4-1.8 1-2.5-.1-.3-.4-1.2.1-2.5 0 0 .8-.2 2.6 1a8.8 8.8 0 0 1 4.8 0c1.8-1.2 2.6-1 2.6-1 .5 1.3.2 2.2.1 2.5.6.7 1 1.5 1 2.5 0 3.6-2.2 4.3-4.3 4.5.5.4.8 1 .8 2.1v3.1c0 .3.2.6.8.5a10 10 0 0 0 6.4-9.3 10 10 0 0 0-5.9-9.1 9.8 9.8 0 0 0-1-.4 9.7 9.7 0 0 0-3-.5Z" fill="currentColor" /></svg></a></div><p class="anonime-footer-note">(c) 2026 Anonime. Your identity, your rules.</p></section>`)
	body.WriteString(`<nav class="anonime-footer-links" aria-label="Footer navigation"><div><h3>Product</h3><ul><li><a href="#features">Features</a></li><li><a href="#pricing">Pricing</a></li><li><a href="#security">Security</a></li><li><a href="#nyns">Nyns</a></li></ul></div><div><h3>Resources</h3><ul><li><a href="`)
	body.WriteString(html.EscapeString(articlesPath(site)))
	body.WriteString(`">Blog</a></li><li><a href="#guides">Guides</a></li><li><a href="#help-center">Help Center</a></li><li><a href="#status">Status</a></li></ul></div><div><h3>Company</h3><ul><li><a href="#about">About</a></li><li><a href="#privacy">Privacy</a></li><li><a href="#terms">Terms</a></li><li><a href="#contact">Contact</a></li></ul></div></nav>`)
	body.WriteString(`<aside class="anonime-footer-card"><div class="anonime-badge anonime-badge-ghost">One platform</div><h3>One platform. Separate lives.</h3><p>Take control of your privacy with end-to-end encrypted communication.</p><a class="anonime-button anonime-button-primary" href="#get-started">Get Started</a></aside>`)
	body.WriteString(`</footer>`)
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
			font-family: Manrope, Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
		}
		body.anonime-layout .page-shell {
			width: min(1110px, calc(100vw - 32px));
			padding: 18px 0 32px;
		}
		.anonime-template {
			position: relative;
			isolation: isolate;
		}
		.anonime-template a { text-decoration: none; color: inherit; }
		.anonime-template img { display: block; max-width: 100%; }
		.anonime-shell { width: 100%; }
		.anonime-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 24px;
			padding: 10px 0 24px;
		}
		.anonime-brand { display: inline-flex; align-items: center; gap: 12px; font-size: 1.4rem; font-weight: 800; letter-spacing: -0.04em; }
		.anonime-brand-mark {
			width: 34px; height: 34px; display: grid; place-items: center; border-radius: 12px;
			background: linear-gradient(180deg, #19c27b 0%, #0b9461 100%); color: #fff; box-shadow: 0 10px 22px rgba(16, 178, 108, 0.24);
		}
		.anonime-nav { display: inline-flex; align-items: center; gap: 30px; color: #334155; }
		.anonime-nav a { position: relative; padding: 8px 2px; }
		.anonime-nav a[aria-current="page"] { color: #10b26c; font-weight: 700; }
		.anonime-nav a[aria-current="page"]::after {
			content: ""; position: absolute; left: 50%; bottom: -4px; width: 18px; height: 2px; border-radius: 999px;
			transform: translateX(-50%); background: #10b26c;
		}
		.anonime-header-actions { display: inline-flex; align-items: center; gap: 14px; }
		.anonime-icon-button, .anonime-button, .anonime-pill, .anonime-pagination-button, .anonime-help-button, .anonime-social-link {
			border: 0; cursor: pointer; transition: transform 160ms ease, box-shadow 160ms ease, border-color 160ms ease, background-color 160ms ease;
		}
		.anonime-icon-button {
			width: 42px; height: 42px; display: grid; place-items: center; border-radius: 999px; color: #5b6474; background: transparent;
		}
		.anonime-button {
			display: inline-flex; align-items: center; justify-content: center; gap: 10px; min-height: 42px; padding: 0 18px;
			border-radius: 12px; font-weight: 700; letter-spacing: -0.01em;
		}
		.anonime-button-primary { color: #fff; background: linear-gradient(180deg, #19c27b 0%, #0b9461 100%); box-shadow: 0 16px 34px rgba(16, 178, 108, 0.24); }
		.anonime-button-secondary { color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-hero {
			display: grid; grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr); align-items: center; gap: 40px; padding: 26px 0 30px;
		}
		.anonime-eyebrow { margin: 0 0 14px; font-size: 0.78rem; font-weight: 800; letter-spacing: 0.18em; text-transform: uppercase; color: #10b26c; }
		.anonime-title, .anonime-featured-title, .anonime-article-heading { margin: 0; letter-spacing: -0.06em; }
		.anonime-title { font-size: clamp(2.8rem, 5.6vw, 4.8rem); line-height: 0.98; }
		.anonime-title strong { color: #10b26c; font-style: italic; }
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
			margin-top: 8px; padding: 24px; border-radius: 32px; color: #fff;
			background: radial-gradient(circle at 76% 34%, rgba(16, 178, 108, 0.28), transparent 18%), linear-gradient(135deg, #0a1410 0%, #041f18 52%, #0d1514 100%);
			box-shadow: 0 30px 90px rgba(14, 23, 38, 0.16);
		}
		.anonime-featured-inner { display: grid; grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.08fr); gap: 24px; align-items: center; }
		.anonime-featured-summary, .anonime-callout-copy, .anonime-article-excerpt, .anonime-subtitle { color: #5b6474; }
		.anonime-featured-copy { max-width: 470px; }
		.anonime-badge-row { display: flex; flex-wrap: wrap; gap: 10px; align-items: center; margin-bottom: 18px; }
		.anonime-badge { display: inline-flex; align-items: center; min-height: 28px; padding: 0 12px; border-radius: 999px; font-size: 0.72rem; font-weight: 800; letter-spacing: 0.1em; text-transform: uppercase; }
		.anonime-badge-solid, .anonime-badge-ghost { color: #fff; background: rgba(16, 178, 108, 0.18); }
		.anonime-featured-title { font-size: clamp(2rem, 4vw, 3rem); line-height: 1.05; letter-spacing: -0.05em; }
		.anonime-featured-summary { margin: 16px 0 0; max-width: 48ch; line-height: 1.8; color: rgba(255, 255, 255, 0.76); }
		.anonime-featured-meta, .anonime-article-footer, .anonime-meta-inline, .anonime-article-feedback { display: flex; align-items: center; flex-wrap: wrap; gap: 10px; margin-top: 18px; font-size: 0.88rem; }
		.anonime-meta-separator { width: 4px; height: 4px; border-radius: 999px; background: currentColor; opacity: 0.4; }
		.anonime-avatar { width: 32px; height: 32px; border-radius: 999px; object-fit: cover; border: 2px solid rgba(255, 255, 255, 0.24); }
		.anonime-avatar-fallback { display: grid; place-items: center; border: 0; background: rgba(16, 178, 108, 0.1); color: #059669; font-size: 0.72rem; font-weight: 800; }
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
		.anonime-card, .anonime-panel, .anonime-callout, .anonime-footer-card, .anonime-article-footer-card { border: 1px solid rgba(18, 28, 43, 0.1); background: rgba(255, 255, 255, 0.92); box-shadow: 0 12px 30px rgba(15, 23, 42, 0.04); backdrop-filter: blur(8px); border-radius: 24px; }
		.anonime-panel { padding: 20px; }
		.anonime-panel-title { margin: 0 0 16px; font-size: 1rem; letter-spacing: -0.03em; }
		.anonime-panel-title.is-tight { margin-bottom: 0; }
		.anonime-section-header { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-bottom: 12px; }
		.anonime-section-link { font-weight: 700; }
		.anonime-article-stack { display: grid; gap: 16px; }
		.anonime-article-card { display: grid; grid-template-columns: minmax(260px, 0.95fr) minmax(0, 1.05fr); gap: 18px; padding: 14px; overflow: hidden; }
		.anonime-article-card.is-compact { grid-template-columns: 72px minmax(0, 1fr); align-items: center; padding: 12px; }
		.anonime-article-art, .anonime-sidebar-thumb { position: relative; min-height: 176px; border-radius: 18px; overflow: hidden; background: linear-gradient(135deg, rgba(10, 20, 17, 0.94), rgba(7, 68, 52, 0.9)); box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.06); }
		.anonime-sidebar-thumb { min-height: 64px; border-radius: 14px; }
		.anonime-article-art::before, .anonime-article-art::after, .anonime-sidebar-thumb::before, .anonime-sidebar-thumb::after { content: ""; position: absolute; border-radius: 999px; background: rgba(16, 178, 108, 0.22); }
		.anonime-article-art::before { width: 110px; height: 110px; left: 50%; top: 16%; transform: translateX(-50%); box-shadow: 0 0 0 16px rgba(16, 178, 108, 0.12), 0 0 0 40px rgba(16, 178, 108, 0.08); }
		.anonime-article-art::after { width: 56px; height: 68px; left: 50%; top: 44%; transform: translateX(-50%); border-radius: 16px; background: rgba(255, 255, 255, 0.12); }
		.anonime-article-copy { min-width: 0; padding: 4px 6px 2px 0; }
		.anonime-article-category { margin: 0 0 8px; font-size: 0.72rem; font-weight: 800; letter-spacing: 0.08em; text-transform: uppercase; color: #10b26c; }
		.anonime-article-category.is-tight { margin-bottom: 4px; }
		.anonime-article-title { margin: 0; font-size: 1.34rem; line-height: 1.2; letter-spacing: -0.04em; }
		.anonime-article-excerpt { margin: 12px 0 0; line-height: 1.72; }
		.anonime-article-footer { color: #5b6474; }
		.anonime-article-arrow { display: inline-grid; place-items: center; width: 34px; height: 34px; margin-left: auto; border-radius: 999px; color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-sidebar-panel, .anonime-article-sidebar { display: grid; gap: 16px; position: sticky; top: 20px; }
		.anonime-topic-list { display: grid; gap: 12px; }
		.anonime-topic-entry { display: flex; align-items: center; justify-content: space-between; gap: 12px; color: #334155; font-size: 0.95rem; }
		.anonime-topic-count { color: #10b26c; font-weight: 800; }
		.anonime-callout { position: relative; overflow: hidden; padding: 22px; border-radius: 24px; background: linear-gradient(180deg, #ffffff 0%, #edf9f2 100%); }
		.anonime-callout.dark { color: #fff; background: radial-gradient(circle at 82% 24%, rgba(16, 178, 108, 0.24), transparent 22%), linear-gradient(135deg, #08111a 0%, #0a2a22 54%, #091518 100%); }
		.anonime-callout-title { margin: 0; max-width: 13ch; font-size: 1.3rem; line-height: 1.1; letter-spacing: -0.05em; }
		.anonime-callout-copy { margin: 12px 0 0; line-height: 1.7; }
		.anonime-callout.dark .anonime-callout-copy { color: rgba(255, 255, 255, 0.72); }
		.anonime-callout-actions { margin-top: 18px; }
		.anonime-pagination { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px; margin-top: 18px; }
		.anonime-pagination-button {
			min-width: 36px; min-height: 36px; padding: 0 12px; border-radius: 11px; border: 1px solid rgba(18, 28, 43, 0.1); background: rgba(255, 255, 255, 0.9); color: #334155;
		}
		.anonime-pagination-button.is-active { border-color: rgba(16, 178, 108, 0.18); color: #059669; background: rgba(16, 178, 108, 0.08); }
		.anonime-footer { display: grid; grid-template-columns: minmax(0, 1fr) minmax(0, 1.1fr) minmax(0, 1fr); gap: 32px; margin-top: 24px; padding: 22px 0 12px; }
		.anonime-footer-brand { max-width: 260px; }
		.anonime-footer-brand p, .anonime-footer-card p, .anonime-featured-summary, .anonime-muted { color: #5b6474; }
		.anonime-social-row { display: flex; gap: 10px; margin-top: 20px; }
		.anonime-social-link { width: 38px; height: 38px; display: grid; place-items: center; border-radius: 999px; color: #334155; background: rgba(255, 255, 255, 0.84); border: 1px solid rgba(18, 28, 43, 0.1); }
		.anonime-footer-links { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 18px; }
		.anonime-footer-links h3 { margin: 0 0 12px; font-size: 0.92rem; }
		.anonime-footer-links ul { margin: 0; padding: 0; list-style: none; display: grid; gap: 10px; color: #5b6474; }
		.anonime-footer-card { align-self: start; padding: 18px; border-radius: 20px; }
		.anonime-footer-card h3 { margin: 10px 0 8px; font-size: 1.02rem; line-height: 1.25; }
		.anonime-footer-card p { margin: 0 0 14px; line-height: 1.7; }
		.anonime-footer-card .anonime-button { width: 100%; }
		.anonime-footer-note { margin-top: 6px; color: #5b6474; font-size: 0.84rem; }
		.anonime-article-layout { display: grid; grid-template-columns: minmax(0, 1fr) 280px; gap: 24px; align-items: start; }
		.anonime-article-main { min-width: 0; }
		.anonime-article-header { padding: 24px 0 0; }
		.anonime-breadcrumbs { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; margin-bottom: 16px; color: #5b6474; font-size: 0.86rem; }
		.anonime-breadcrumbs a { color: #059669; }
		.anonime-article-heading { max-width: 15ch; font-size: clamp(2.8rem, 5vw, 4.2rem); line-height: 0.97; }
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
		.anonime-article-footer-card { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 20px; align-items: center; padding: 22px 24px; border-radius: 22px; color: #fff; background: linear-gradient(135deg, #081115 0%, #04231c 55%, #081513 100%); box-shadow: 0 30px 90px rgba(14, 23, 38, 0.16); }
		.anonime-article-footer-card.is-flat { padding: 0; background: transparent; box-shadow: none; }
		.anonime-article-footer-card p { margin: 0; max-width: 32ch; color: rgba(255, 255, 255, 0.8); line-height: 1.7; }
		.anonime-article-footer-card h3 { margin: 0 0 8px; font-size: 1.35rem; line-height: 1.12; letter-spacing: -0.05em; }
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
		@media (max-width: 1080px) {
			.anonime-content-grid, .anonime-article-layout, .anonime-footer { grid-template-columns: 1fr; }
			.anonime-sidebar-panel, .anonime-article-sidebar { position: static; }
			.anonime-featured-inner, .anonime-article-card { grid-template-columns: 1fr; }
			.anonime-hero { grid-template-columns: 1fr; }
		}
		@media (max-width: 820px) {
			.anonime-header { flex-wrap: wrap; justify-content: center; }
			.anonime-nav { order: 3; width: 100%; justify-content: center; gap: 18px; }
			.anonime-featured, .anonime-panel, .anonime-card, .anonime-callout, .anonime-footer-card, .anonime-article-footer-card { border-radius: 20px; }
			.anonime-footer-links { grid-template-columns: 1fr; }
		}
		@media (max-width: 620px) {
			.anonime-shell { width: min(1110px, calc(100vw - 20px)); }
			.anonime-header-actions { width: 100%; justify-content: space-between; }
			.anonime-featured-copy, .anonime-panel, .anonime-card, .anonime-callout, .anonime-footer-card, .anonime-article-footer-card { padding-left: 16px; padding-right: 16px; }
			.anonime-article-card { padding: 12px; }
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
	if strings.TrimSpace(article.PublishedAt) != "" {
		return article.PublishedAt
	}
	return "Feb 7, 2026"
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
