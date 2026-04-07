package paper

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/globals"
)

func (me *Paper) ContourF(img cv.Mat) (gocv.PointVector, cv.Mat) {
	result := gocv.NewPointVector()
	scaled := gocv.NewMat()
	blur := gocv.NewMat()
	noiseBlurSize := 17
	factor := 8
	erodeDillateSize := 9
	d := gocv.NewMat()
	imgThresh := gocv.NewMat()

	gocv.Resize(img.Get(), &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea) // faster detection

	/*
		gMap := cuda.NewGpuMatFromMat(img.Get())
		dstMap := cuda.NewGpuMat()
		gscaled := cuda.NewGpuMat()
		gblur := cuda.NewGpuMat()
		cuda.GammaCorrection(gMap, &dstMap, true)

		cuda.Resize(dstMap, &gscaled, image.Point{dstMap.Cols() / factor, dstMap.Rows() / factor}, 0, 0, gocv.InterpolationArea)

		cuda.CvtColor(gscaled, &gscaled, gocv.ColorBGRAToGray)
		filter := cuda.NewGaussianFilter(gscaled.Type(), gscaled.Type(), image.Pt(23, 23), 30)
		defer filter.Close()
		filter.Apply(gscaled, &gblur)
	*/

	gocv.CvtColor(scaled, &scaled, gocv.ColorBGRAToGray)
	gocv.GaussianBlur(scaled, &blur, image.Pt(noiseBlurSize, noiseBlurSize), 0, 0, gocv.BorderDefault) // reduce noise

	gocv.Dilate(blur, &d, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(erodeDillateSize, erodeDillateSize)))
	gocv.Threshold(d, &imgThresh, 5, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	gocv.Erode(imgThresh, &imgThresh, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(erodeDillateSize, erodeDillateSize)))

	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	newmat := gocv.Zeros(img.Rows(), img.Cols(), img.Type())
	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}
	var points []image.Point
	if maxContourIndex != -1 {
		points = contours.At(maxContourIndex).ToPoints()

		for i := 0; i < len(points); i++ {
			points[i].X *= factor
			points[i].Y *= factor
		}
		result = gocv.NewPointVectorFromPoints(points)
	}

	gocv.DrawContours(&newmat, contours, maxContourIndex, color.RGBA{0, 0, 255, 120}, int(28.0))
	gocv.Resize(newmat, &newmat, image.Point{img.Cols(), img.Rows()}, 0, 0, gocv.InterpolationArea) // faster detection
	/*
		if len(me.debugChannel) == cap(me.debugChannel) {
			old := <-me.debugChannel
			old.Close()
		}
		me.debugChannel <- cv.NewFromMat("paper debugChannel", newmat)
	*/
	cMat := cv.NewFromMat("paper contour image", newmat)
	return result, cMat
}

func (me *Paper) Contour(img cv.Mat) (gocv.PointVector, cv.Mat) {

	result := gocv.NewPointVector()
	scaled := gocv.NewMat()
	blur := gocv.NewMat()
	eroded := gocv.NewMat()
	morphMatSize := 5

	kernel := gocv.Ones(morphMatSize, morphMatSize, gocv.MatTypeCV8U)
	if globals.Globals().PaperFindContourFactor <= 0 {
		globals.Globals().PaperFindContourFactor = 1
	}

	globals.Globals().PaperFindContourFactor = 1

	factor := int(1 / globals.Globals().PaperFindContourFactor)
	noiseBlurSize := globals.Globals().PaperFindContourNoiseBlurSize
	erodeDillateSize := globals.Globals().ErodeDillateSize

	if noiseBlurSize%2 == 0 {
		noiseBlurSize++
	}
	if erodeDillateSize%2 == 0 {
		erodeDillateSize++
	}
	factor = 5

	gocv.Resize(img.Get(), &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea) // faster detection
	gocv.GaussianBlur(scaled, &blur, image.Pt(noiseBlurSize, noiseBlurSize), 0, 0, gocv.BorderDefault)                   // reduce noise
	channels := gocv.Split(blur)
	merged := gocv.NewMat()
	countChannels := len(channels)
	if countChannels > 3 {
		countChannels = 3
	}

	for i := 0; i < countChannels; i++ {
		if ((globals.Globals().FindContourChannelMask) & ((1) << i)) != 0 {
			gocv.MorphologyExWithParams(channels[i], &eroded, gocv.MorphErode, kernel, 3, gocv.BorderDefault)
			if globals.Globals().ShowImage == 501 {
				pImage := gocv.NewMat()
				gocv.CvtColor(eroded, &pImage, gocv.ColorGrayToBGR)
				// me.pipeUIImage(pImage)
				pImage.Close()
			}

			d := gocv.NewMat()
			imgThresh := gocv.NewMat()

			gocv.Dilate(channels[i], &d, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(erodeDillateSize, erodeDillateSize)))
			gocv.Threshold(d, &imgThresh, 140, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
			gocv.Erode(imgThresh, &imgThresh, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(erodeDillateSize, erodeDillateSize)))
			//gocv.Dilate(imgThresh, &imgThresh, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(13, 13)))

			if merged.Empty() {
				merged = imgThresh.Clone()
			} else {
				gocv.Add(merged, imgThresh, &merged)
			}

			channels[i].Close()
			channels[i] = imgThresh.Clone()

			if globals.Globals().ShowImage == 502 {
				pImage := gocv.NewMat()
				gocv.CvtColor(imgThresh, &pImage, gocv.ColorGrayToBGR)
				// me.pipeUIImage(pImage)
				pImage.Close()
			}
			imgThresh.Close()
			d.Close()

		}

	}
	if globals.Globals().ShowImage == 503 {
		pImage := gocv.NewMat()
		gocv.CvtColor(merged, &pImage, gocv.ColorGrayToBGR)
		// me.pipeUIImage(pImage)
		pImage.Close()
	}

	contours := gocv.FindContours(merged, gocv.RetrievalCComp, gocv.ChainApproxSimple)

	newmat := gocv.Zeros(img.Rows()/factor, img.Cols()/factor, img.Type())
	// gocv.PutText(&newmat, "Contouren", image.Pt(100, 300), gocv.FontHersheyPlain, 18.2, color.RGBA{0, 255, 0, 25}, 12)
	// log.Println("newmat", newmat.Channels())

	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}
	var points []image.Point
	if maxContourIndex != -1 {
		points = contours.At(maxContourIndex).ToPoints()

		for i := 0; i < len(points); i++ {
			points[i].X *= factor
			points[i].Y *= factor
		}
		result = gocv.NewPointVectorFromPoints(points)
	}

	gocv.DrawContours(&newmat, contours, maxContourIndex, color.RGBA{0, 0, 255, 120}, int(28.0))
	gocv.Resize(newmat, &newmat, image.Point{img.Cols(), img.Rows()}, 0, 0, gocv.InterpolationArea) // faster detection

	/*
		if len(me.debugChannel) == cap(me.debugChannel) {
			old := <-me.debugChannel
			old.Close()
		}
		me.debugChannel <- cv.NewFromMat("paper debugChannel", newmat)
	*/
	// clean up
	for i := 0; i < len(channels); i++ {
		channels[i].Close()
	}

	scaled.Close()
	blur.Close()
	eroded.Close()
	kernel.Close()

	merged.Close()
	contours.Close()

	/*
		if len(me.contourChannel) == cap(me.contourChannel) {
			o := <-me.contourChannel
			/*
				if !o.Contour.IsNil() {
					o.Contour.Close()
				}
			* /
			o.Mat.Close()
		}
		me.contourChannel <- CC{
			Points: result,
			Mat:    img.Clone("me.contourChannel"),
		}*/
	// img.Close()

	cMat := cv.NewFromMat("paper contour image", newmat)
	return result, cMat

}
