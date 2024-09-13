package grab

import (
	"gocv.io/x/gocv"
	"fmt"
	"log"
	"image"
	"strings"
	"encoding/json"
	"encoding/base64"
	api "io.tualo.bp/api"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) processBarcodes(paper gocv.Mat){


	barcodeTest := this.img.Clone()
	codes := this.findBarcodes(this.scanner,barcodeTest)
	barcodeTest.Close()

	if len(codes) > 0 {

		this.currentState = this.setState("findBarcodes",this.currentState)

		for _, code := range codes {

			//fmt.Println("code **",code.Type,code.Data)
			if code.Type == "CODE-39" {
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

					this.currentState = this.setState("findStackBarcodes",this.currentState)

				}
			}
			if code.Type == "CODE-128" {

				if code.Data != this.lastBarcode {
					this.lastBarcode = code.Data
					if len(this.ballotBarcode) == cap(this.ballotBarcode) {
						<-this.ballotBarcode
					}
					this.ballotBarcode <- code.Data

					this.tesseractNeeded = true
					this.doFindCircles = false
					this.checkMarkList = []structs.CheckMarkList{}
					this.debugMarkList = []structs.CheckMarkList{}
						this.currentState = this.setState("ballotPaperCode",this.currentState)

					this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
					this.pipeDetectedCodes()
					

				}

				if this.tesseractNeeded {
					this.processTesseract(paper)
				}

				if this.doFindCircles {

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

					/*
					log.Println("listOfRoiIndexes",listOfRoiIndexes)
					log.Println("==================================")
					*/

					if len(listOfRoiIndexes) >= 1 {
					/*for pRoiIndex := 0; pRoiIndex < len(listOfRoiIndexes); pRoiIndex++ {
						titles := []string{}

						for i := 0; i < len(lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
							titles = append(titles, lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
						}
						foundIndex := IndexOf(titles, lastTesseractResult.Title)
							*/
							if (foundIndex>-1) {

							res := this.processRegionsOfInterest(this.lastTesseractResult,paper,listOfRoiIndexes)


							if false {
								log.Println("res.Marks",res.Marks,listOfRoiIndexes,res.IsCorrect)
							}
							//	log.Println("res.Id",lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id)
							
							this.currentState = this.setState("ballotPaperMarksAnalysed",this.currentState)


							for i := 0; i < len(res.Marks); i++ {
								if i >= len(this.debugMarkList) {

									offestX := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
									offestY := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)

									fmt.Println("XYZ",
										offestX + res.Marks[i].X, 
										offestY + res.Marks[i].Y,
									)

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
								// log.Println("IsCorrect",res)
								// lastTesseractResult=res
								this.currentState = this.setState("isCorrect",this.currentState)
								this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)

								
								for i := 0; i < len(res.Marks); i++ {
									if i >= len(this.checkMarkList) {

										offestX := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
										offestY := int(float64(this.lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)

										/*
										fmt.Println("XYZ",
											offestX + res.Marks[i].X, 
											offestY + res.Marks[i].Y,
										)
										*/

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

								if len(this.checkMarkList)>0 && this.checkMarkList[0].Count>6 {
									//

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

									res,err := api.SendReading(
										this.strCurrentBoxBarcode,
										this.strCurrentStackBarcode,
										res.Barcode,
										this.lastTesseractResult.PageRois[listOfRoiIndexes[0]].Types[foundIndex].Id,
										b.String(),
										"data:image/jpeg;base64,"+image_base64,
									)
									
									if err != nil {
										log.Println("SendReading ERROR",err)
										
										this.currentState = this.setState("sendError",this.currentState)
										this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
									}else{
										// log.Println(">>>>",res.Msg)
										if res.Success {
											this.doFindCircles = false


											
											this.currentState = this.setState("sendDone",this.currentState)
											this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
										}else{
											log.Println("SendReading ERROR",res.Msg)
										}
									}

								}
							}else{
								// log.Println("IsCorrect NO!")
							}
							}
					}
				}else{
					// Suchen der Kreise ist nicht mehr notwendig
					this.currentState = this.setState("doFindCirclesDone",this.currentState)
					//log.Println("doFindCirclesDone",doFindCircles)
					// green = 50
					// red = 0
					// blue = 50
				}
				//log.Println("code use tesseract",code.Data,tesseractNeeded,lastTesseractResult)
			}
		}
		// gocv.IMWrite("paper.png",paper)
	}else{
		this.currentState = this.setState("noBarcodeFound",this.currentState)
		this.checkMarkList = []structs.CheckMarkList{}
		this.debugMarkList = []structs.CheckMarkList{}
	}

}