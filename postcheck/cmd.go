package postcheck

import (
	"fmt"
	"image"
	"log"
	"math"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/args"
	"tualo.de/deep-test/blurry"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/dbconnector"
	"tualo.de/deep-test/marks"
	"tualo.de/deep-test/net"
	"tualo.de/deep-test/nn"
	"tualo.de/deep-test/sheet"
	"tualo.de/deep-test/structs"
)

var window *gocv.Window
var dbconn = dbconnector.DBConnector{}

func show(mat gocv.Mat) {
	if window != nil {
		if window.IsOpen() {
			window.IMShow(mat)
			window.WaitKey(0)
		}
	}
}

func page(mat *gocv.Mat, width int, height int, pagination string) {

	var sheet = sheet.Sheet{}
	sheet.SetImage(mat)
	sheet.Measurements(width, height)

	show(*mat)
	scale := sheet.Rescale()
	if *args.Parser().BlurryThreshold > 0 {
		bM := cv.NewFromMat("blurry", mat.Clone())
		blurred, _ := blurry.Static().Detect(&bM)
		if blurred {
			log.Println(pagination, "too", "blurry!")
			mat.Close()
			return
		} else {
			log.Println(pagination, "not", "blurry")
		}
		bM.Close()
	}
	rois := dbconn.ReadFields(pagination)
	index := 0
	for r := 0; r < len(rois.Roi); r++ {
		currentRoi := rois.Roi[r]

		var toSave []structs.SaveMarker
		hasError := false

		for i := 0; i < currentRoi.ItemCountY; i++ {
			scaleX := scale.PixelScaleX
			scaleY := scale.PixelScaleY

			y := (currentRoi.Y + (((currentRoi.Height) + currentRoi.YCap) * float64(i)))
			rect := image.Rect(
				int(currentRoi.X*scaleX),
				int(y*scaleY),
				int(currentRoi.X*scaleX+currentRoi.Width*scaleY),
				int(y*scaleY+currentRoi.Height*scaleY),
			)
			item := mat.Region(rect)
			mark := marks.Marks{}
			mark.Set(&item, scale)
			gocv.CvtColor(item, &item, gocv.ColorBGRToGray)
			var best gocv.Mat
			mark.SetBestContourTrigger(func(m gocv.Mat) {
				best = m
			})
			res := mark.Check()

			if !best.Closed() && !best.Empty() {

				ratioDifference := math.Abs(math.Round((1 - float64(best.Cols())/float64(best.Rows())) * 100))

				if *args.Parser().SaveBestMarkerImage {
					var marked string
					marked = "O"
					if res.Marked {
						marked = "X"
					}
					if res.Value > 0.5 {
						marked = "Q"
					}
					toSave = append(toSave, structs.SaveMarker{
						Mat:        best.Clone(),
						Marked:     marked,
						Pagination: pagination,
						R:          r,
						I:          i,
						V:          int(res.Value * 255),
					})
				}
				if true {
					dbconn.StoreRestult(
						pagination,
						"bildverarbeitung",
						index,
						res.Value*255.0,
					)
				}
				hasError = ratioDifference > 5

				if *args.Parser().EnableNNServiceQuery {
					f, _ := nn.Service().RequestImageFloat(&best)
					dbconn.StoreRestult(
						pagination,
						"tensorflow 32x32x150",
						index,
						f,
					)
				}

				if *args.Parser().EnableLocalNN {
					s, f := net.Static().Predict(&best)
					v := 0.0
					if s == "X" {
						v = f
					}

					dbconn.StoreRestult(
						pagination,
						"local nn",
						index,
						v,
					)
				}

				if *args.Parser().EnableNNServiceQuery {
					f, _ := nn.Service().RequestImageFloat(&best)
					dbconn.StoreRestult(
						pagination,
						"tensorflow 32x32x150",
						index,
						f,
					)
				}
				best.Close()
			} else {
				hasError = true
			}
			item.Close()
			index++
		}

		if !hasError && *args.Parser().SaveBestMarkerImage {
			for i := 0; i < len(toSave); i++ {
				gocv.IMWrite(
					fmt.Sprintf(
						"%s/%s/%s.%d.%d.%d.jpg",
						*args.Parser().SaveBestMarkerImagePath,
						toSave[i].Marked,
						toSave[i].Pagination,
						toSave[i].R,
						toSave[i].I,
						toSave[i].V,
					),
					toSave[i].Mat,
				)
			}
		}
		for i := 0; i < len(toSave); i++ {
			toSave[i].Mat.Close()
		}

	}

	mat.Close()
}

func Run() {
	arguments := args.Parser()
	if *arguments.EnableCVWindow {
		window = gocv.NewWindow("Hello")
	}
	dbconn.Run(*arguments.DBConnection, *arguments.DBFilter, page)
}
