package grab

import (
	"gocv.io/x/gocv"

)


func (this *GrabcameraClass) processPaper( paper gocv.Mat){
	
	//gocv.IMWrite("paper.png",paper)


	area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(this.img.Size()[0]) / float64(this.img.Size()[1])
	// log.Println("extractPaper done %s %f",time.Since(start),area)
	if area > 0.1 {
		
		this.currentState = this.setState("detectedPaper",this.currentState)

		this.processBarcodes(paper)
	}else{

		this.currentState = this.setState("findPaperContourFailed",this.currentState)

	}

}