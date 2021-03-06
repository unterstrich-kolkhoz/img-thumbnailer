package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestImageResizeSizes(t *testing.T) {
	// original size, resize params, expected new size
	sizes := [][]int{
		{100, 200, 10, 20, 10, 20},
		{100, 200, 10, 0, 10, 20},
		{100, 200, 0, 20, 10, 20},
		{100, 200, 5, 20, 5, 20},
		{100, 200, 500, 0, 500, 1000},
		{200, 50, 100, 0, 100, 25},
		{200, 50, 100, 100, 100, 100},
	}

	os.Mkdir("test_tmp", 0700)
	defer os.RemoveAll("test_tmp")

	for _, options := range sizes {
		origW, origH := options[0], options[1]
		paramW, paramH := options[2], options[3]
		expectedW, expectedH := options[4], options[5]

		// create a new image
		origImg := image.NewGray(image.Rect(0, 0, origW, origH))

		// save it as jpeg
		src, _ := os.Create("test_tmp/src.jpg")
		jpeg.Encode(src, origImg, nil)
		src.Close()

		path, err := resize("test_tmp/src.jpg", "jpg", int(paramW), int(paramH), 95)

		if err != nil {
			t.FailNow()
		}

		dst, _ := os.Open(path)
		img, _ := jpeg.Decode(dst)
		dst.Close()
		if img.Bounds().Max.X != expectedW || img.Bounds().Max.Y != expectedH {
			t.FailNow()
		}
	}
}

func testFormatConversion(t *testing.T) {
	os.Mkdir("test_tmp", 0700)
	defer os.RemoveAll("test_tmp")

	// create a new image
	origImg := image.NewGray(image.Rect(0, 0, 100, 100))

	// save it as png
	src, _ := os.Create("test_tmp/src.png")
	png.Encode(src, origImg)
	src.Close()

	// convert it to jpeg
	path, err := resize("test_tmp/src.png", "jpg", 100, 100, 95)

	if err != nil {
		t.FailNow()
	}

	// reconvert to image.Image
	dst, _ := os.Open(path)
	img, _ := jpeg.Decode(dst)
	dst.Close()

	// check that both images are the same
	if origImg != img {
		t.FailNow()
	}
}

func TestImageResizeCropping(t *testing.T) {
	os.Mkdir("test_tmp", 0700)
	defer os.RemoveAll("test_tmp")

	white := color.Gray{255}
	gray := color.Gray{100}
	origImg := image.NewGray(image.Rect(0, 0, 100, 100))
	origImg.Set(24, 24, white)
	origImg.Set(100-25, 100-25, white)
	origImg.Set(100-25, 100-25, white)
	origImg.Set(25, 25, gray)
	origImg.Set(100-26, 100-26, gray)

	src, _ := os.Create("test_tmp/src.png")
	png.Encode(src, origImg)
	src.Close()

	path, err := resize("test_tmp/src.png", "jpg", 50, 100, 25)

	if err != nil {
		t.FailNow()
	}

	dst, _ := os.Open(path)
	img, _ := png.Decode(dst)
	dst.Close()

	// check cropping
	// check that white pixels have been cropped away
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if colorsEqual(img.At(x, y), white) {
				t.FailNow()
			}
		}
	}

	// check that gray pixels are still in the image
	if !colorsEqual(img.At(0, 25), gray) {
		t.FailNow()
	}

	if !colorsEqual(img.At(49, 100-26), gray) {
		t.FailNow()
	}
}

func colorsEqual(first, second color.Color) bool {
	r1, g1, b1, a1 := first.RGBA()
	r2, g2, b2, a2 := second.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
