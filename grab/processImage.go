package grab

import (
	"log"
	// "gocv.io/x/gocv"
	// "image"
	// "image/color"
	// "time"
	"github.com/bieber/barcode"
	// "fmt"
	// api "io.tualo.bp/api"
	// "strings"
	// "encoding/json"
	// "encoding/base64"
	structs "io.tualo.bp/structs"
)








func (this *GrabcameraClass) processImage(){
	this.scanner = barcode.NewScanner()
	this.scanner.SetEnabledAll(false)
	this.scanner.SetEnabledSymbology(barcode.Code39,true)
	this.scanner.SetEnabledSymbology(barcode.Code128,true)
	// log.Println("processImage starting ")
	this.tesseractNeeded = true
	this.doFindCircles = false

	this.lastTesseractResult = structs.TesseractReturnType{}
	this.checkMarkList = []structs.CheckMarkList{}
	this.debugMarkList = []structs.CheckMarkList{}
	this.currentState = structs.ImageProcessorState{};
	this.currentState = this.setState("default",this.currentState)
	this.sendNeeded = true
	for {
		if !this.runVideo {
			break
		}

		this.img,this.paperChannelImageOK = <-this.paperChannelImage
		if this.paperChannelImageOK {

			if !this.img.Empty() {
				this.processImageChannelData()
			}else{
				log.Println("img empty")
			}
			this.img.Close()


			// WRONG POSITION
			/*
			log.Println("img ***")
			if len(this.imageChannelPaper)==cap(this.imageChannelPaper) {
				mat,_:=<-this.imageChannelPaper
				mat.Close()
			}

			cloned := this.img.Clone()
			gocv.Resize(cloned, &cloned, image.Point{}, 0.3, 0.3 , gocv.InterpolationLinear)
			this.imageChannelPaper <- cloned
			this.img.Close()
			*/
		}
			

		
		
	}
}