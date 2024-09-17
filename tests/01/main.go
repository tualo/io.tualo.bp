package main

import (
	//"fmt"
	"image"
	// "image/color"
	"log"
	"time"
	//"math"
	"os"
	//"sort"
	"gocv.io/x/gocv"
)


func findPaperContourX(img gocv.Mat) gocv.PointVector {

	factor := 8

	scaled := gocv.NewMat()
	// Merge the channels back together
	merged:= gocv.NewMat()

	gocv.Resize(img, &scaled, image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	blur := gocv.NewMat()
	gocv.GaussianBlur(scaled, &blur, image.Pt(83, 83), 1, 1, gocv.BorderDefault)

	
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
	
		gocv.Dilate(imgThresh,&imgThresh,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(33, 33)))

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


	points := contours.At(maxContourIndex).ToPoints()

	for i:=0; i<len(points); i++ {
		points[i].X *= factor
		points[i].Y *= factor
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

func main() {
	filename := os.Args[1]
	window := gocv.NewWindow("fin")
	mat := gocv.IMRead(filename, gocv.IMReadColor)
	start := time.Now()

	findPaperContourX(mat)
	/*
	factor := 8

	scaled := gocv.NewMat()
	gocv.Resize(mat, &scaled, image.Point{mat.Cols() / factor, mat.Rows() / factor}, 0, 0, gocv.InterpolationArea)
	// mat.Close() // Close the original mat to free memory
	// mat=scaled
	// scaled.Close()	



	blur := gocv.NewMat()
	gocv.GaussianBlur(scaled, &blur, image.Pt(83, 83), 1, 1, gocv.BorderDefault)

	channels := gocv.Split(scaled)
	scaled.Close() // Close the scaled mat to free memory

	for i := 0; i < len(channels); i++ {
		eroded := gocv.NewMat()
		kernel := gocv.Ones(5, 5, gocv.MatTypeCV8U)


		 

		gocv.MorphologyExWithParams(channels[i], &eroded, gocv.MorphErode, kernel, 3, gocv.BorderDefault)
	 
	
	
		d := gocv.NewMat()
		imgThresh := gocv.NewMat()


		gocv.Dilate(channels[i],&d,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))
		gocv.Threshold(d, &imgThresh, 100, 255, gocv.ThresholdBinary +gocv.ThresholdOtsu ) 
		gocv.Erode(imgThresh,&imgThresh,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(5, 5)))
	
		gocv.Dilate(imgThresh,&imgThresh,gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(33, 33)))
		 
		channels[i]=imgThresh.Clone()
	}

	// Merge the channels back together
	merged:= gocv.NewMat()
	//gocv.Merge(channels, &merged)

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


	points := contours.At(maxContourIndex).ToPoints()

	for i:=0; i<len(points); i++ {
		points[i].X *= factor
		points[i].Y *= factor
	}

	//contours[maxContourIndex]=gocv.NewPointVectorFromPoints(points)


	if maxContourIndex == -1 {
		 
	}else{
		// Draw the largest contour
		gocv.DrawContours(&mat,contours, maxContourIndex, color.RGBA{255, 255, 0, 0}, 20)
	}
		*/
	log.Println("largest contour %s %f",time.Since(start))
	//maxContour := gocv.NewPointVectorFromPoints( points )



	for {
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}


}
