package corpus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateImageFiles проверяет, что все изображения из manifest существуют на диске.
func ValidateImageFiles(root string, records []Record) error {
	for _, record := range records {
		image := strings.TrimSpace(record.Image)
		if image == "" {
			return fmt.Errorf("record %q: empty image path", record.ID)
		}

		path := image
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}

		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf(
					"record %q: image does not exist: %s",
					record.ID,
					path,
				)
			}

			return fmt.Errorf(
				"record %q: stat image %s: %w",
				record.ID,
				path,
				err,
			)
		}

		if info.IsDir() {
			return fmt.Errorf(
				"record %q: image path is a directory: %s",
				record.ID,
				path,
			)
		}
	}

	return nil
}
