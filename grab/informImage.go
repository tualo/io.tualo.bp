package grab

import (
	"log"
	"gocv.io/x/gocv"
	"encoding/base64"
)


func (this *GrabcameraClass) informImage(paper gocv.Mat){
	if len(this.escapedImage)>0 {
		escapedImage,escapedImageOk := <-this.escapedImage
		if escapedImageOk {
			if true {
				log.Println("escaped image found",escapedImage)
			}
			image_bytes, _ := gocv.IMEncode(gocv.JPEGFileExt, paper)
			image_base64 := base64.StdEncoding.EncodeToString(image_bytes.GetBytes())
			image_bytes.Close()
			this.sendImageItem(
				this.strCurrentBoxBarcode,
				this.strCurrentStackBarcode,
				this.lastBarcode,
				0/**/, 
				"[]",
				"data:image/jpeg;base64,"+image_base64,
			)

			this.currentState = this.setState("escaped",this.currentState)
			this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)

		}
	}
}