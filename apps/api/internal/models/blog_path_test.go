package models

import "testing"

func TestCanonicalBlogPathNormalizesValidPaths(t *testing.T) {
	for _, testCase := range []struct {
		input string
		want  string
	}{
		{input: "", want: "/articles"},
		{input: "articles/", want: "/articles"},
		{input: "/blog/news/", want: "/blog/news"},
	} {
		got, err := CanonicalBlogPath(testCase.input)
		if err != nil {
			t.Fatalf("CanonicalBlogPath(%q) returned error: %v", testCase.input, err)
		}
		if got != testCase.want {
			t.Fatalf("CanonicalBlogPath(%q) = %q, want %q", testCase.input, got, testCase.want)
		}
	}
}

func TestCanonicalBlogPathRejectsUnsafePaths(t *testing.T) {
	for _, input := range []string{"/", "/blog/../private", "/blog//news"} {
		if _, err := CanonicalBlogPath(input); err == nil {
			t.Fatalf("expected CanonicalBlogPath(%q) to reject an unsafe path", input)
		}
	}
}
