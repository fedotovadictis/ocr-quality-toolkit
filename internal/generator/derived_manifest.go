package generator

import "ocr-quality-toolkit/internal/corpus"

// BuildDerivedManifest формирует записи manifest
// для производных изображений заданных профилей.
func BuildDerivedManifest(
	parent corpus.Record,
	profiles []string,
	seed string,
) []corpus.Record {
	records := make([]corpus.Record, 0, len(profiles))

	reference := ""
	if len(parent.References) > 0 {
		reference = parent.References[0]
	}

	for _, profile := range profiles {
		id := parent.ID + "__" + profile
		imagePath := "images/" + id + ".png"

		record := MakeManifestRecord(
			id,
			parent.ID,
			imagePath,
			reference,
			profile,
			seed,
		)

		records = append(records, record)
	}

	return records
}
