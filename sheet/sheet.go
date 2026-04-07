package sheet

import (
	"image"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/structs"
)

type Sheet struct {
	mat                 *gocv.Mat
	width_in_mm         int
	height_in_mm        int
	detectBallotpaperId int
	detectPaginationId  int
}

func (me *Sheet) SetImage(mat *gocv.Mat) {
	me.mat = mat
}

func (me *Sheet) GetImage() *gocv.Mat {
	return me.mat
}

func (me *Sheet) Measurements(width int, height int) {
	me.width_in_mm = width
	me.height_in_mm = height
}

func (me *Sheet) Rescale() structs.Scale {
	var goalRatioX float64
	var goalRatioY float64

	goalRatioX = 1.0
	goalRatioY = 1.0
	originalFactor := (float64(me.width_in_mm) / float64(me.height_in_mm))
	imageFactor := (float64(me.mat.Cols()) / float64(me.mat.Rows()))

	if imageFactor > originalFactor {
		goalRatioY = (imageFactor / originalFactor)
	} else {
		goalRatioX = 1 / (imageFactor / originalFactor)
	}

	gocv.Resize(
		*me.mat,
		me.mat,
		image.Point{
			int(float64(me.mat.Cols()) * goalRatioX),
			int(float64(me.mat.Rows()) * goalRatioY),
		},
		0,
		0,
		gocv.InterpolationArea,
	)

	scale := structs.Scale{
		PixelScaleX: float64(me.mat.Cols()) / float64(me.width_in_mm),
		PixelScaleY: float64(me.mat.Rows()) / float64(me.height_in_mm),
	}
	return scale
}
