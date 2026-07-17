package handlers

import "testing"

func TestValidateAIConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{name: "openai with environment reference", config: `{"provider":"openai","model":"gpt-4.1-mini","apiKeySecretRef":"OPENAI_API_KEY"}`},
		{name: "compatible provider requires secure base URL", config: `{"provider":"openai_compatible","model":"model","apiKeySecretRef":"AI_KEY","baseUrl":"http://example.com"}`, wantErr: true},
		{name: "rejects stored credentials", config: `{"provider":"openai","apiKeySecretRef":"sk-live-key"}`, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAIConfig(test.config)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateAIConfig() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
