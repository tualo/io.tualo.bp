package grab

import (
	"log"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"time"
	"github.com/bieber/barcode"
	"fmt"
	api "io.tualo.bp/api"
	"strings"
	"encoding/json"
	"encoding/base64"
	structs "io.tualo.bp/structs"
)

/*
func NewLuminanceSourceFromImage(img image.Image) LuminanceSource {

func (this *GrabcameraClass) findBarcodes(scanner *barcode.ImageScanner, img gocv.Mat)[]structs.BarcodeSymbol{

	luminance := make([]byte, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			luminance[y*width+x] = byte((r + g + b) / 3)
		}
	}

	&GoImageLuminanceSource{&RGBLuminanceSource{
		LuminanceSourceBase{width, height},
		luminance,
		width,
		height,
		0,
		0,
	}}

	NewBinaryBitmap(NewHybridBinarizer(src))

	dec := oned.NewCode128Reader()
    i, err := gozxing.NewBinaryBitmapFromImage(bitmapimg)
    if err != nil {
        return content, err
    }
    r, err := dec.DecodeWithoutHints(i)
    if err != nil {
        return content, err
    }
}
*/
func (this *GrabcameraClass) findBarcodes(scanner *barcode.ImageScanner, img gocv.Mat)[]structs.BarcodeSymbol{
	syms := []structs.BarcodeSymbol{}
	if img.Empty() {
		return syms
	}
	barcodeScale := 1

	smaller:=gocv.NewMat()
	gocv.CvtColor(img, &smaller, gocv.ColorBGRToGray)
	if smaller.Cols() > 800 {
		gocv.GaussianBlur(smaller, &smaller, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
		gocv.Resize(smaller, &smaller, image.Point{smaller.Cols() / barcodeScale, smaller.Rows() / barcodeScale}, 0, 0, gocv.InterpolationArea)
	}
	if false {
		log.Println("barcodeScale",barcodeScale,smaller.Cols())
	}
	symbols, err := scanner.ScanMat(&smaller)
	if err != nil {
		panic(err)
	}
	
	
	/*
	log.Println("findBarcodes",len(symbols))
	if len(symbols) == 0 {
		gocv.IMWrite("noBarcode.png",img)
	}else{
		gocv.IMWrite("barcode.png",img)
	
	}
	*/
	
	
	for _, s := range symbols {
		syms = append(syms,structs.BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
		if false {
			log.Println("BarcodeSymbol",s.Type.Name(),s.Data,s.Quality,s.Boundary)
		}
	}
	smaller.Close()
	return syms
}



func (this *GrabcameraClass) processRegionsOfInterest(tr structs.TesseractReturnType,img gocv.Mat, useRois []int) structs.TesseractReturnType{
	

	//log.Println("processRegionsOfInterest Ratio ",float64(img.Cols()) / float64(img.Rows()),  float64(tr.Pagesize.Width) / float64(tr.Pagesize.Height) )
						
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

	/*
	if false {
		log.Println("processRegionsOfInterest",tr.PageRois[useRoi].X, 
			tr.PageRois[useRoi].Y, tr.PageRois[useRoi].Width, 
			tr.PageRois[useRoi].Height, tr.PageRois[useRoi].ExcpectedMarks, 
			tr.PageRois[useRoi].Types[0].Title,
			"pixelScale",this.pixelScale,
			"pixelScaleY",this.pixelScaleY,
			"circleSize",circleSize,
			"minDist",minDist,
		)
	}
*/
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


				/*
				
				*/
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

func (this *GrabcameraClass) setState(name string,oldState structs.ImageProcessorState) structs.ImageProcessorState{
	state := structs.ImageProcessorState{}
	if false {
		log.Println("setState",name,oldState.Name);
	}

	if oldState.Name == "sendDone" && name != "ballotPaperCode"  && name != "noBarcodeFound" {
		return oldState
	} 

	state.Name = name

	if len(this.currentStateChannel)==cap(this.currentStateChannel) {
		txt,_:=<-this.currentStateChannel
		if false {
			log.Println("setState",txt)
		}
	}
	this.currentStateChannel <- name



	if (name == "default") {
		state.Red = 0
		state.Green = 0
		state.Blue = 0
	}
	if (name == "findPaperContour") {
		state.Red = 110
		state.Green = 110
		state.Blue = 110
	}

	if (name == "findPaperContourFailed") {
		state.Red = 155
		state.Green = 110
		state.Blue = 110
	}

	if (name == "detectedPaper") {
		state.Red = 200
		state.Green = 110
		state.Blue = 110
	}

	if (name == "findBarcodes") {
		state.Red = 110
		state.Green = 200
		state.Blue = 110
	}

	if (name == "noBarcodeFound") {
		state.Red = 255
		state.Green = 0
		state.Blue = 0
	}

	if (name == "findBoxBarcodes") {
		state.Red = 110
		state.Green = 110
		state.Blue = 200
	}

	if (name == "findStackBarcodes") {
		state.Red = 200
		state.Green = 200
		state.Blue = 110
	}

	if (name == "ballotPaperCode") {
		state.Red = 0
		state.Green = 0
		state.Blue = 255
	}

	if (name == "ballotPaperDetected") {
		state.Red = 255
		state.Green = 255
		state.Blue = 255
	}

	if (name == "ballotPaperNotDetected") {
		state.Red = 100
		state.Green = 255
		state.Blue = 100
	}

	if (name == "ballotPaperMarksAnalysed") {
		state.Red = 255
		state.Green = 100
		state.Blue = 255
	}

	if (name == "doFindCirclesDone") {
		state.Red = 155
		state.Green = 50
		state.Blue = 155
	}

	if (name == "sendError") {
		state.Red = 255
		state.Green = 120
		state.Blue = 120
	}

	if (name == "sendDone") {
		state.Red = 0
		state.Green = 255
		state.Blue = 0
	}

	if (name == "isCorrect") {
		state.Red = 10
		state.Green = 55
		state.Blue = 10
	}


	return state
}

func (this *GrabcameraClass) setHistoryItem(barcode string,boxcode string,stackcode string, currentState structs.ImageProcessorState){
	histItem := structs.HistoryListItem{
		Barcode: barcode,
		BoxBarcode: boxcode,
		StackBarcode: stackcode,
		State: currentState.Name,
		StateColor: color.RGBA{uint8(currentState.Red),uint8(currentState.Green),uint8(currentState.Blue),120},
	}

	if len(this.listItemChannel)==cap(this.listItemChannel) {
		oldItem,_:=<-this.listItemChannel
		if false {
			log.Println("setState oldItem",oldItem)
		}
	}
	this.listItemChannel <- histItem
}

func (this *GrabcameraClass) processImage(){
	scanner := barcode.NewScanner()
	scanner.SetEnabledAll(false)
	scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)
	// log.Println("processImage starting ")
	tesseractNeeded := true
	lastTesseractResult := structs.TesseractReturnType{}
	doFindCircles := false
	checkMarkList := []structs.CheckMarkList{}
	debugMarkList := []structs.CheckMarkList{}
	
	

	currentState := structs.ImageProcessorState{};



	currentState = this.setState("default",currentState)
	for {
		if !this.runVideo {
			break
		}
		start:=time.Now()
		if false {
			log.Println("processImage ************")
		}
		//for range grabVideoCameraTicker.C {	
		img,ok := <-this.paperChannelImage
		if ok {
			if false {
				log.Println("got image",ok,img.Size(),len(this.paperChannelImage))
			}
//			gocv.Resize(img, &img, image.Point{}, 0.75, 0.75 , gocv.InterpolationLinear)
			if !img.Empty() {
				gocv.CvtColor(img, &img, gocv.ColorBGRToBGRA)
				currentState = this.setState("findPaperContour",currentState)
				barcodeTest := img.Clone()
				contour := findPaperContour(img)
				if contour.Size() == 0 {
					contour.Close()
				}else{
					approx := gocv.ApproxPolyDP(contour, 0.02*gocv.ArcLength(contour, true), true)
					if !(approx.Size() >= 4 &&  approx.Size() <= 7) {
						if false {
							log.Println("findPaperContour done %s %v",time.Since(start),approx.Size())
						}
						approx.Close()
						contour.Close()
					}else{
						approx.Close()

						cornerPoints := getCornerPoints(contour)
						topLeftCorner := cornerPoints["topLeftCorner"]
						bottomRightCorner := cornerPoints["bottomRightCorner"]
						if false {
							log.Printf("template: %d %d",  bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y )
						}

						paper,invM := extractPaper(img, contour, bottomRightCorner.X-topLeftCorner.X, bottomRightCorner.Y-topLeftCorner.Y, cornerPoints)
						
						if paper.Empty() {
							if false {
								log.Printf("paper empty")
							}
							contour.Close()

							cloned := img.Clone()
							this.imageChannelPaper <- cloned
							img.Close()
							
							img.Close()
							paper.Close()
							invM.Close()
							continue
						}
						playGround := paper.Clone()
						// gocv.IMWrite("playGround.png",playGround)


						area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(img.Size()[0]) / float64(img.Size()[1])
						// log.Println("extractPaper done %s %f",time.Since(start),area)
						if area > 0.1 {
							
							currentState = this.setState("detectedPaper",currentState)

							codes := this.findBarcodes(scanner,barcodeTest)
							barcodeTest.Close()
							if false {
								log.Println("findBarcodes done %s %v",time.Since(start),codes)
							}
							if len(codes) > 0 {

								currentState = this.setState("findBarcodes",currentState)

								for _, code := range codes {

									//fmt.Println("code **",code.Type,code.Data)
									if code.Type == "CODE-39" {
										if code.Data[0:3]=="FC4" {
											if len(this.currentBoxBarcode) == cap(this.currentBoxBarcode) {
												<-this.currentBoxBarcode
											}
											this.currentBoxBarcode <- code.Data
											this.strCurrentBoxBarcode = code.Data

											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}
											debugMarkList = []structs.CheckMarkList{}

											currentState = this.setState("findBoxBarcodes",currentState)

										}
										if code.Data[0:3]=="FC3" {
											if len(this.currentStackBarcode) == cap(this.currentStackBarcode) {
												<-this.currentStackBarcode
											}
											this.currentStackBarcode <- code.Data
											this.strCurrentStackBarcode = code.Data

											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}
											debugMarkList = []structs.CheckMarkList{}

											currentState = this.setState("findStackBarcodes",currentState)

										}
									}
									if code.Type == "CODE-128" {

										if code.Data != this.lastBarcode {
											this.lastBarcode = code.Data
											if len(this.ballotBarcode) == cap(this.ballotBarcode) {
												<-this.ballotBarcode
											}
											this.ballotBarcode <- code.Data

											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}
											debugMarkList = []structs.CheckMarkList{}
											
								 

											currentState = this.setState("ballotPaperCode",currentState)
											this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,currentState)

										}

										if tesseractNeeded {
											
											result := this.tesseract(paper,this.currentOCRChannel)
											if len(result.PageRois)>0 {
												tesseractNeeded = false
												lastTesseractResult = result
												doFindCircles = true
												checkMarkList = []structs.CheckMarkList{}
												debugMarkList = []structs.CheckMarkList{}
												// fmt.Println("lastTesseractResult **",lastTesseractResult.Title)
												 

												currentState = this.setState("ballotPaperDetected",currentState)


											}else{
												 
												currentState = this.setState("ballotPaperNotDetected",currentState)
												fmt.Println("tesseract no bp found")
											}
											
										}

										if doFindCircles {

											neededRegions := 0
											for pRoiIndex := 0; pRoiIndex < len(lastTesseractResult.PageRois); pRoiIndex++ {
												titles := []string{}
					
												for i := 0; i < len(lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
													titles = append(titles, lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
												}
												foundIndex := IndexOf(titles, lastTesseractResult.Title)
												if (foundIndex>-1) {
													neededRegions++
												}
											}

											listOfRoiIndexes := []int{}


											foundIndex := -1
											for pRoiIndex := 0; pRoiIndex < len(lastTesseractResult.PageRois); pRoiIndex++ {
												titles := []string{}
					
												for i := 0; i < len(lastTesseractResult.PageRois[pRoiIndex].Types); i++ {
													titles = append(titles, lastTesseractResult.PageRois[pRoiIndex].Types[i].Title)
												}
												fIndex := IndexOf(titles, lastTesseractResult.Title)
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

													res := this.processRegionsOfInterest(lastTesseractResult,paper,listOfRoiIndexes)


													if false {
														log.Println("res.Marks",res.Marks,listOfRoiIndexes,res.IsCorrect)
													}
													//	log.Println("res.Id",lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id)
													
													currentState = this.setState("ballotPaperMarksAnalysed",currentState)


													for i := 0; i < len(res.Marks); i++ {
														if i >= len(debugMarkList) {

															offestX := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
															offestY := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)

															fmt.Println("XYZ",
																offestX + res.Marks[i].X, 
																offestY + res.Marks[i].Y,
															)

															debugMarkList = append(debugMarkList, structs.CheckMarkList{
																Point: image.Point{
																	offestX + res.Marks[i].X, 
																	offestY + res.Marks[i].Y,
																},
																Pixelsize: res.Marks[i].Radius,
															})
														}
														if res.Marks[i].Checked {
															debugMarkList[i].Sum += 1
														}
														debugMarkList[i].Count++
														debugMarkList[i].AVG = float64(debugMarkList[i].Sum) / float64(debugMarkList[i].Count)
														debugMarkList[i].Checked = debugMarkList[i].AVG > this.globals.SumMarksAVG
													}

													if res.IsCorrect {
														// log.Println("IsCorrect",res)
														// lastTesseractResult=res
														currentState = this.setState("isCorrect",currentState)
														this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,currentState)

														
														for i := 0; i < len(res.Marks); i++ {
															if i >= len(checkMarkList) {

																offestX := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
																offestY := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)

																/*
																fmt.Println("XYZ",
																	offestX + res.Marks[i].X, 
																	offestY + res.Marks[i].Y,
																)
																*/

																checkMarkList = append(checkMarkList, structs.CheckMarkList{
																	Point: image.Point{
																		offestX + res.Marks[i].X, 
																		offestY + res.Marks[i].Y,
																	},
																	Pixelsize: res.Marks[i].Radius,
																})
															}
															if res.Marks[i].Checked {
																checkMarkList[i].Sum += 1
															}
															checkMarkList[i].Count++
															checkMarkList[i].AVG = float64(checkMarkList[i].Sum) / float64(checkMarkList[i].Count)
															checkMarkList[i].Checked = checkMarkList[i].AVG > this.globals.SumMarksAVG
														}

														if len(checkMarkList)>0 && checkMarkList[0].Count>6 {
															//

															outList:=[]string{}
															for i := 0; i < len(checkMarkList); i++ {
																
																if checkMarkList[i].Checked {
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
																lastTesseractResult.PageRois[listOfRoiIndexes[0]].Types[foundIndex].Id,
																b.String(),
																"data:image/jpeg;base64,"+image_base64,
															)
															
															if err != nil {
																log.Println("SendReading ERROR",err)
																
																currentState = this.setState("sendError",currentState)
																this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,currentState)
															}else{
																// log.Println(">>>>",res.Msg)
																if res.Success {
																	doFindCircles = false


																	
																	currentState = this.setState("sendDone",currentState)
																	this.setHistoryItem(code.Data,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,currentState)
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
											currentState = this.setState("doFindCirclesDone",currentState)
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
								currentState = this.setState("noBarcodeFound",currentState)
								checkMarkList = []structs.CheckMarkList{}
								debugMarkList = []structs.CheckMarkList{}
							}

						}else{

							currentState = this.setState("findPaperContourFailed",currentState)

						}

										


						
						// gocv.Line(&xp, image.Point{0,0}, image.Point{200,img.Rows()},color.RGBA{0,0,255,0}, 20)

						// gocv.CvtColor(img, &img, gocv.ColorBGRToBGRA)

						/*
						for i := 0; i < len(debugMarkList); i++ {
							//playGround
							if debugMarkList[i].Checked {
								gocv.Circle(&playGround, debugMarkList[i].Point, debugMarkList[i].Pixelsize, color.RGBA{0, 155, 0, 0}, 10)
							}else{
								gocv.Circle(&playGround, debugMarkList[i].Point, debugMarkList[i].Pixelsize, color.RGBA{155, 0, 0, 0}, 10)
							}
						}
						*/

						// gocv.WarpPerspective(playGround, &img, invM, image.Point{img.Cols(), img.Rows()})
						
						// fmt.Printf("checkMarkList: %v   \n", checkMarkList )
						// if !doFindCircles && !tesseractNeeded && !paper.Empty()  {
							for i := 0; i < len(checkMarkList); i++ {
								//playGround
								if checkMarkList[i].Checked {
									gocv.Circle(&playGround, checkMarkList[i].Point, checkMarkList[i].Pixelsize, color.RGBA{0, 255, 0, 120}, int(3.0*this.pixelScale))
								}else{
									gocv.Circle(&playGround, checkMarkList[i].Point, checkMarkList[i].Pixelsize, color.RGBA{255, 0, 0, 120}, int(3.0*this.pixelScale))
								}
							}

							gocv.WarpPerspective(playGround, &img, invM, image.Point{img.Cols(), img.Rows()})
							
							
						// }

						drawContours := gocv.NewPointsVector()
						drawContours.Append(contour)
						gocv.DrawContours(&img, drawContours, -1, color.RGBA{uint8(currentState.Red), uint8(currentState.Green), uint8(currentState.Blue), 120}, int(8.0*this.pixelScale))
						drawContours.Close()


						


						invM.Close()
						paper.Close()
						playGround.Close()
						contour.Close()
					}
				}
			}else{
				log.Println("img empty")
			}
			log.Println("img ***")
			if len(this.imageChannelPaper)==cap(this.imageChannelPaper) {
				mat,_:=<-this.imageChannelPaper
				mat.Close()
			}

			cloned := img.Clone()
			gocv.Resize(cloned, &cloned, image.Point{}, 0.3, 0.3 , gocv.InterpolationLinear)
			this.imageChannelPaper <- cloned
			img.Close()
		}

		/*
		if len(this.escapedImage)>0 {
			escapedImage,escapedImageOk := <-this.escapedImage
			if escapedImageOk {
				log.Println("EscapedImage",escapedImage)

			}
		}
		*/
		//log.Println("processImage done %s",time.Since(start))
	}
	//log.Println("processImage exit",runVideo)
}