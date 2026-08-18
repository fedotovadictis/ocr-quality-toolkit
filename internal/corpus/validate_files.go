package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ValidateRecordFiles проверяет файлы изображений и их метаданные.
func ValidateRecordFiles(root string, records []Record) []error {
	var errs []error

	for _, record := range records {
		imagePath := strings.TrimSpace(record.Image)
		if imagePath == "" {
			continue
		}

		if !isSafeRelativePath(imagePath) {
			continue
		}

		fullPath := filepath.Join(
			root,
			filepath.FromSlash(
				strings.ReplaceAll(imagePath, "\\", "/"),
			),
		)

		data, err := os.ReadFile(fullPath)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: read image %q: %w",
					record.ID,
					imagePath,
					err,
				),
			)
			continue
		}

		if record.SHA256 != "" {
			sum := sha256.Sum256(data)
			actualSHA256 := hex.EncodeToString(sum[:])

			if actualSHA256 != record.SHA256 {
				errs = append(
					errs,
					fmt.Errorf(
						"record %q: image checksum mismatch",
						record.ID,
					),
				)
			}
		}

		file, err := os.Open(fullPath)
		if err != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: open image %q: %w",
					record.ID,
					imagePath,
					err,
				),
			)
			continue
		}

		config, format, decodeErr := image.DecodeConfig(file)

		closeErr := file.Close()

		if decodeErr != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: decode image %q: %w",
					record.ID,
					imagePath,
					decodeErr,
				),
			)
			continue
		}

		if closeErr != nil {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: close image %q: %w",
					record.ID,
					imagePath,
					closeErr,
				),
			)
		}

		if record.Width != config.Width {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: width mismatch: manifest=%d actual=%d",
					record.ID,
					record.Width,
					config.Width,
				),
			)
		}

		if record.Height != config.Height {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: height mismatch: manifest=%d actual=%d",
					record.ID,
					record.Height,
					config.Height,
				),
			)
		}

		if record.Format != "" && record.Format != format {
			errs = append(
				errs,
				fmt.Errorf(
					"record %q: format mismatch: manifest=%q actual=%q",
					record.ID,
					record.Format,
					format,
				),
			)
		}
	}

	sort.SliceStable(
		errs,
		func(i, j int) bool {
			return errs[i].Error() < errs[j].Error()
		},
	)

	return errs
}
