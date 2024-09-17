package grab

import (
	"image"
	"gocv.io/x/gocv"
)


func calculateBestThresh(img gocv.Mat) float32 {
	x:=img.Cols() - 1
	max_val := 0
	for i:=0; i<img.Rows(); i++ {
		for j:=0; j<1; j++ {
			mx:=img.GetUCharAt(  i, x+j)
			if int(mx) > (max_val) {
				max_val = int(mx)
			}

		}
	}
	return float32(max_val)*0.9
}



func findPaperContour(img gocv.Mat) gocv.PointVector {

	factor := 12
	//vig:=gocv.IMRead("vig.png",gocv.IMReadColor)

	scaled := gocv.NewMat()
	//X := gocv.NewMat()
	// Merge the channels back together
	merged:= gocv.NewMat()

	//gocv.Add(img, vig, &X)

	gocv.Resize(img, &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	blur := gocv.NewMat()
	gocv.GaussianBlur(scaled, &blur, image.Pt(3, 3), 1, 1, gocv.BorderDefault)

	
	channels := gocv.Split(scaled)
	countChannels := len(channels)
	if countChannels > 3 {
		countChannels = 3
	}
	for i := 0; i < countChannels; i++ {
		eroded := gocv.NewMat()
		kernel := gocv.Ones(5, 5, gocv.MatTypeCV8U)
		gocv.MorphologyExWithParams(channels[i], &eroded, gocv.MorphErode, kernel, 3, gocv.BorderDefault)
		d := gocv.NewMat()
		imgThresh := gocv.NewMat()


		gocv.Dilate(channels[i],&d,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))
		gocv.Threshold(d, &imgThresh, 100, 255, gocv.ThresholdBinary +gocv.ThresholdOtsu ) 
		gocv.Erode(imgThresh,&imgThresh,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(5, 5)))
	
		gocv.Dilate(imgThresh,&imgThresh,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))

		channels[i].Close()
		channels[i]=imgThresh.Clone()

		eroded.Close()
		kernel.Close()
		d.Close()
		imgThresh.Close()

	}


	gocv.Add(channels[0], channels[1], &merged)
	gocv.Add(merged, channels[2], &merged)
	gocv.Threshold(merged, &merged, 100, 255, gocv.ThresholdBinary +gocv.ThresholdOtsu ) 

	
	contours := gocv.FindContours(merged, gocv.RetrievalCComp, gocv.ChainApproxSimple)
	  
	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}

	if maxContourIndex != -1 {
		points := contours.At(maxContourIndex).ToPoints()

		for i:=0; i<len(points); i++ {
			points[i].X *= factor
			points[i].Y *= factor
		}
	}

	//contours[maxContourIndex]=gocv.NewPointVectorFromPoints(points)

	for i:=0; i<len(channels); i++ {
		channels[i].Close()
	}
	scaled.Close()
	blur.Close()
	merged.Close()
	contours.Close()


	if maxContourIndex == -1 {
		
		return gocv.NewPointVector()
		
	}
	return gocv.NewPointVectorFromPoints(points)



}



func findPaperContour16092024(img gocv.Mat) gocv.PointVector {
	var thresh float32

	factor:=4
	imgGray := gocv.NewMat()
	imgBlur := gocv.NewMat()
	imgThresh := gocv.NewMat()

	scaled := gocv.NewMat()
	gocv.Resize(img, &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationLanczos4)
	// gocv.GaussianBlur(scaled, &imgBlur, image.Point{85, 85}, 0, 0, gocv.BorderDefault)
	gocv.CvtColor(scaled, &imgGray, gocv.ColorBGRToGray)
	  
	
	thresh = calculateBestThresh(imgGray)
	gocv.Threshold(imgGray, &imgThresh, thresh , 255 , gocv.ThresholdBinary  /* +gocv.ThresholdOtsu */ ) 
	//gocv.GaussianBlur(imgThresh, &imgBlur, image.Point{5, 5}, 0, 0, gocv.BorderDefault)

	//contours := gocv.FindContours(imgGray, gocv.RetrievalCComp, gocv.ChainApproxSimple)
	contours := gocv.FindContours(imgThresh, gocv.RetrievalList, gocv.ChainApproxTC89KCOS)
	  
	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}
	if maxContourIndex == -1 {
		imgGray.Close()
		imgBlur.Close()
		imgThresh.Close()
		return gocv.NewPointVector()
	}
	// gocv.Resize(contours, &contours, image.Point{contours.Cols() / factor, contours.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	points := contours.At(maxContourIndex).ToPoints()

	for i:=0; i<len(points); i++ {
		points[i].X *= factor
		points[i].Y *= factor
	}

	maxContour := gocv.NewPointVectorFromPoints( points )
	imgGray.Close()
	imgBlur.Close()
	imgThresh.Close()
	contours.Close()
	scaled.Close()
	return maxContour

}

func findPaperContourChannels(img gocv.Mat) gocv.PointVector {
	var thresh float32
	imgGray := gocv.NewMat()
//	chanIMG := gocv.NewMat()
	chanIMGM:= gocv.NewMat()
	/*
	gocv.IMWrite("imgThresh_orig.jpg", img)
	gocv.CvtColor(img, &imgGray, gocv.ColorBGRToGray)
*/
	factor:=2

	scaled := gocv.NewMat()
	gocv.Resize(img, &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	defer scaled.Close()	


	channels := gocv.Split(scaled)
 

	thresh = calculateBestThresh(channels[0])
	gocv.GaussianBlur(channels[0], &channels[0], image.Point{15, 15}, 0, 0, gocv.BorderDefault)
	//gocv.IMWrite("channels0a.jpg", channels[0])
	gocv.Threshold(channels[0], &channels[0], thresh, 255, gocv.ThresholdBinary  ) 
	//gocv.IMWrite("channels0.jpg", channels[0])

	thresh = calculateBestThresh(channels[1])
	gocv.GaussianBlur(channels[1], &channels[1], image.Point{15, 15}, 0, 0, gocv.BorderDefault)
	//gocv.IMWrite("channels1a.jpg", channels[1])
	gocv.Threshold(channels[1], &channels[1],thresh, 255, gocv.ThresholdBinary  ) 
	//gocv.IMWrite("channels1.jpg", channels[1])


	thresh = calculateBestThresh(channels[2])
	gocv.GaussianBlur(channels[2], &channels[2], image.Point{25, 25}, 0, 0, gocv.BorderDefault)
	//gocv.IMWrite("channels2a.jpg", channels[2])

	gocv.Threshold(channels[2], &channels[2], thresh , 255 , gocv.ThresholdBinary    ) 
	//gocv.IMWrite("channels2.jpg", channels[2])



	gocv.Add(channels[0], channels[1], &chanIMGM)
	gocv.Add(chanIMGM, channels[2], &chanIMGM)
	//gocv.IMWrite("chanIMGmult.jpg",chanIMGM)

	//gocv.CvtColor(chanIMGM, &imgGray, gocv.ColorBGRToGray)
	//gocv.CvtColor(chanIMGM, &imgGray, gocv.ColorBGRToGray)
	
	imgGray = chanIMGM.Clone()
	
	chanIMGM.Close()
	//chanIMG.Close()


		
	for _, channel := range channels {
		channel.Close()
	}

	// gocv.Resize(imgGray, &imgGray, image.Point{imgGray.Cols() / 8, imgGray.Rows() / 8}, 0, 0, gocv.InterpolationArea)


	imgBlur := gocv.NewMat()
	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{55, 55}, 0, 0, gocv.BorderDefault)
	imgThresh := gocv.NewMat()
	// gocv.Threshold(imgBlur, &imgThresh, 140, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	// gocv.Threshold(imgBlur, &imgThresh, 4, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)
	// gocv.IMWrite("imgThresh0.jpg", imgBlur)
	gocv.Threshold(imgBlur, &imgThresh, 15, 255, gocv.ThresholdBinary +gocv.ThresholdOtsu ) 
	// gocv.IMWrite("imgThresh1.jpg", imgThresh)
	
/*

	gocv.GaussianBlur(imgThresh, &imgBlur, image.Point{15, 15}, 0, 0, gocv.BorderDefault)
	gocv.Threshold(imgBlur, &imgThresh, 140, 255, gocv.ThresholdBinary+gocv.ThresholdOtsu)

*/

	// gocv.IMWrite("imgThresh2.jpg", imgThresh)
	contours := gocv.FindContours(imgThresh, gocv.RetrievalCComp, gocv.ChainApproxSimple)
	  
	maxArea := 0.0
	maxContourIndex := -1
	for i := 0; i < contours.Size(); i++ {
		contourArea := gocv.ContourArea(contours.At(i))
		if contourArea > maxArea {
			maxArea = contourArea
			maxContourIndex = i
		}
	}
	if maxContourIndex == -1 {
		imgGray.Close()
		imgBlur.Close()
		imgThresh.Close()
		return gocv.NewPointVector()
	}
	// gocv.Resize(contours, &contours, image.Point{contours.Cols() / factor, contours.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	points := contours.At(maxContourIndex).ToPoints()

	for i:=0; i<len(points); i++ {
		points[i].X *= factor
		points[i].Y *= factor
	}

	maxContour := gocv.NewPointVectorFromPoints( points )
	imgGray.Close()
	imgBlur.Close()
	imgThresh.Close()
	contours.Close()
	return maxContour
}
