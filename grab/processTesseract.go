package grab

import (
	"log"
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)


func (this *GrabcameraClass) processTesseract(paper gocv.Mat){

	result := this.tesseract(paper,this.currentOCRChannel)
	if len(result.PageRois)>0 {
		this.tesseractNeeded = false
		this.lastTesseractResult = result
		this.doFindCircles = true
		this.checkMarkList = []structs.CheckMarkList{}
		this.debugMarkList = []structs.CheckMarkList{}
		this.currentState = this.setState("ballotPaperDetected",this.currentState)
	}else{
		this.currentState = this.setState("ballotPaperNotDetected",this.currentState)
		log.Println("tesseract no bp found")
	}

}