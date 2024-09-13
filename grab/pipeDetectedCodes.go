package grab

import (
	"log"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) pipeDetectedCodes(){
	if len(this.detectedCodesChannel)==cap(this.detectedCodesChannel) {
		nonce,_:=<-this.detectedCodesChannel
		if false {
			log.Println("detectedCodesChannel",nonce)
		}

	}
	boxCode:=this.strCurrentBoxBarcode
	if boxCode == "" {
		boxCode="UNKNOWN"
	}
	stackCode:=this.strCurrentStackBarcode
	if stackCode == "" {
		stackCode="UNKNOWN"
	}
	this.detectedCodesChannel <- structs.DetectedCodes{BoxBarcode:boxCode,StackBarcode:stackCode,Barcode:this.lastBarcode}

}