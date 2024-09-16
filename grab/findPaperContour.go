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
	var thresh float32

	factor:=2
	imgGray := gocv.NewMat()
	imgBlur := gocv.NewMat()
	imgThresh := gocv.NewMat()

	scaled := gocv.NewMat()
	gocv.Resize(img, &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	gocv.CvtColor(scaled, &imgGray, gocv.ColorBGRToGray)
	  
	
	thresh = calculateBestThresh(imgGray)
	gocv.Threshold(imgGray, &imgGray, thresh , 255 , gocv.ThresholdBinary    ) 

	contours := gocv.FindContours(imgGray, gocv.RetrievalCComp, gocv.ChainApproxSimple)
	  
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

	/*
	mean := gocv.NewMat()
	stddev := gocv.NewMat()

	gocv.MeanStdDev(channels[0], &mean, &stddev)
	fmt.Println("1) mean: ", mean.GetFloatAt(0, 0), "stddev: ", stddev.GetFloatAt(0, 0))
	gocv.MeanStdDev(channels[1], &mean, &stddev)
	fmt.Println("2) mean: ", mean.GetFloatAt(0, 0), "stddev: ", stddev.GetFloatAt(0, 0))
	gocv.MeanStdDev(channels[2], &mean, &stddev)
	fmt.Println("3) mean: ", mean , "stddev: ", stddev )
	*/


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


	/*
	var s float32
	s =0
	bytes := channels[0].ToBytes()
	for _, b := range bytes {
		s+=float32(b)
	}
	fmt.Println("sum: ", s, "len: ", len(bytes), "mean: ", s/float32(len(bytes)))
	*/



	 /*
	gocv.Merge([]gocv.Mat{channels[0], channels[1], channels[2]}, &chanIMG)
	gocv.IMWrite("chanIMG.jpg",chanIMG)
	*/

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
