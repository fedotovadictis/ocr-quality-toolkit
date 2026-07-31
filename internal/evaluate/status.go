package evaluate

type Status string

const (
	StatusSuccess           Status = "success"
	StatusOCRError          Status = "ocr_error"
	StatusMissingHypothesis Status = "missing_hypothesis"
)
