package grab

import (
	"gocv.io/x/gocv"
	//"fmt"
	"log"
	"image"
	"strings"
	"encoding/json"
	"encoding/base64"
	//api "io.tualo.bp/api"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) processMarks(paper gocv.Mat){


	neededRegions := 0
	for pRoiIndex := 0; pRoiIndex < len(this.lastTesseractResult.PageRois); pRoiIndex++ {
		titles := []string{}

		for i := 0; i < len(this.lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
			titles = append(titles, this.lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
		}
		foundIndex := IndexOf(titles, this.lastTesseractResult.Title)
		if (foundIndex>-1) {
			neededRegions++
		}
	}

	listOfRoiIndexes := []int{}


	foundIndex := -1
	for pRoiIndex := 0; pRoiIndex < len(this.lastTesseractResult.PageRois); pRoiIndex++ {
		titles := []string{}

		for i := 0; i < len(this.lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
			titles = append(titles, this.lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
		}
		fIndex := IndexOf(titles, this.lastTesseractResult.Title)
		if (fIndex>-1) {
			foundIndex = fIndex
			listOfRoiIndexes = append(listOfRoiIndexes,pRoiIndex)
		}
	}


	if len(listOfRoiIndexes) >= 1 {
		if (foundIndex>-1) {

			res := this.processRegionsOfInterest(this.lastTesseractResult,paper,listOfRoiIndexes)

			this.currentBallotPaperId = this.lastTesseractResult.PageRois[listOfRoiIndexes[0]].Types[foundIndex].Id

			if false {
				log.Println("res.Marks",res.Marks,listOfRoiIndexes,res.IsCorrect)
			}
			//	log.Println("res.Id",lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id)
			
			this.currentState = this.setState("ballotPaperMarksAnalysed",this.currentState)


			for i := 0; i < len(res.Marks); i++ {
				if i >= len(this.debugMarkList) {

					offestX := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
					offestY := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)


					this.debugMarkList = append(this.debugMarkList, structs.CheckMarkList{
						Point: image.Point{
							offestX + res.Marks[i].X, 
							offestY + res.Marks[i].Y,
						},
						Pixelsize: res.Marks[i].Radius,
					})
				}
				if res.Marks[i].Checked {
					this.debugMarkList[i].Sum += 1
				}
				this.debugMarkList[i].Count++
				this.debugMarkList[i].AVG = float64(this.debugMarkList[i].Sum) / float64(this.debugMarkList[i].Count)
				this.debugMarkList[i].Checked = this.debugMarkList[i].AVG > this.globals.SumMarksAVG
			}

			if res.IsCorrect {
				this.currentState = this.setState("isCorrect",this.currentState)
				this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
				for i := 0; i < len(res.Marks); i++ {
					if i >= len(this.checkMarkList) {
						offestX := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
						offestY := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)
						this.checkMarkList = append(this.checkMarkList, structs.CheckMarkList{
							Point: image.Point{
								offestX + res.Marks[i].X, 
								offestY + res.Marks[i].Y,
							},
							Pixelsize: res.Marks[i].Radius,
						})
					}
					if res.Marks[i].Checked {
						this.checkMarkList[i].Sum += 1
					}
					this.checkMarkList[i].Count++
					this.checkMarkList[i].AVG = float64(this.checkMarkList[i].Sum) / float64(this.checkMarkList[i].Count)
					this.checkMarkList[i].Checked = this.checkMarkList[i].AVG > this.globals.SumMarksAVG
				}

				log.Println("IsCorrect COUNTER: ",this.checkMarkList[0].Count, this.sendNeeded )

				if len(this.checkMarkList)>0 && this.checkMarkList[0].Count>2 {
					outList:=[]string{}
					for i := 0; i < len(this.checkMarkList); i++ {
						
						if this.checkMarkList[i].Checked {
							outList = append(outList, "X")
						} else {
							outList = append(outList, "O")
						}
					}	
					res.Barcode = this.lastBarcode

					b := new(strings.Builder)
					json.NewEncoder(b).Encode(outList)


					

					image_bytes, _ := gocv.IMEncode(gocv.JPEGFileExt, paper)
					image_base64 := base64.StdEncoding.EncodeToString(image_bytes.GetBytes())
					//fmt.Println("pic",image_base64[0:100])
					image_bytes.Close()

					if this.sendNeeded {
						status := this.sendImageItem(this.strCurrentBoxBarcode,this.strCurrentStackBarcode,res.Barcode,this.lastTesseractResult.PageRois[listOfRoiIndexes[0]].Types[foundIndex].Id,b.String(),"data:image/jpeg;base64,"+image_base64)
						if status {
							this.sendNeeded = false
							this.currentState = this.setState("sendDone",this.currentState)
							this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
						}else{
							this.currentState = this.setState("sendError",this.currentState)
							this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
						}
					}


				}
			}else{
				// log.Println("IsCorrect NO!")
			}
		}
	}
	this.informImage(paper)

}