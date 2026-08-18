package generator

import (
	"fmt"
	"path/filepath"
	"strconv"

	"ocr-quality-toolkit/internal/corpus"
	ocrhash "ocr-quality-toolkit/internal/hash"
	"ocr-quality-toolkit/internal/imageinfo"
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
		sourcePath := filepath.Join(
			root,
			filepath.FromSlash(parent.Image),
		)

		id := parent.ID + "__" + profile

		imagePath := filepath.ToSlash(
			filepath.Join(
				"synthetic",
				id+".png",
			),
		)

		targetPath := filepath.Join(
			root,
			filepath.FromSlash(imagePath),
		)

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

		format, width, height, err := imageinfo.Read(targetPath)
		if err != nil {
			return nil, fmt.Errorf(
				"read synthetic image info for %q: %w",
				parent.ID,
				err,
			)
		}

		checksum, err := ocrhash.FileSHA256(targetPath)
		if err != nil {
			return nil, fmt.Errorf(
				"calculate synthetic image SHA-256 for %q: %w",
				parent.ID,
				err,
			)
		}

		record.Format = format
		record.Width = width
		record.Height = height
		record.SHA256 = checksum

		records = append(records, record)
	}

	return records, nil
}
