package blurry

import (
	"math"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/args"
	"tualo.de/deep-test/cv"
)

type Blurry struct {
	blurryThreshold float64
}

func Static() *Blurry {
	return &Blurry{
		blurryThreshold: float64(*args.Parser().BlurryThreshold),
	}
}

func (me *Blurry) variance_of_laplacian(mat *cv.Mat) float64 {
	image := mat.Clone("variance_of_laplacian")
	if image.GetPointer().Channels() > 1 {
		gocv.CvtColor(image.Get(), image.GetPointer(), gocv.ColorBGRAToGray)
	}
	dst := gocv.NewMat()
	gocv.Laplacian(image.Get(), &dst, gocv.MatTypeCV16S, 1, 1, 0, gocv.BorderDefault)

	mean := gocv.NewMat()
	stdDev := gocv.NewMat()
	gocv.MeanStdDev(dst, &mean, &stdDev)
	flts, _ := stdDev.DataPtrFloat64()
	//log.Println(flts, stdDev.ToBytes()[0])
	variance := math.Pow(flts[0], 2)

	dst.Close()
	image.Close()
	stdDev.Close()
	mean.Close()
	return variance
}

func (me *Blurry) Detect(mat *cv.Mat) (bool, float64) {
	// start := time.Now()
	variance := me.variance_of_laplacian(mat)
	// log.Println("variance", variance, time.Since(start))

	if variance < me.blurryThreshold {
		return true, variance
	} else {
		return false, variance
	}
}
