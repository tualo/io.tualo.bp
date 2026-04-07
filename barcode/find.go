package barcode

import (
	"image"
	"image/color"
	"log"

	"github.com/bieber/barcode"
	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/layer"
	"tualo.de/deep-test/structs"
)

type BCScanner struct {
	scanner *barcode.ImageScanner
}

var static *BCScanner

func Static() *BCScanner {
	if static == nil {
		static = &BCScanner{}
		static.scanner = barcode.NewScanner()
		static.scanner.SetEnabledAll(false)
		static.scanner.SetEnabledSymbology(barcode.Code39, true)
		static.scanner.SetMinLengthSymbology(barcode.Code128, 6)
		static.scanner.SetMaxLengthSymbology(barcode.Code128, 10)
		static.scanner.SetEnabledSymbology(barcode.Code128, true)
	}
	return static
}

func (me *BCScanner) Text(
	img *cv.Mat,
	layerName string,
	baseLayerName string,
	text string,
	posX float64,
	posY float64,
	color color.RGBA,
) {
	newmat := cv.NewFromMat("BCText", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))
	gocv.PutText(newmat.GetPointer(), text,
		image.Pt(
			int(float64(newmat.Cols())*posX),
			int(float64(newmat.Cols())*posY),
		), gocv.FontHersheyPlain,
		8.2,
		color,
		22,
	)
	layer.Static().Set(layerName, newmat, false, baseLayerName, 0.5, 1.0, 0.0)
	newmat.Close()
}

func (me *BCScanner) Find(img *cv.Mat) ([]structs.BarcodeSymbol, bool) {
	syms := []structs.BarcodeSymbol{}
	if img.GetPointer().Empty() {
		return syms, false
	}
	barcodeScale := 1

	smaller := gocv.NewMat()
	gocv.CvtColor(img.Get(), &smaller, gocv.ColorBGRToGray)
	if true {
		if smaller.Cols() > 800 {
			gocv.GaussianBlur(smaller, &smaller, image.Point{5, 5}, 0, 0, gocv.BorderDefault)
			gocv.Resize(smaller, &smaller, image.Point{smaller.Cols() / barcodeScale, smaller.Rows() / barcodeScale}, 0, 0, gocv.InterpolationArea)
		}
	}

	// preparegozxing.NewBinaryBitmapFromImage(img)

	// gocv.IMWrite("page.jpg", smaller)

	symbols, err := me.scanner.ScanMat(&smaller)
	if err != nil {
		// panic(err)
		log.Println("Error scanning for barcodes:", err)
		smaller.Close()
		return syms, false
	}

	log.Println("Symbols found:", len(symbols))
	for _, s := range symbols {
		syms = append(syms, structs.BarcodeSymbol{Type: s.Type.Name(), Data: s.Data, Quality: s.Quality, Boundary: s.Boundary})

		rect := image.Rectangle{}
		minX := img.Cols()
		minY := img.Rows()
		maxX := 0
		maxY := 0
		for i := 0; i < len(s.Boundary); i++ {

			if s.Boundary[i].X < minX {
				minX = s.Boundary[i].X
			}
			if s.Boundary[i].Y < minY {
				minY = s.Boundary[i].Y
			}
			if s.Boundary[i].X > maxX {
				maxX = s.Boundary[i].X
			}
			if s.Boundary[i].Y > maxY {
				maxY = s.Boundary[i].Y
			}
		}
		rect.Min.Y = minY
		rect.Min.X = minX
		rect.Max.X = maxX
		rect.Max.Y = maxY

		if true {
			me.Text(img, "BCText", "inverseMat", s.Data, 0.1, 0.1, color.RGBA{255, 100, 100, 255})
			// log.Println("BarcodeSymbol", s.Type.Name(), s.Data, s.Quality, s.Boundary, rect)

			newmat := cv.NewFromMat("BCText", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))
			gocv.Rectangle(
				newmat.GetPointer(),
				rect,
				color.RGBA{255, 100, 100, 255},
				19,
			)
			layer.Static().Set("rect_bc", newmat, false, "inverseMat", 0.5, 1.0, 0.0)
			newmat.Close()

		}
	}
	smaller.Close()
	return syms, len(syms) > 0
}
