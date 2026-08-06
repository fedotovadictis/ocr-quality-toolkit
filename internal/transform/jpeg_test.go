package transform

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestJPEG70(t *testing.T) {
	src := image.NewRGBA(image.Rect(0, 0, 128, 128))

	// Сделаем изображение неоднотонным, чтобы JPEG действительно сжал его.
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			src.Set(x, y, color.RGBA{
				R: uint8(x * 2),
				G: uint8(y * 2),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}

	got, err := JPEG70(src)
	if err != nil {
		t.Fatalf("JPEG70 returned error: %v", err)
	}

	if got.Bounds() != src.Bounds() {
		t.Fatalf("bounds mismatch: got %v want %v", got.Bounds(), src.Bounds())
	}

	// Проверяем, что результат действительно можно закодировать как JPEG.
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, got, nil); err != nil {
		t.Fatalf("result is not encodable as JPEG: %v", err)
	}

	// Повторная генерация должна быть воспроизводимой.
	got2, err := JPEG70(src)
	if err != nil {
		t.Fatalf("JPEG70 returned error: %v", err)
	}

	var b1, b2 bytes.Buffer

	if err := jpeg.Encode(&b1, got, nil); err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(&b2, got2, nil); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b1.Bytes(), b2.Bytes()) {
		t.Fatal("JPEG70 is not deterministic")
	}
}
