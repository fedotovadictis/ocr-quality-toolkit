package corpus

// Hypothesis describes OCR result for a single corpus record
type Hypothesis struct {
	ID         string `json:"id"`
	Text       string `json:"text,omitempty"`
	Engine     string `json:"engine"`
	Model      string `json:"model"`
	DurationMS int    `json:"duration_ms,omitempty"`
	Error      string `json:"error,omitempty"`
}
