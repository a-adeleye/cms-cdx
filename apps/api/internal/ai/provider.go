package ai

type Idea struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Outline struct {
	Sections []string `json:"sections"`
}

type Draft struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SEOResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GenerateIdeasInput struct {
	Topic string `json:"topic"`
}

type GenerateOutlineInput struct {
	Title string `json:"title"`
}

type GenerateDraftInput struct {
	Title string `json:"title"`
}

type GenerateSEOInput struct {
	Title string `json:"title"`
}

type GenerateTagsInput struct {
	Title string `json:"title"`
}

type GenerateImagePromptInput struct {
	Title string `json:"title"`
}

type AIProvider interface {
	GenerateIdeas(input GenerateIdeasInput) ([]Idea, error)
	GenerateOutline(input GenerateOutlineInput) (Outline, error)
	GenerateDraft(input GenerateDraftInput) (Draft, error)
	GenerateSEO(input GenerateSEOInput) (SEOResult, error)
	GenerateTags(input GenerateTagsInput) ([]string, error)
	GenerateImagePrompt(input GenerateImagePromptInput) (string, error)
}

type NoopProvider struct{}

func (NoopProvider) GenerateIdeas(input GenerateIdeasInput) ([]Idea, error) {
	return []Idea{{Title: input.Topic, Description: "Draft idea from placeholder provider"}}, nil
}

func (NoopProvider) GenerateOutline(input GenerateOutlineInput) (Outline, error) {
	return Outline{Sections: []string{"Intro", "Body", "Conclusion"}}, nil
}

func (NoopProvider) GenerateDraft(input GenerateDraftInput) (Draft, error) {
	return Draft{Title: input.Title, Content: "AI draft placeholder"}, nil
}

func (NoopProvider) GenerateSEO(input GenerateSEOInput) (SEOResult, error) {
	return SEOResult{Title: input.Title, Description: "SEO placeholder"}, nil
}

func (NoopProvider) GenerateTags(input GenerateTagsInput) ([]string, error) {
	return []string{"example", "placeholder"}, nil
}

func (NoopProvider) GenerateImagePrompt(input GenerateImagePromptInput) (string, error) {
	return "Editorial image prompt for " + input.Title, nil
}
