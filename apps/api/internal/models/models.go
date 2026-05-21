package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name,omitempty"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Site struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Slug           string         `json:"slug"`
	Domain         string         `json:"domain,omitempty"`
	BlogPath       string         `json:"blog_path"`
	Status         string         `json:"status"`
	TemplateKey    string         `json:"template_key"`
	ThemeConfig    map[string]any `json:"theme_config"`
	DeployProvider string         `json:"deploy_provider,omitempty"`
	DeployConfig   map[string]any `json:"deploy_config"`
	AIConfig       map[string]any `json:"ai_config"`
	StorageConfig  map[string]any `json:"storage_config"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Page struct {
	ID          string         `json:"id"`
	SiteID      string         `json:"site_id"`
	Type        string         `json:"type"`
	Title       string         `json:"title"`
	Slug        string         `json:"slug,omitempty"`
	ContentJSON map[string]any `json:"content_json"`
	SEOTitle    string         `json:"seo_title,omitempty"`
	SEODesc     string         `json:"seo_description,omitempty"`
	Status      string         `json:"status"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type LandingSection struct {
	ID           string         `json:"id"`
	SiteID       string         `json:"site_id"`
	SectionKey   string         `json:"section_key"`
	Title        string         `json:"title,omitempty"`
	Subtitle     string         `json:"subtitle,omitempty"`
	ContentJSON  map[string]any `json:"content_json"`
	DisplayOrder int            `json:"display_order"`
	IsEnabled    bool           `json:"is_enabled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Article struct {
	ID             string         `json:"id"`
	SiteID         string         `json:"site_id"`
	AuthorID       string         `json:"author_id,omitempty"`
	CategoryID     string         `json:"category_id,omitempty"`
	Title          string         `json:"title"`
	Slug           string         `json:"slug"`
	Excerpt        string         `json:"excerpt,omitempty"`
	ContentMarkdown string        `json:"content_markdown"`
	CoverImageURL  string         `json:"cover_image_url,omitempty"`
	Status         string         `json:"status"`
	IsFeatured     bool           `json:"is_featured"`
	PublishedAt    *time.Time     `json:"published_at,omitempty"`
	SEOTitle       string         `json:"seo_title,omitempty"`
	SEODescription string         `json:"seo_description,omitempty"`
	CanonicalURL   string         `json:"canonical_url,omitempty"`
	GeneratedByAI  bool           `json:"generated_by_ai"`
	HumanReviewed  bool           `json:"human_reviewed"`
	AIPrompt       string         `json:"ai_prompt,omitempty"`
	AIModel        string         `json:"ai_model,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type Build struct {
	ID           string     `json:"id"`
	SiteID       string     `json:"site_id"`
	Status       string     `json:"status"`
	BuildType    string     `json:"build_type"`
	Logs         string     `json:"logs,omitempty"`
	OutputPath   string     `json:"output_path,omitempty"`
	DeployProvider string   `json:"deploy_provider,omitempty"`
	DeployStatus string     `json:"deploy_status,omitempty"`
	DeployURL    string     `json:"deploy_url,omitempty"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type AuditLog struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id,omitempty"`
	SiteID     string         `json:"site_id,omitempty"`
	Action     string         `json:"action"`
	EntityType string         `json:"entity_type,omitempty"`
	EntityID   string         `json:"entity_id,omitempty"`
	Metadata   map[string]any `json:"metadata"`
	CreatedAt  time.Time      `json:"created_at"`
}

