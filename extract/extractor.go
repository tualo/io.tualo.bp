package extract

import (
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"

	"encoding/base64"

	_ "github.com/go-sql-driver/mysql"
	"gocv.io/x/gocv"
)

var currentImageName string
var currentRoiRect int

type Payload struct {
	Stuff Data
}
type Data struct {
	Rows    []string
	Success bool
}

type CircleState struct {
	Found  bool
	Marked bool
	Value  float64
	Image  gocv.Mat
}

type Extractor struct {
	db  *sql.DB
	wnd *gocv.Window
}

func (me *Extractor) Run(w *gocv.Window) {
	me.wnd = w
	me.db = me.connectDB()
	defer me.db.Close()

	me.getPaginations()
}

// bestContour obtains the biggest contour in the frame provided is bigger
// than the minArea.
func (me *Extractor) bestContour(frame gocv.Mat, minArea float64) (image.Rectangle, gocv.PointsVector, int) {
	circleIndex := -1
	frameContour := 3
	gocv.Rectangle(&frame, image.Rectangle{image.Point{0, 0}, image.Point{frame.Cols(), frame.Rows()}}, color.RGBA{255, 255, 255, 120}, frameContour)
	frameX := frame.Clone()
	frame.ConvertTo(&frameX, gocv.MatTypeCV32SC1)
	cnts := gocv.FindContours(frameX, gocv.RetrievalFloodfill, gocv.ChainApproxNone)
	var (
		bestCnt  gocv.PointVector
		bestArea = minArea
		found    = false
	)
	largest := 0.0
	for i := 0; i < cnts.Size(); i++ {
		if area := gocv.ContourArea(cnts.At(i)); area > largest {
			largest = (area)
		}
	}
	for i := 0; i < cnts.Size(); i++ {
		cnt := cnts.At(i)
		area := gocv.ContourArea(cnt)
		if area != largest {
			//if currentImageName == "20001" && currentRoiRect == 19 {
			// MinEnclosingCircle
			appoxCnt := gocv.ApproxPolyDP(cnt, 0.03*gocv.ArcLength(cnt, true), true)
			log.Println("appoxCnt", appoxCnt.Size())
			r := rand.New(rand.NewSource(255))
			log.Println("r", r)
			// gocv.DrawContours(&show, cnts, i, color.RGBA{0, 0, 255, 0}, 3)

			if appoxCnt.Size() == 8 {
				_, _, radius := gocv.MinEnclosingCircle(appoxCnt)
				circleArea := radius * radius * math.Pi
				log.Println("XYZ", circleArea, area)
				if math.Abs(math.Round(((float64(circleArea)/area)-1)*100)) < 5 {
					circleIndex = i
				}
			}

			appoxCnt.Close()

			// gocv.DrawContour( {appoxCnt}
			// }
			if area > bestArea {
				bestArea = area
				bestCnt = cnt
				found = true
			}
		}
	}

	//}
	if !found {
		log.Panicln(">>>>>>>>>>>>>>>>>>>>>>")
		return image.Rect(0, 0, frame.Cols(), frame.Rows()), cnts, circleIndex
	}

	pnts := bestCnt.ToPoints()
	x1 := 9999
	y1 := 9999
	x2 := 0
	y2 := 0
	for i := 0; i < len(pnts); i++ {
		x1 = int(math.Min(float64(x1), float64(pnts[i].X)))
		y1 = int(math.Min(float64(y1), float64(pnts[i].Y)))
		x2 = int(math.Max(float64(x2), float64(pnts[i].X)))
		y2 = int(math.Max(float64(y2), float64(pnts[i].Y)))
	}

	return image.Rect(x1, y1, x2, y2), cnts, circleIndex
}

func (me *Extractor) connectDB() *sql.DB {
	var err error
	// Replace "username", "password", "dbname" with your database credentials
	connectionString := "thomashoffmann:@tcp(127.0.0.1:3306)/bwnuernberg"
	me.db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return me.db
}

func (me *Extractor) circles(img gocv.Mat, pixelScaleX float64, pixelScaleY float64) CircleState {
	croppedMatGray := img.Clone()
	useResult := img.Clone()
	circles := gocv.NewMat()

	dpHoughCircles := 1.0
	thresholdHoughCircles := 90.0
	accumulatorThresholdHoughCircles := 10.0

	//innerOverdraw := int(float64(0.01) * math.Min(pixelScaleX, pixelScaleY))
	circleSize := int(float64(5) * math.Min(pixelScaleX, pixelScaleY))
	minDist := (float64(15) * math.Min(pixelScaleX, pixelScaleY))

	strokeWidth := (float64(4) * math.Max(pixelScaleX, pixelScaleY))

	gocv.HoughCirclesWithParams(
		croppedMatGray,
		&circles,
		gocv.HoughGradient,
		dpHoughCircles,                   // dp
		minDist,                          //float64(croppedMatGray.Rows()/50), // minDist
		thresholdHoughCircles,            // param1
		accumulatorThresholdHoughCircles, // param2
		circleSize,                       // minRadius
		circleSize,                       // maxRadius
	)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			sx := x - r - int(pixelScaleX*0.5)
			sy := y - r - int(pixelScaleY*0.5)
			ex := x + r + int(pixelScaleX*0.5)
			ey := y + r + int(pixelScaleY*0.5)

			if sx < 0 {
				sx = 0
			}
			if sx > croppedMatGray.Cols() {
				// ToDo TriState: No Circle Found
				return CircleState{
					Marked: false,
					Found:  false,
					Value:  -1,
				}
			}
			if sy < 0 {
				sy = 0
			}
			if sy > croppedMatGray.Rows() {
				// ToDo TriState: No Circle Found
				return CircleState{
					Marked: false,
					Found:  false,
					Value:  -1,
				}
			}

			if ex > croppedMatGray.Cols() {
				ex = croppedMatGray.Cols()
			}
			if ey > croppedMatGray.Rows() {
				ey = croppedMatGray.Rows()
			}

			rect := image.Rect(
				sx,
				sy,
				ex,
				ey,
			)

			me.wnd.IMShow(croppedMatGray)
			me.wnd.WaitKey(0)
			gocv.Threshold(croppedMatGray, &croppedMatGray, 45, 255, gocv.ThresholdBinary)
			me.wnd.IMShow(croppedMatGray)
			me.wnd.WaitKey(0)
			gocv.Dilate(croppedMatGray, &croppedMatGray, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))
			gocv.Erode(croppedMatGray, &croppedMatGray, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(7, 7)))
			me.wnd.IMShow(croppedMatGray)
			me.wnd.WaitKey(0)
			croppedMatGray = croppedMatGray.Region(rect)
			useResult = useResult.Region(rect)
			best, contours, circleIndex := me.bestContour(croppedMatGray, float64(r*r)*0.5)

			log.Println("circleIndex>>>>", circleIndex)
			if circleIndex > -1 {
				// wenn keine kreiskontur gefunden wurde (-1),
				// gint es verbindungen zum kreis > ergo eine markierung
				gocv.DrawContours(&useResult, contours, circleIndex, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
			} else {

				gocv.HoughCirclesWithParams(
					useResult,
					&circles,
					gocv.HoughGradient,
					dpHoughCircles,                   // dp
					minDist,                          //float64(croppedMatGray.Rows()/50), // minDist
					thresholdHoughCircles,            // param1
					accumulatorThresholdHoughCircles, // param2
					circleSize,                       // minRadius
					circleSize,                       // maxRadius
				)

				for i := 0; i < circles.Cols(); i++ {
					v := circles.GetVecfAt(0, i)
					if len(v) > 2 {
						x := int(v[0])
						y := int(v[1])
						r := int(v[2])
						gocv.Circle(&useResult, image.Pt(x, y), r, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
					}
				}
			}
			gocv.Erode(useResult, &useResult, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))

			croppedMatGray = croppedMatGray.Region(best)
			useResult = useResult.Region(best)

			gocv.Threshold(useResult, &useResult, 150, 255, gocv.ThresholdBinary)
			gocv.BitwiseNot(useResult, &useResult)

			cx := (float64(useResult.Cols()) * float64(useResult.Rows()))
			avg := ((useResult.Sum().Val1) / cx) / 255.0
			log.Println("avg****", avg, useResult.Sum().Val1/cx, cx)

			/*
				me.wnd.IMShow(croppedMatGray)
				me.wnd.WaitKey(0)
			*/
			me.wnd.IMShow(useResult)
			me.wnd.WaitKey(0)

			return CircleState{
				Marked: (useResult.Sum().Val1 / cx) > 0.0,
				Found:  true,
				Value:  avg,
				Image:  useResult,
			}

			/*
				} else {
					log.Println(rect)
					me.wnd.IMShow(croppedMatGray)
					me.wnd.WaitKey(0)
				}
			*/
		}
	}

	// gocv.Resize(useResult, &useResult, image.Point{128, 128}, 0, 0, gocv.InterpolationArea)

	return me.circleRating(useResult, pixelScaleX, pixelScaleY)

}

func (me *Extractor) circleRating(img gocv.Mat, pixelScaleX float64, pixelScaleY float64) CircleState {
	circles := gocv.NewMat()
	res := 0.0
	dpHoughCircles := 1.0
	thresholdHoughCircles := 90.0
	accumulatorThresholdHoughCircles := 10.0

	innerOverdraw := int(float64(1) * math.Min(pixelScaleX, pixelScaleY))
	circleSize := int(float64(5) * math.Min(pixelScaleX, pixelScaleY))
	minDist := (float64(15) * math.Min(pixelScaleX, pixelScaleY))

	gocv.HoughCirclesWithParams(
		img,
		&circles,
		gocv.HoughGradient,
		dpHoughCircles,                   // dp
		minDist,                          //float64(croppedMatGray.Rows()/50), // minDist
		thresholdHoughCircles,            // param1
		accumulatorThresholdHoughCircles, // param2
		circleSize,                       // minRadius
		circleSize,                       // maxRadius
	)

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			r_out := int(math.Ceil(math.Sqrt((float64(r) * float64(r) * 2))))
			gocv.Circle(&img, image.Pt(x, y), r_out, color.RGBA{255, 255, 255, 0}, r_out-r+int(float64(innerOverdraw)*2))
			gocv.Circle(&img, image.Pt(x, y), r-int(float64(innerOverdraw)*1.4/2), color.RGBA{255, 255, 255, 0}, int(float64(innerOverdraw)*2))
			// log.Println("Radius", r)

			/**
			==========================================
			für NN bitte zentrieren und ausschneiden
			*/

			stepSize := r / 10
			index := 0

			for s := stepSize; s < r; s += stepSize {
				img2 := img.Clone()
				sx := s + stepSize*3

				// log.Println("center", image.Pt(x, y))
				gocv.Circle(&img2, image.Pt(x, y), sx, color.RGBA{255, 255, 255, 0}, r_out-sx-int(float64(innerOverdraw)))
				gocv.Threshold(img2, &img2, 150, 255, gocv.ThresholdBinary)

				/*
					avg := 1/(img2.Sum().Val1/(float64(img2.Cols())*float64(img2.Rows()))) - 0.00392156862745098
					log.Println("avg* img2", s, avg)
					v += avg / math.Pow10(index)
				*/
				avg := 1 - (img2.Sum().Val1/(float64(img2.Cols())*float64(img2.Rows())))/255.0
				res += avg / math.Pow10(index)

				//
				img2.Close()
				index++

			}

			res *= math.Pow10(4)

			// avg := img.Sum().Val1 / (float64(img.Cols()) * float64(img.Rows()))
			// log.Println("avg v>>>>", res > 1.0)
		}
	}
	/*
		me.wnd.IMShow(img)
		me.wnd.WaitKey(0)
	*/
	cImg := img.Clone()

	return CircleState{
		Marked: res > 0.001,
		Found:  true,
		Value:  res,
		Image:  cImg,
	}
}

func (me *Extractor) getPaginations() {
	rows, err := me.db.Query("SELECT pagination_id FROM papervote_optical where /*edited_marks not like '%W%' and pagination_id not like '40%' and*/  pagination_id  like '22%' limit 10000")
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var names []string
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			break
		}
		data_rows, err2 := me.db.Query("SELECT  replace(data,' ','+') data FROM papervote_optical_data where pagination_id = " + name)
		if err2 != nil {
			fmt.Println("Error Found:", err2)
			break
		}
		for data_rows.Next() {
			var data string
			var marks []bool
			var xmarks []bool
			var fmarks []float64
			err3 := data_rows.Scan(&data)

			b64data := data[strings.Index(data, ",")+1:]
			if err3 != nil {
				fmt.Println("Error Found:", err3)
				return
			}

			dst := make([]byte, base64.StdEncoding.DecodedLen(len(b64data)))
			_, errX := base64.StdEncoding.Decode(dst, []byte(b64data))
			if errX != nil {
				fmt.Println("Error Found:", errX)
				return
			}

			orerr := os.WriteFile(name+".jpg", dst, 0644)
			if orerr != nil {
				fmt.Println("Error Found:", orerr)
				return
			}

			sql := `
			with roi as (
				select 
					sz_rois.id roi_id,
					sz_rois.name roi_name,
					sz_rois.x   roi_x,
					sz_rois.y   roi_y,
					sz_rois.width roi_width,
					sz_rois.height roi_height,
					sz_rois.item_height roi_item_height,
					sz_rois.item_cap_y roi_item_cap_y,
					sz_page_sizes.width page_width,
					sz_page_sizes.height page_height
				from 
					stimmzettel 
					join stimmzettel_roi 
						on stimmzettel_roi.stimmzettel_id = stimmzettel.id
					join sz_rois 
						on stimmzettel_roi.sz_rois_id = sz_rois.id
					join sz_to_region 
						on sz_to_region.id_sz = stimmzettel.id
					join sz_titel_regions 
						on  sz_titel_regions.id = sz_to_region.id_sz_titel_regions
					join sz_to_page_sizes 
						on sz_to_page_sizes.id_sz = stimmzettel.id
					join sz_page_sizes 
						on  sz_to_page_sizes.id_sz_page_sizes = sz_page_sizes.id
	
			)
				select 
					roi.roi_x,
					roi.roi_y,
					roi.roi_item_height,
					roi.roi_item_cap_y,
					roi.page_width,
					roi.page_height,
					view_papervote_optical_result_ballotpaper.marked,
					rank() over (
						partition by pagination_id,
						sz_rois_id
						order by result_index
					) roi_pos 
				from 
					view_papervote_optical_result_ballotpaper 
					join roi on roi.roi_id = view_papervote_optical_result_ballotpaper.sz_rois_id
 				where pagination_id=`

			roi_rows, err_roi_rows := me.db.Query(sql + name)
			if err_roi_rows != nil {
				fmt.Println("Error Found:", err_roi_rows)
				break
			}

			img := gocv.IMRead(name+".jpg", gocv.IMReadGrayScale)
			if img.Empty() {
				fmt.Printf("Failed to read image: %s\n", name+".jpg")
				os.Exit(1)
			}

			currentImageName = name

			log.Println("============== ", name, " ==============")
			idx := 0
			for roi_rows.Next() {

				var roi_pos int
				var roi_x int
				var roi_y int
				var roi_item_height int
				var page_width int
				var page_height int
				var roi_item_cap_y float64
				var marked string
				roi_rows_err3 := roi_rows.Scan(&roi_x, &roi_y, &roi_item_height, &roi_item_cap_y, &page_width, &page_height, &marked, &roi_pos)
				log.Println("marked ", marked, roi_pos, roi_x, roi_y, roi_item_height, roi_item_cap_y, roi_rows_err3)

				log.Println("Ratio 1", float64(img.Cols())/float64(img.Rows()), float64(page_width)/float64(page_height))

				var goalRatioX float64
				var goalRatioY float64

				goalRatioX = 1.0
				goalRatioY = 1.0

				if (math.Abs(math.Round(((float64(img.Cols())/float64(img.Rows()))/(float64(page_width)/float64(page_height))-1)*10) / 10)) != 0 {
					log.Println("RX", math.Abs(math.Round(((float64(img.Cols())/float64(img.Rows()))/(float64(page_width)/float64(page_height))-1)*10)/10))
					/*
						if (float64(img.Cols()) / float64(img.Rows())) < (float64(page_width) / float64(page_height)) {
							goalRatioX = 1 / (float64(page_width) / float64(page_height))
						} else {
							goalRatioY = 1 / (float64(img.Cols()) / float64(img.Rows()))
						}


					*/

					originalFactor := (float64(page_width) / float64(page_height))
					imageFactor := (float64(img.Cols()) / float64(img.Rows()))

					log.Println("***************")
					log.Println(page_width, page_height)
					log.Println(img.Cols(), img.Rows())

					log.Println("originalFactor", originalFactor)
					log.Println("imageFactor", imageFactor)

					log.Println("X")

					if imageFactor > originalFactor {
						goalRatioY = (imageFactor / originalFactor)
					} else {
						goalRatioX = 1 / (imageFactor / originalFactor)
					}

					if imageFactor > originalFactor {

					}

					log.Println("resultFactor", ((float64(img.Cols()) * goalRatioX) / (float64(img.Rows()) * goalRatioY)))

					log.Println("***************")

				}

				if idx == 0 {
					log.Println(goalRatioX, int(float64(img.Cols())*goalRatioX))
					log.Println(goalRatioY, int(float64(img.Rows())*goalRatioY))
					/*
						me.wnd.IMShow(img)
						me.wnd.WaitKey(0)
					*/

					/*
						if "22004" == name {
							me.wnd.IMShow(img)
							me.wnd.WaitKey(0)
						}
					*/

					gocv.Resize(
						img,
						&img,
						image.Point{
							int(float64(img.Cols()) * goalRatioX),
							int(float64(img.Rows()) * goalRatioY),
						},
						0,
						0,
						gocv.InterpolationArea,
					)
					log.Println("Ratio 2", float64(img.Cols())/float64(img.Rows()), float64(page_width)/float64(page_height))

				}

				scaleX := float64(img.Cols()) / float64(page_width)
				scaleY := float64(img.Rows()) / float64(page_height)

				y := (roi_y + int((float64(roi_item_height)+roi_item_cap_y)*float64(roi_pos-1)))
				rect := image.Rect(
					int(float64(roi_x)*scaleX),
					int(float64(y)*scaleY),
					int(float64(roi_x)*scaleX+float64(roi_item_height)*scaleY),
					int(float64(y)*scaleY+float64(roi_item_height)*scaleY),
				)
				croppedMat := img.Region(rect)
				/*
					me.wnd.IMShow(croppedMat)
					me.wnd.WaitKey(0)
				*/
				currentRoiRect = idx
				res := me.circles(croppedMat, scaleX, scaleY)
				idx++
				// log.Println("res", res.Found, res.Marked, res.Value, name, page_width, page_height)
				/*if idx == 2 {
					me.wnd.IMShow(croppedMat)
					me.wnd.WaitKey(0)
					//os.Exit(1)
				}*/
				if !res.Found {
					/*
						me.wnd.IMShow(img)
						me.wnd.WaitKey(0)
						me.wnd.IMShow(croppedMat)
						me.wnd.WaitKey(0)
					*/
				} else {

					xmarks = append(xmarks, marked == "X")
					marks = append(marks, res.Marked)
					fmarks = append(fmarks, res.Value)

					gocv.IMWrite(
						fmt.Sprintf("data/%s/%s.%d.%d.jpg", marked, name, roi_pos, roi_x),
						res.Image)
					res.Image.Close()
				}
				croppedMat.Close()

			}

			failure := false
			for i := 0; i < len(xmarks); i++ {
				if xmarks[i] != marks[i] {
					failure = true
				}
			}

			log.Println("index", "alt", "neu", "wert")
			for i := 0; i < len(xmarks); i++ {
				log.Println(i, xmarks[i], marks[i], fmarks[i])
			}

			/*
				log.Println("marks NNNN*", xmarks)
				log.Println("marks NNNN#", marks)

				log.Println("marks NNNNN", fmarks)
			*/
			if failure {
				me.wnd.IMShow(img)
				me.wnd.WaitKey(0)
			}

			img.Close()
			roi_rows.Close()
			os.Remove(name + ".jpg")
		}
		data_rows.Close()
		names = append(names, name)
	}

	defer rows.Close()

}
