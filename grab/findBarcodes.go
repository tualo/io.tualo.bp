
package grab

import (
	"log"
	"image"
	"github.com/bieber/barcode"
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)

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
	
	for _, s := range symbols {
		syms = append(syms,structs.BarcodeSymbol{Type:s.Type.Name(),Data:s.Data,Quality:s.Quality,Boundary:s.Boundary})
		if false {
			log.Println("BarcodeSymbol",s.Type.Name(),s.Data,s.Quality,s.Boundary)
		}
	}
	smaller.Close()
	return syms
}