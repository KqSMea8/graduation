package tools

import (
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
)

func ImageCompress(r io.Reader, jpegW, pngW io.Writer, format constdef.ImageFormat) error {
	var im image.Image
	var err error
	switch format {
	case constdef.Jpeg:
		im, err = JpegDecode(r)
	case constdef.Png:
		im, err = PngDecode(r)
	default:
		logrus.Errorf("ImageFormat invalid: %v", format)
		return fmt.Errorf("invalid ImageFormat: %v", format)
	}
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		JpegCompress(im, jpegW)
	}()
	go func() {
		defer wg.Done()
		PngCompress(im, pngW)
	}()
	wg.Wait()
	return nil
}

// 如果图片不是 jpeg 格式则返回 Error
func JpegCompress(im image.Image, w io.Writer) (err error) {
	if err = jpeg.Encode(w, im, nil); err != nil {
		logrus.Errorf("jpeg.Encode Error: %s", err)
	}
	return err
}

// 如果图片不是 png 格式则返回 Error
func PngCompress(im image.Image, w io.Writer) (err error) {
	if err = png.Encode(w, im); err != nil {
		logrus.Errorf("png.Encode Error: %s", err)
	}
	return err
}

// 如果图片不是 jpeg 格式，则返回 error
func JpegDecode(r io.Reader) (im image.Image, err error) {
	im, err = jpeg.Decode(r)
	if err != nil {
		logrus.Errorf("jpeg.Decode Error: %s", err)
	}
	return
}

func PngDecode(r io.Reader) (im image.Image, err error) {
	im, err = png.Decode(r)
	if err != nil {
		logrus.Errorf("png.Decode Error: %s", err)
	}
	return
}