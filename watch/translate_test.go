package watch

import (
	"fmt"
	"strings"
	"testing"
)

func TestGoogleTranslate(t *testing.T) {
	// Test basic Spanish to English translation
	translated, err := googleTranslate("Hola, ¿cómo estás?", "es")
	if err != nil {
		t.Fatalf("googleTranslate failed: %v", err)
	}
	if translated == "" {
		t.Fatal("got empty translation")
	}
	lower := strings.ToLower(translated)
	if !strings.Contains(lower, "hello") && !strings.Contains(lower, "hi") {
		t.Errorf("unexpected translation: %q (expected something with 'hello')", translated)
	}
	fmt.Printf("ES->EN: %q\n", translated)
}

func TestGoogleTranslateWithSeparator(t *testing.T) {
	// Test that ||| separator survives translation
	input := "Bonjour ||| Comment allez-vous ||| Merci beaucoup"
	translated, err := googleTranslate(input, "fr")
	if err != nil {
		t.Fatalf("googleTranslate failed: %v", err)
	}

	parts := strings.Split(translated, "|||")
	fmt.Printf("FR->EN with separator: %q -> %q (%d parts)\n", input, translated, len(parts))

	if len(parts) < 2 {
		t.Errorf("separator was lost in translation, got %d parts: %q", len(parts), translated)
	}
}

func TestTranslateSRT(t *testing.T) {
	// Test SRT parsing and translation with a small Spanish SRT
	srt := `1
00:00:01,000 --> 00:00:04,000
Hola, bienvenidos a la fiesta.

2
00:00:05,000 --> 00:00:08,000
¿Cómo están todos esta noche?

3
00:00:09,000 --> 00:00:12,000
Estoy muy contento de estar aquí.
`

	translated, err := translateSRT(srt, "es")
	if err != nil {
		t.Fatalf("translateSRT failed: %v", err)
	}

	fmt.Printf("Translated SRT:\n%s\n", translated)

	// Verify structure is preserved
	blocks := strings.Split(strings.TrimSpace(translated), "\n\n")
	if len(blocks) != 3 {
		t.Errorf("expected 3 SRT blocks, got %d", len(blocks))
	}

	// Verify timestamps are preserved
	if !strings.Contains(translated, "00:00:01,000 --> 00:00:04,000") {
		t.Error("first timestamp was corrupted")
	}
	if !strings.Contains(translated, "00:00:05,000 --> 00:00:08,000") {
		t.Error("second timestamp was corrupted")
	}
}

func TestPickBestSubtitle(t *testing.T) {
	subtitles := SubtitlesResponse{
		{Url: "http://example.com/fi.srt", Language: "fi"},
		{Url: "http://example.com/es.srt", Language: "es"},
		{Url: "http://example.com/de.srt", Language: "de"},
		{Url: "http://example.com/ar.srt", Language: "ar"},
	}

	url, lang := pickBestSubtitle(subtitles)
	if lang != "es" {
		t.Errorf("expected Spanish (es) as best language, got %q", lang)
	}
	if url != "http://example.com/es.srt" {
		t.Errorf("wrong URL selected: %q", url)
	}
}
