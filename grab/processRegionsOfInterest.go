
package grab

import (
	"image"
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) processRegionsOfInterest(tr structs.TesseractReturnType,img gocv.Mat, useRois []int) structs.TesseractReturnType{
		
	this.pixelScale =  float64(img.Cols()) /  float64(tr.Pagesize.Width)
	this.pixelScaleY =  float64(img.Rows()) /  float64(tr.Pagesize.Height)

	if this.pixelScale==0 {
		this.pixelScale=1
	}
	if this.pixelScaleY==0 {
		this.pixelScaleY=1
	}
	circleSize := int(float64(tr.CircleSize) * this.pixelScale)
	minDist :=float64(tr.CircleMinDistance) * this.pixelScale

	
	marks:=[]structs.CheckMarks{}
	for useRoi := 0; useRoi < len(useRois); useRoi++ {
		if useRoi<len(tr.PageRois) {
			// for pRoiIndex := 0; pRoiIndex < len(tr.PageRois); pRoiIndex++ {
			X := int(float64(tr.PageRois[useRois[useRoi]].X) * this.pixelScale)
			Y := int(float64(tr.PageRois[useRois[useRoi]].Y) * this.pixelScaleY)
			W := int(float64(tr.PageRois[useRois[useRoi]].Width) * this.pixelScale)
			H := int(float64(tr.PageRois[useRois[useRoi]].Height) * this.pixelScaleY)

			rect:=image.Rect( X, Y, X+W, Y+H)
			croppedMat := img.Region(rect)
			
			if !croppedMat.Empty() {
				fMarks:=this.findCircles(croppedMat, circleSize,minDist ,useRois[useRoi] )
				for i := 0; i < len(fMarks); i++ {
					marks = append(marks, fMarks[i])
				}
			}
			croppedMat.Close()
		}
	}
	tr.Marks=marks
	if tr.PageRois[useRois[0]].ExcpectedMarks==len(marks) {
		tr.IsCorrect=true
	}

	return tr
	

}