package corpus

// Record describes a single image in the OCR corpus.
type Record struct {
	ID         string   `json:"id"`
	ParentID   string   `json:"parent_id,omitempty"`
	Image      string   `json:"image"`
	References []string `json:"references"`
	Language   string   `json:"language"`
	Task       string   `json:"task"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Format     string   `json:"format"`
	Tags       []string `json:"tags"`
	SHA256     string   `json:"sha256"`
}
