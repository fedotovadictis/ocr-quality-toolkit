package corpus

type mwsRecord struct {
	FileName    string   `json:"file_name"`
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	DatasetName string   `json:"dataset_name"`
	Question    string   `json:"question"`
	Answers     []string `json:"answers"`
}
