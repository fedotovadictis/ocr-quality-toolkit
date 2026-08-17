package generator

import (
	"fmt"
	"path/filepath"
	"strconv"

	"ocr-quality-toolkit/internal/corpus"
)

// BuildSyntheticWorkset создаёт synthetic-изображения
// и соответствующие записи manifest.
func BuildSyntheticWorkset(
	root string,
	parents []corpus.Record,
	profile string,
	seed int64,
) ([]corpus.Record, error) {
	records := make([]corpus.Record, 0, len(parents))

	for _, parent := range parents {
		sourcePath := filepath.Join(root, filepath.FromSlash(parent.Image))

		id := parent.ID + "__" + profile
		imagePath := filepath.ToSlash(
			filepath.Join("synthetic", id+".png"),
		)

		targetPath := filepath.Join(root, filepath.FromSlash(imagePath))

		if err := BuildSyntheticImage(
			sourcePath,
			targetPath,
			profile,
			seed,
		); err != nil {
			return nil, fmt.Errorf(
				"build synthetic image for %q: %w",
				parent.ID,
				err,
			)
		}

		record, err := BuildSyntheticRecord(
			parent,
			imagePath,
			profile,
			strconv.FormatInt(seed, 10),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"build synthetic record for %q: %w",
				parent.ID,
				err,
			)
		}

		records = append(records, record)
	}

	return records, nil
}
