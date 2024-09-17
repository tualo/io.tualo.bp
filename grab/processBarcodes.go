package grab

import (
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) processBarcodes(paper gocv.Mat){


	barcodeTest := this.img.Clone()
	codes := this.findBarcodes(this.scanner,barcodeTest)
	barcodeTest.Close()

	if len(codes) > 0 {

		this.currentState = this.setState("findBarcodes",this.currentState)

		for _, code := range codes {

			if code.Type == "CODE-39" {
				if len(code.Data) >= 3 {
					if code.Data[0:3]=="FC4" {
						if len(this.currentBoxBarcode) == cap(this.currentBoxBarcode) {
							<-this.currentBoxBarcode
						}
						this.currentBoxBarcode <- code.Data
						this.strCurrentBoxBarcode = code.Data

						this.tesseractNeeded = true
						this.doFindCircles = false
						this.checkMarkList = []structs.CheckMarkList{}
						this.debugMarkList = []structs.CheckMarkList{}
						this.currentBallotPaperId = 0
						this.sendNeeded = true

						this.currentState = this.setState("findBoxBarcodes",this.currentState)

					}
					if code.Data[0:3]=="FC3" {
						if len(this.currentStackBarcode) == cap(this.currentStackBarcode) {
							<-this.currentStackBarcode
						}
						this.currentStackBarcode <- code.Data
						this.strCurrentStackBarcode = code.Data

						this.tesseractNeeded = true
						this.doFindCircles = false
						this.checkMarkList = []structs.CheckMarkList{}
						this.debugMarkList = []structs.CheckMarkList{}						
						this.currentBallotPaperId = 0

						this.sendNeeded = true

						this.currentState = this.setState("findStackBarcodes",this.currentState)

					}
				}
			}
			if code.Type == "CODE-128" {

				if len(code.Data)>=5 && code.Data != this.lastBarcode {
					this.lastBarcode = code.Data
					if len(this.ballotBarcode) == cap(this.ballotBarcode) {
						<-this.ballotBarcode
					}
					this.ballotBarcode <- code.Data

					this.tesseractNeeded = true
					this.doFindCircles = false
					this.checkMarkList = []structs.CheckMarkList{}
					this.debugMarkList = []structs.CheckMarkList{}						
					this.currentBallotPaperId = 0
					this.sendNeeded = true
					this.currentState = this.setState("ballotPaperCode",this.currentState)

					this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
					this.pipeDetectedCodes()


					
					

				}

				if this.tesseractNeeded {
					this.processTesseract(paper)
				}

				if this.doFindCircles {
					this.processMarks(paper)
				}else{
					// Suchen der Kreise ist nicht mehr notwendig
					this.currentState = this.setState("doFindCirclesDone",this.currentState)
				}
			}
		}
	}else{
		this.currentState = this.setState("noBarcodeFound",this.currentState)
		this.checkMarkList = []structs.CheckMarkList{}
		this.debugMarkList = []structs.CheckMarkList{}						
		this.currentBallotPaperId = 0
		this.sendNeeded = true
	}

}