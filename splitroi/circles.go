package splitroi

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func circles(img gocv.Mat, pixelScaleX float64, pixelScaleY float64) {
	croppedMatGray := gocv.NewMat()
	gocv.CvtColor(img, &croppedMatGray, gocv.ColorBGRToGray)
	circles := gocv.NewMat()

	dpHoughCircles := 1.0
	thresholdHoughCircles := 90.0
	accumulatorThresholdHoughCircles := 10.0

	innerOverdraw := int(float64(0.01) * math.Min(pixelScaleX, pixelScaleY))
	circleSize := int(float64(5) * math.Min(pixelScaleX, pixelScaleY))
	minDist := (float64(15) * math.Min(pixelScaleX, pixelScaleY))

	gocv.HoughCirclesWithParams(
		croppedMatGray,
		&circles,
		gocv.HoughGradient,
		dpHoughCircles,                   // dp
		minDist,                          //float64(croppedMatGray.Rows()/50), // minDist
		thresholdHoughCircles,            // param1
		accumulatorThresholdHoughCircles, // param2
		circleSize,                       // minRadius
		circleSize,                       // maxRadius
	)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			_color := color.RGBA{220, 220, 250, 0}
			gocv.Circle(&img, image.Pt(x, y), r-innerOverdraw, _color, 40)
		}
	}

}
