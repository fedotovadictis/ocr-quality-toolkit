package imageinfo

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func Read(path string) (format string, width int, height int, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, 0, err
	}
	defer file.Close()

	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return "", 0, 0, err
	}
	return format, config.Width, config.Height, nil
}
