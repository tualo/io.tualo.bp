package grab

import (
	"gocv.io/x/gocv"
	"fmt"
)

func (this *GrabcameraClass) processImageChannelData(){

	gocv.CvtColor(this.img, &this.img, gocv.ColorBGRToBGRA)
	this.currentState = this.setState("findPaperContour",this.currentState)
	this.contour = findPaperContour(this.img)
	if this.contour.Size() == 0 {
		// contour.Close()
		this.pipeUIImage(this.img)
	}else{
		approx := gocv.ApproxPolyDP(this.contour, 0.02*gocv.ArcLength(this.contour, true), true)
		if !(approx.Size() >= 4 &&  approx.Size() <= 7) {
			this.pipeUIImage(this.img)
			
		}else{

			cornerPoints := getCornerPoints(this.contour)
			fmt.Println("cornerPoints")
			fmt.Println(cornerPoints)
			fmt.Println("====================================")
			topLeftCorner := cornerPoints["topLeftCorner"]
			bottomRightCorner := cornerPoints["bottomRightCorner"]

			paper := gocv.NewMat()
			paper,this.invM = extractPaper(this.img, this.contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
			if paper.Empty() {
				this.pipeUIImage(this.img)
			}else{
				this.playGround = paper.Clone()
				this.processPaper(paper);
				this.drawBackResults()
				this.pipeUIImage(this.img)
				this.playGround.Close()
			}
			paper.Close()
			this.invM.Close()

							


		}
		approx.Close()
	}
	this.contour.Close()


}