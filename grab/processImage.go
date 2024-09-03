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

func (this *GrabcameraClass) processImage(){
	scanner := barcode.NewScanner()
	scanner.SetEnabledAll(false)
	scanner.SetEnabledSymbology(barcode.Code39,true)
	scanner.SetEnabledSymbology(barcode.Code128,true)
	log.Println("processImage starting ")
	tesseractNeeded := true
	lastTesseractResult := structs.TesseractReturnType{}
	lastBarcode := "wlekfjwuqezgzw"
	doFindCircles := false
	checkMarkList := []structs.CheckMarkList{}
	debugMarkList := []structs.CheckMarkList{}
	strCurrentBoxBarcode := ""
	strCurrentStackBarcode := ""
	
	//lastCheckMarkList := []structs.CheckMarkList{}

	green := 0
	red := 0
	blue := 0
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

//			green = 0
//			red = 0
//			blue = 0

			if !img.Empty() {


				/*
				meanStart := time.Now()
				img_mean := img.Mean()
				log.Println("Mean: ",time.Since(meanStart),img_mean.Val1)
				*/

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
							img.Close()
							paper.Close()
							invM.Close()
							continue
						}
						playGround := paper.Clone()
						

						area := float64(paper.Size()[0]) * float64(paper.Size()[1]) / float64(img.Size()[0]) / float64(img.Size()[1])
						// log.Println("extractPaper done %s %f",time.Since(start),area)
						if area > 0.1 {
							codes := this.findBarcodes(scanner,paper)
							if false {
								log.Println("findBarcodes done %s %v",time.Since(start),codes)
							}
							if len(codes) > 0 {
								for _, code := range codes {

									//fmt.Println("code **",code.Type,code.Data)
									if code.Type == "CODE-39" {
										if code.Data[0:3]=="FC4" {
											if len(this.currentBoxBarcode) == cap(this.currentBoxBarcode) {
												<-this.currentBoxBarcode
											}
											this.currentBoxBarcode <- code.Data
											strCurrentBoxBarcode = code.Data

											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}
										}
										if code.Data[0:3]=="FC3" {
											if len(this.currentStackBarcode) == cap(this.currentStackBarcode) {
												<-this.currentStackBarcode
											}
											this.currentStackBarcode <- code.Data
											strCurrentStackBarcode = code.Data

											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}

										}
									}
									if code.Type == "CODE-128" {

										if code.Data != lastBarcode {
											lastBarcode = code.Data
											if len(this.ballotBarcode) == cap(this.ballotBarcode) {
												<-this.ballotBarcode
											}
											this.ballotBarcode <- code.Data

											log.Println(">>>>> RESET code",lastBarcode)
											tesseractNeeded = true
											doFindCircles = false
											checkMarkList = []structs.CheckMarkList{}
											
											green = 0
											red = 0
											blue = 255
										}

										if tesseractNeeded {
											
											result := this.tesseract(paper)
											if len(result.PageRois)>0 {
												tesseractNeeded = false
												lastTesseractResult = result
												doFindCircles = true
												checkMarkList = []structs.CheckMarkList{}
												debugMarkList = []structs.CheckMarkList{}
												// fmt.Println("lastTesseractResult **",lastTesseractResult.Title)
												green = 255
												red = 255
												blue = 255

											}else{
												green = 100
												red = 255
												blue = 100
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

											log.Println("listOfRoiIndexes",listOfRoiIndexes)
											log.Println("==================================")

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


													if true {
														log.Println("res.Marks",res.Marks,listOfRoiIndexes,res.IsCorrect)
													}
													//	log.Println("res.Id",lastTesseractResult.PageRois[pRoiIndex].Types[foundIndex].Id)
													green = 255
													red = 100
													blue = 255


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

														
														for i := 0; i < len(res.Marks); i++ {
															if i >= len(checkMarkList) {

																offestX := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].X) * this.pixelScale)	
																offestY := int(float64(lastTesseractResult.PageRois[res.Marks[i].RoiIndex].Y) * this.pixelScaleY)

																fmt.Println("XYZ",
																	offestX + res.Marks[i].X, 
																	offestY + res.Marks[i].Y,
																)

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
															res.Barcode = lastBarcode
															fmt.Printf("Box: %s, Stack: %s, Barcode: %s, Title: %s, List: %v \n",res.BoxBarcode,res.StackBarcode, res.Barcode , lastTesseractResult.Title, outList)
															//checkMarkList = sumMarks(checkMarkList, res)
															//lastCheckMarkList = checkMarkList
															b := new(strings.Builder)
															json.NewEncoder(b).Encode(outList)


															

															image_bytes, _ := gocv.IMEncode(gocv.JPEGFileExt, paper)
															image_base64 := base64.StdEncoding.EncodeToString(image_bytes.GetBytes())
															//fmt.Println("pic",image_base64[0:100])
															image_bytes.Close()

															res,err := api.SendReading(
																strCurrentBoxBarcode,
																strCurrentStackBarcode,
																res.Barcode,
																lastTesseractResult.PageRois[listOfRoiIndexes[0]].Types[foundIndex].Id,
																b.String(),
																"data:image/jpeg;base64,"+image_base64,
															)
															
															if err != nil {
																log.Println("SendReading ERROR",err)
																green = 0
																red = 255
																blue = 0
															}else{
																// log.Println(">>>>",res.Msg)
																if res.Success {
																	doFindCircles = false


																	green = 250
																	red = 0
																	blue = 0
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
											// green = 50
											// red = 0
											// blue = 50
										}
										//log.Println("code use tesseract",code.Data,tesseractNeeded,lastTesseractResult)
									}
								}
								// gocv.IMWrite("paper.png",paper)
							}

						}else{
							green = 0
							red = 0
							blue = 0
						}

										


						
						// gocv.Line(&xp, image.Point{0,0}, image.Point{200,img.Rows()},color.RGBA{0,0,255,0}, 20)

						gocv.CvtColor(img, &img, gocv.ColorBGRToBGRA)

						for i := 0; i < len(debugMarkList); i++ {
							//playGround
							if debugMarkList[i].Checked {
								gocv.Circle(&playGround, debugMarkList[i].Point, debugMarkList[i].Pixelsize, color.RGBA{0, 155, 0, 0}, 10)
							}else{
								gocv.Circle(&playGround, debugMarkList[i].Point, debugMarkList[i].Pixelsize, color.RGBA{155, 0, 0, 0}, 10)
							}
						}

						gocv.WarpPerspective(playGround, &img, invM, image.Point{img.Cols(), img.Rows()})
						

						if !doFindCircles && !tesseractNeeded && !paper.Empty() && (red + green + blue > 0) {
							for i := 0; i < len(checkMarkList); i++ {
								//playGround
								if checkMarkList[i].Checked {
									gocv.Circle(&playGround, checkMarkList[i].Point, checkMarkList[i].Pixelsize, color.RGBA{0, 255, 0, 0}, 20)
								}else{
									gocv.Circle(&playGround, checkMarkList[i].Point, checkMarkList[i].Pixelsize, color.RGBA{255, 0, 0, 0}, 20)
								}
							}

							gocv.WarpPerspective(playGround, &img, invM, image.Point{img.Cols(), img.Rows()})
							
							
						}

						drawContours := gocv.NewPointsVector()
						drawContours.Append(contour)
						

						// fmt.Printf("red: %i, green: %i, blue: %i  \n", red,green,blue)
															
						gocv.DrawContours(&img, drawContours, -1, color.RGBA{uint8(red), uint8(green), uint8(blue), 120}, int(8.0*this.pixelScale))
						drawContours.Close()


						


						invM.Close()
						paper.Close()
						playGround.Close()
						contour.Close()
					}
				}
			}
			if len(this.imageChannelPaper)==cap(this.imageChannelPaper) {
				mat,_:=<-this.imageChannelPaper
				mat.Close()
			}
			cloned := img.Clone()
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