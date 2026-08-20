package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"ocr-quality-toolkit/internal/corpus"
	ocrhash "ocr-quality-toolkit/internal/hash"
	"ocr-quality-toolkit/internal/imageinfo"
)

// GenerateCorpus создаёт синтетический OCR-корпус из входных текстов.
func GenerateCorpus(
	inputs []TextInput,
	options PageOptions,
	pages int,
	outputDir string,
) ([]corpus.Record, error) {
	if pages < 0 {
		return nil, fmt.Errorf("pages must not be negative")
	}

	if pages == 0 || pages > len(inputs) {
		pages = len(inputs)
	}

	imageDir := filepath.Join(outputDir, "images")

	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		return nil, fmt.Errorf(
			"create image directory: %w",
			err,
		)
	}

	selected := append([]TextInput(nil), inputs[:pages]...)

	sort.Slice(selected, func(i, j int) bool {
		return selected[i].ID < selected[j].ID
	})

	records := make([]corpus.Record, 0, len(selected))

	for _, input := range selected {
		id := MakePageID(
			input.ID,
			input.Text,
			options,
		)

		imagePath := filepath.ToSlash(
			filepath.Join(
				"images",
				id+".png",
			),
		)

		targetPath := filepath.Join(
			outputDir,
			filepath.FromSlash(imagePath),
		)

		page, err := GeneratePage(
			input.Text,
			options,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"generate page %q: %w",
				input.ID,
				err,
			)
		}

		if err := SavePNG(targetPath, page); err != nil {
			return nil, fmt.Errorf(
				"save page %q: %w",
				input.ID,
				err,
			)
		}

		format, width, height, err := imageinfo.Read(targetPath)
		if err != nil {
			return nil, fmt.Errorf(
				"read generated image info %q: %w",
				input.ID,
				err,
			)
		}

		checksum, err := ocrhash.FileSHA256(targetPath)
		if err != nil {
			return nil, fmt.Errorf(
				"calculate generated image SHA-256 %q: %w",
				input.ID,
				err,
			)
		}

		record := corpus.Record{
			ID:         id,
			Image:      imagePath,
			References: []string{input.Text},
			Language:   "ru",
			Task:       "full-page OCR ru",
			Width:      width,
			Height:     height,
			Format:     format,
			Tags:       []string{"synthetic"},
			SHA256:     checksum,
		}

		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})

	manifestPath := filepath.Join(
		outputDir,
		"manifest.jsonl",
	)

	if err := corpus.WriteJSONL(manifestPath, records); err != nil {
		return nil, fmt.Errorf(
			"write generated manifest: %w",
			err,
		)
	}

	return records, nil
}
