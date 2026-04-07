package tesseract

import (
	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

func (me *Tesseract) fileformatBytes(img *cv.Mat) []byte {
	buffer, err := gocv.IMEncodeWithParams(gocv.PNGFileExt, img.Get(), []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes()
}
