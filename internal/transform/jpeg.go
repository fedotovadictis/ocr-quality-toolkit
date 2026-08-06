package transform

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
)

// JPEG70 перекодирует изображение в JPEG с quality 70
// и возвращает декодированный результат.
func JPEG70(source image.Image) (image.Image, error) {
	var buffer bytes.Buffer

	if err := jpeg.Encode(
		&buffer,
		source,
		&jpeg.Options{Quality: 70},
	); err != nil {
		return nil, fmt.Errorf("encode JPEG quality 70: %w", err)
	}

	result, err := jpeg.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("decode JPEG quality 70: %w", err)
	}

	return result, nil
}
