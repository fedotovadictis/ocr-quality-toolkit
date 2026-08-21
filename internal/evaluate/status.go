package evaluate

type Status string

const (
	StatusSuccess           Status = "success"
	StatusEngineError       Status = "engine_error"
	StatusMissingHypothesis Status = "missing_hypothesis"
)
