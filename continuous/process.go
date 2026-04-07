package continuous

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/api"
	"tualo.de/deep-test/args"
	"tualo.de/deep-test/barcode"
	"tualo.de/deep-test/layer2"
	"tualo.de/deep-test/tesseract"

	"tualo.de/deep-test/blurry"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/freezer"
	"tualo.de/deep-test/globals"
	"tualo.de/deep-test/marks"
	"tualo.de/deep-test/net"
	"tualo.de/deep-test/nn"
	"tualo.de/deep-test/paper"
	"tualo.de/deep-test/sheet"
	"tualo.de/deep-test/structs"
)

func Min(a, b string) string {
	if a < b {
		return a
	}
	return b
}

/*
Dies Klasse unterstützt dabei fortlaufende Bilder zu verarbeiten
*/
type Continuous struct {
	running          bool
	freeze           freezer.Freezer
	paper            paper.Paper
	ShowImageChannel chan structs.ShowImageStruct
	DrawLayer        *layer2.Layer2
	sheet            sheet.Sheet

	boxBarcode     string
	stackBarcode   string
	paginationCode string
	szIDCode       int
}

var static *Continuous

func Static() *Continuous {
	if static == nil {
		static = &Continuous{}
		static.running = false
		static.freeze = freezer.Freezer{}
		static.paper = paper.Paper{}
		static.ShowImageChannel = make(chan structs.ShowImageStruct, 1)
		static.DrawLayer = layer2.Static()

		static.freeze.Init()
		static.paper.Init()
	}
	return static
}

func (me *Continuous) IsRunning() bool {
	return me.running
}

func (me *Continuous) Text(
	img *cv.Mat,
	layerName string,
	baseLayerName string,
	text string,
	posX float64,
	posY float64,
	color color.RGBA,
) {
	newmat := cv.NewFromMat("Text", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))
	gocv.PutText(newmat.GetPointer(), text,
		image.Pt(
			int(float64(newmat.Cols())*posX),
			int(float64(newmat.Cols())*posY),
		), gocv.FontHersheyPlain,
		8.2,
		color,
		22,
	)
	me.DrawLayer.Set(layerName, newmat, false, baseLayerName, 0.5, 1.0, 0.0)
	newmat.Close()
}

func (me *Continuous) TextSz(
	img *cv.Mat,
	layerName string,
	baseLayerName string,
	text string,
	posX float64,
	posY float64,
	color color.RGBA,
	size float64,
) {
	newmat := cv.NewFromMat("Text", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))
	gocv.PutText(newmat.GetPointer(), text,
		image.Pt(
			int(float64(newmat.Cols())*posX),
			int(float64(newmat.Cols())*posY),
		), gocv.FontHersheyPlain,
		size,
		color,
		int(size*2),
	)
	me.DrawLayer.Set(layerName, newmat, false, baseLayerName, 1.9, 1, 0.0)
	newmat.Close()
}

func (me *Continuous) TextXY(
	img *cv.Mat,
	layerName string,
	baseLayerName string,
	text string,
	posX int,
	posY int,
	color color.RGBA,
) {
	newmat := cv.NewFromMat("TextXY", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))
	gocv.PutText(newmat.GetPointer(), text,
		image.Pt(
			posX,
			posY,
		), gocv.FontHersheyPlain,
		8.2,
		color,
		22,
	)
	me.DrawLayer.Set(layerName, newmat, false, baseLayerName, 0.5, 1.0, 0.0)
	newmat.Close()
}

/*
takes an image and if there is an free process,
the image will be analysed
*/
func (me *Continuous) Process(img cv.Mat) {
	if me.running {
		return
	}
	me.running = true
	cvx := img.Clone("local")
	start := time.Now()

	if !cvx.GetPointer().Empty() {

		// prüfen, ob das bild noch in bewegung ist
		if me.freeze.Check(&cvx) {
			if *args.Parser().BlurryThreshold > 0 {
				bM := cvx.Clone("blurry check")
				blurred, variance := blurry.Static().Detect(&bM)
				if blurred {

					me.Text(
						&cvx,
						"blurry",
						"",
						fmt.Sprintf("Unscharf: %0.2f", variance),
						0.05,
						0.95,
						color.RGBA{255, 0, 0, 0},
					)

					bM.Close()
					cvx.Close()
					me.running = false
					return
				} else {

					me.Text(
						&cvx,
						"blurry",
						"",
						fmt.Sprintf("Scharf: %0.2f", variance),
						0.05,
						0.95,
						color.RGBA{155, 255, 0, 0},
					)

				}
				bM.Close()
			}

			me.Paper(cvx)
		} else {
			me.DrawLayer.Remove("rois")
		}

		me.Text(
			&cvx,
			"duration",
			"",
			fmt.Sprintf("D: %s", time.Since(start)),
			0.05,
			0.99,
			color.RGBA{0, 255, 0, 0},
		)
		// log.Println("duration", time.Since(start))

	}
	cvx.Close()
	me.running = false
}

// Paper processes the paper image and extracts the contour
// and checks for the paper type
// it also sets the text on the image
// and draws the contour if available

func (me *Continuous) Paper(img cv.Mat) {
	cvx := img.Clone("local")

	cnt, cmat := me.paper.Contour(cvx)

	if cnt.Size() > 0 {

		f := gocv.ArcLength(cnt, true)
		approx := gocv.ApproxPolyDP(cnt, globals.Globals().ApproxPolyDPFactor*f, true)

		if approx.Size() == 4 {
			points := me.paper.GetCornerPoints(cnt)

			me.Text(
				&cvx,
				"shape",
				"",
				fmt.Sprintf("Obj: %d %0.2f", approx.Size(), f),
				0.60,
				0.05,
				color.RGBA{0, 255, 255, 0},
			)

			if cmat.GetPointer().Channels() == 3 {
				me.DrawLayer.Set("contour", cmat, false, "", 0.5, 1.0, 0.0)
				cmat.Close()
			}

			diffs := me.paper.EdgeDifference(points) - 1
			if math.Round(diffs*10) < 10 {

				extract, inverseMat := me.paper.Extract(&cvx, points)
				me.DrawLayer.Set("inverseMat", inverseMat, true, "", 0, 0, 0)
				if len(me.ShowImageChannel) == cap(me.ShowImageChannel) {
					o := <-me.ShowImageChannel
					o.Mat.Close()
				}
				//gocv.IMWrite("sample.jpg", cvx.Get())
				me.Sheet(&extract)

				inverseMat.Close()
				extract.Close()
				// log.Println(">>>>>>>>>>>>>>>", diffs)
			}
		} else {

			me.Text(
				&cvx,
				"shape",
				"",
				fmt.Sprintf("Obj: %d %0.2f", approx.Size(), f),
				0.60,
				0.05,
				color.RGBA{255, 0, 0, 0},
			)
		}

		approx.Close()

	}
	cmat.Close()
	cnt.Close()
	cvx.Close()
}

func (me *Continuous) CheckBarcodes(codes []structs.BarcodeSymbol) (bool, bool, bool) {

	boxBarcode := "Unkown"
	stackBarcode := "Unkown"
	paginationCode := "Unkown"
	szID := "Unkown"

	changedBoxBarcode := false
	changedStackBarcode := false
	changedPaginationCode := false

	for _, code := range codes {

		log.Println("Code", code.Type, code.Data, len(code.Data))
		if code.Type == "CODE-39" {
			if len(code.Data) >= 5 {
				if code.Data[0:3] == "FC4" {
					boxBarcode = code.Data
				}
				if code.Data[0:3] == "FC3" {
					stackBarcode = code.Data
				}

			}
			if len(code.Data) == 4 {
				// min barcode is the right szID
				log.Println("SZ ID candidate", code.Data)
				szID = Min(szID, code.Data)

			}
		}
		if code.Type == "CODE-128" {

			if len(code.Data) >= 5 {
				paginationCode = code.Data
			}

		}
	}

	if me.boxBarcode != boxBarcode {
		me.boxBarcode = boxBarcode
		changedBoxBarcode = true
	}
	if me.stackBarcode != stackBarcode {
		me.stackBarcode = stackBarcode
		changedStackBarcode = true
	}
	if me.paginationCode != paginationCode {
		me.paginationCode = paginationCode
		changedPaginationCode = true
	}

	return changedBoxBarcode, changedStackBarcode, changedPaginationCode

	// changedPaginationCode = true
	// me.paginationCode = "13912753612"

	// return changedBoxBarcode, changedStackBarcode, changedPaginationCode

}

func (me *Continuous) Sheet(img *cv.Mat) {
	var codes []structs.BarcodeSymbol
	var foundCodes bool
	cvx := img.Clone("Sheet")

	id := -1

	if len(api.Static().CandidateBarcodeList) == 0 {
		id = tesseract.Static().DetectBallotpaperType(img)
	} else {
		codes, foundCodes = barcode.Static().FindCandidateCode(img)

		if foundCodes {
			candidateBarcodeList := api.Static().CandidateBarcodeList
			mixedBallotpaper := false
			foundBallotpaper := false
			lastId := -1
			for _, code := range codes {
				// log.Println("Candidate Code", code.Type, code.Data, len(code.Data))
				for _, candidate := range candidateBarcodeList {
					if code.Data == candidate.Barcode {
						if lastId != -1 && lastId != candidate.Ballotpaper {
							log.Println("Multiple candidate codes with different ballotpapers found", lastId, candidate.Ballotpaper)
							mixedBallotpaper = true
						} else {
							lastId = candidate.Ballotpaper
							foundBallotpaper = true
						}
					}
				}
			}
			if foundBallotpaper && !mixedBallotpaper {
				id = lastId
			} else {
				log.Println("No unique candidate code found, try again")
			}

		}
		// log.Println("***Codes Found:", codes, foundCodes)
		// id = api.Static().DetectBallotpaperTypeByBarcode(me.szIDCode)
	}

	log.Println("BP Id", id)
	// id = 0
	if id != -1 {
		rois := api.Static().RoisByBallopPaperID(id)
		me.sheet.SetImage(cvx.GetPointer())

		bp := api.Static().GetBallotpaper(id)
		if bp == nil {
			log.Println("No ballotpaper found for id", id)
			return
		}

		title := bp.Title
		// log.Println("BP", title, cvx.GetPointer().Channels())
		me.TextSz(
			&cvx,
			"bptitle",
			"inverseMat",
			title,
			0.05,
			0.95,
			color.RGBA{255, 0, 0, 0},
			3.2,
		)

		me.sheet.Measurements(
			bp.PageWidth,
			bp.PageHeight,
		)
		scale := me.sheet.Rescale()

		codes, foundCodes = barcode.Static().Find(img)
		log.Println("Pagination Codes Found:", foundCodes)
		foundCodes = true
		if foundCodes {
			// log.Println(codes)
			changedBoxBarcode, changedStackBarcode, changedPaginationCode := me.CheckBarcodes(codes)
			if changedPaginationCode {
				log.Println("Pagination Code", me.paginationCode)
			}
			me.szIDCode = id

			log.Println("Found SZ ID:", me.szIDCode, " title:", bp.Title, " scale:", scale, " rois:", rois)
			res := me.SheetRois(&cvx, &scale, rois)
			Display(res)
			// }
			if false {
				log.Println(changedBoxBarcode, changedStackBarcode)
			}
		}

		// barcode
		// tesseract
	}
	cvx.Close()
}

type Results struct {
	/*
		Index int
		Identifier string
	*/
	Done        bool
	Marked      bool
	MarkedValue float64
}

/*
var mat_monitor map[string]*MonitorData
mat_monitor = make(map[string]*MonitorData)
*/

func Display(dis map[string][]Results) {

	fmt.Println("================ results ================")
	s := fmt.Sprintf("%-10s", "Index")
	v1 := fmt.Sprintf("%*s", 10, "BV")
	v2 := fmt.Sprintf("%*s", 10, "NN")
	v3 := fmt.Sprintf("%*s", 10, "LN")
	fmt.Println(s, v1, v2, v3)
	for i := 0; i < len(dis["bv"]); i++ {

		fmt.Println(
			fmt.Sprintf("%-10d", i),
			fmt.Sprintf("%*f", 10, dis["bv"][i].MarkedValue),
			fmt.Sprintf("%*f", 10, dis["nn"][i].MarkedValue),
			fmt.Sprintf("%*f", 10, dis["ln"][i].MarkedValue),
		)
	}
	fmt.Println("=========================================")
}

/*
SheetRois processes the ROIS for a sheet image

	and returns a map with results for each roi type
*/
func (me *Continuous) SheetRois(img *cv.Mat, scale *structs.Scale, rois []structs.ROI) map[string][]Results {

	result := make(map[string][]Results)

	sum := 0
	res_bv := []Results{}
	res_nn := []Results{}
	res_ln := []Results{}
	for r := 0; r < len(rois); r++ {
		currentRoi := rois[r]
		sum += currentRoi.ItemCountX * currentRoi.ItemCountY
		for i := 0; i < currentRoi.ItemCountY; i++ {

			res_bv = append(res_bv, Results{
				Done:        false,
				Marked:      false,
				MarkedValue: 0.0,
			})

			res_nn = append(res_nn, Results{
				Done:        false,
				Marked:      false,
				MarkedValue: 0.0,
			})

			res_ln = append(res_ln, Results{
				Done:        false,
				Marked:      false,
				MarkedValue: 0.0,
			})
		}
	}
	result["bv"] = res_bv
	result["nn"] = res_nn
	result["ln"] = res_ln

	// log.Println("img.Type()", img.Type().String())

	newmat := cv.NewFromMat("rois image", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))

	index := 0
	for r := 0; r < len(rois); r++ {
		currentRoi := rois[r]
		log.Println("Roi", r, "of", len(rois), "ItemCountY", currentRoi.ItemCountY, "ItemCountX", currentRoi.ItemCountX, currentRoi)

		for i := 0; i < currentRoi.ItemCountY; i++ {
			scaleX := scale.PixelScaleX
			scaleY := scale.PixelScaleY

			y := (currentRoi.Y + (((currentRoi.Height) + currentRoi.YCap + 0.5) * float64(i)))
			rect := image.Rect(
				int(currentRoi.X*scaleX),
				int(y*scaleY),
				int(currentRoi.X*scaleX+currentRoi.Width*scaleY),
				int(y*scaleY+currentRoi.Height*scaleY),
			)

			gocv.Rectangle(newmat.GetPointer(), rect, color.RGBA{255, 255, 0, 0}, 15)

			// log.Println(rect)
			if rect.Min.X < 0 {
				newmat.Close()
				return result
			}
			if rect.Min.Y < 0 {
				newmat.Close()
				return result
			}
			if rect.Max.X > img.GetPointer().Cols() {
				newmat.Close()
				return result
			}
			if rect.Max.Y > img.GetPointer().Rows() {
				newmat.Close()
				return result
			}
			item := img.GetPointer().Region(rect)

			mark := marks.Marks{
				InternalId: i,
			}
			mark.Set(&item, *scale)
			gocv.CvtColor(item, &item, gocv.ColorBGRToGray)

			var best gocv.Mat
			mark.SetBestContourTrigger(func(m gocv.Mat) {
				best = m
			})

			result["bv"][index].Done = false
			result["bv"][index].Marked = false
			result["bv"][index].MarkedValue = 0.0

			markResult := mark.Check()

			if !best.Closed() && !best.Empty() {

				result["bv"][index].Done = markResult.Found
				result["bv"][index].Marked = markResult.Marked
				result["bv"][index].MarkedValue = markResult.Value

				x := (rect.Max.X - rect.Min.X) / 2
				y := (rect.Max.Y - rect.Min.Y) / 2
				r := x / 2

				if markResult.Marked {
					gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r, color.RGBA{0, 255, 0, 255}, 15)
				} else {
					gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r, color.RGBA{255, 255, 255, 255}, 15)
					gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r, color.RGBA{255, 0, 0, 255}, 15)
				}

				if *args.Parser().EnableLocalNN {
					s, f := net.Static().Predict(&best)
					// log.Println(s, f)
					v := 0.0
					if s == "X" {
						v = f
					}
					result["ln"][index].Done = true
					result["ln"][index].Marked = s == "X"
					result["ln"][index].MarkedValue = v

					if s == "X" {
						gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r+10, color.RGBA{0, 255, 0, 255}, 15)

						testSize := gocv.GetTextSize(
							fmt.Sprintf("LN: %0.2f%%", f),
							gocv.FontHersheyPlain,
							3.2,
							6,
						)
						gocv.PutText(
							newmat.GetPointer(),
							fmt.Sprintf("LN: %0.2f%%", f),
							image.Point{
								rect.Min.X - int(float64(testSize.X)*1.3),
								rect.Min.Y + y + testSize.Y/2,
							},
							gocv.FontHersheyPlain,
							3.2,
							color.RGBA{100, 100, 100, 255},
							6,
						)

					} else {
						gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r+10, color.RGBA{255, 255, 255, 255}, 15)
						gocv.Circle(newmat.GetPointer(), image.Point{rect.Min.X + x, rect.Min.Y + y}, r+10, color.RGBA{255, 0, 0, 255}, 15)

						/*
							testSize := gocv.GetTextSize(
								fmt.Sprintf("LN: %0.2f%%", f),
								gocv.FontHersheyPlain,
								3.2,
								6,
							)
							gocv.PutText(
								newmat.GetPointer(),
								fmt.Sprintf("LN: %0.2f%%", f),
								image.Point{
									rect.Min.X - int(float64(testSize.X)*1.3),
									rect.Min.Y + y + testSize.Y/2,
								},
								gocv.FontHersheyPlain,
								3.2,
								color.RGBA{180, 100, 100, 255},
								6,
							)
						*/
					}

				}

				if *args.Parser().EnableNNServiceQuery {
					f, _ := nn.Service().RequestImageFloat(&best)
					result["nn"][index].Done = true
					result["nn"][index].Marked = f > 30
					result["nn"][index].MarkedValue = f
				}
				best.Close()
			}
			item.Close()
			index++
		}

	}
	me.DrawLayer.Set("rois", newmat, false, "inverseMat", 2, 1, 5.0)

	newmat.Close()

	return result

}
