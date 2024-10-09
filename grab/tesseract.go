package grab

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)

func fileformatBytes(img gocv.Mat) []byte {
	buffer, err := gocv.IMEncodeWithParams(gocv.PNGFileExt, img, []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes()
}

type SortIndexDistance struct {
	Index    int
	Distance float64
}

func (this *GrabcameraClass) uniqueCharacters(str string) string {
	charSet := make(map[rune]bool)

	for _, char := range str {
		charSet[char] = true
	}

	keys := make([]rune, 0, len(charSet))
	res := ""
	for k := range charSet {
		keys = append(keys, k)
		if (k >= 65 && k <= 90) || (k >= 97 && k <= 122) || (k >= 48 && k <= 57) || (k == 252) || (k == 246) || (k == 228) {
			res += string(k)
		}
	}

	return res
}

func (this *GrabcameraClass) printableCharacters(str string) string {
	res := ""
	for _, char := range str {
		if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || (char >= 48 && char <= 57) || (char == 252) || (char == 246) || (char == 228) {
			res += string(char)
		}
	}
	return res
}

func (this *GrabcameraClass) tesseract(img gocv.Mat, currentOCRChannel chan string) structs.TesseractReturnType {

	start := time.Now()

	result := structs.TesseractReturnType{}
	result.Point = image.Point{0, 0}
	result.Title = ""
	result.IsCorrect = false
	result.PageRois = []structs.DocumentConfigurationPageRoi{}
	result.Pagesize = structs.DocumentConfigurationPageSize{}
	result.CircleSize = 1
	result.CircleMinDistance = 100
	result.Marks = []structs.CheckMarks{}

	client := gosseract.NewClient()
	defer client.Close()

	if this.globals.TesseractPrefix != "" {
		client.SetTessdataPrefix(this.globals.TesseractPrefix)
	}
	client.SetLanguage("deu")

	documentConfigurations := this.documentConfigurations

	/*
		for i := 0; i < len(documentConfigurations); i++ {

			log.Println("ID", documentConfigurations[i].Titles[0])
			log.Println("PageSize", fmt.Sprintf("%.3f", float64(documentConfigurations[i].Pagesize.Width)/float64(documentConfigurations[i].Pagesize.Height)))
			log.Println("PaperSize", fmt.Sprintf("%.3f", float64(img.Cols())/float64(img.Rows())))

		}
	*/

	indexSort := []SortIndexDistance{}
	for i := 0; i < len(documentConfigurations); i++ {

		/*
			itm := SortIndexDistance{
				int(math.Abs(
					math.Round(float64(img.Cols())/float64(img.Rows())*100) -
						math.Round(float64(documentConfigurations[i].Pagesize.Width)/float64(documentConfigurations[i].Pagesize.Height)*100),
				)), i,
			}
		*/
		itm := SortIndexDistance{}
		itm.Distance = math.Abs(float64(img.Cols())/float64(img.Rows()) - float64(documentConfigurations[i].Pagesize.Width)/float64(documentConfigurations[i].Pagesize.Height))
		itm.Index = i

		indexSort = append(indexSort, itm)
	}
	log.Println("Unsorted", indexSort)

	sort.SliceStable(indexSort, func(i, j int) bool {
		return indexSort[i].Distance < indexSort[j].Distance
	})

	log.Println("Sorted", indexSort)

	for i := 0; i < len(indexSort); i++ {

		currentCocumentConfigurations := documentConfigurations[indexSort[i].Index]
		/*

			log.Println("PaperSize", fmt.Sprintf("%.3f", float64(img.Cols())/float64(img.Rows())))
		*/

		log.Println(
			"PR",
			i,
			indexSort[i].Index,
			indexSort[i].Distance,
			math.Round(float64(img.Cols())/float64(img.Rows())*10)-
				math.Round(float64(currentCocumentConfigurations.Pagesize.Width)/float64(currentCocumentConfigurations.Pagesize.Height)*10),
			currentCocumentConfigurations.Titles[0],
		)

		//if math.Round(float64(img.Cols())/float64(img.Rows())*100) == math.Round(float64(currentCocumentConfigurations.Pagesize.Width)/float64(currentCocumentConfigurations.Pagesize.Height)*100) {

		result.CircleSize = currentCocumentConfigurations.CircleSize
		result.CircleMinDistance = currentCocumentConfigurations.CircleMinDistance
		result.Pagesize = currentCocumentConfigurations.Pagesize

		X := currentCocumentConfigurations.TitleRegion.X * img.Cols() / result.Pagesize.Width
		Y := currentCocumentConfigurations.TitleRegion.Y * img.Rows() / result.Pagesize.Height
		W := currentCocumentConfigurations.TitleRegion.Width * img.Cols() / result.Pagesize.Width
		H := currentCocumentConfigurations.TitleRegion.Height * img.Rows() / result.Pagesize.Height

		croppedMat := img.Region(image.Rect(X, Y, W+X, H+Y))

		if false {
			gocv.IMWrite(fmt.Sprintf("tesseract_%d_%d%d.png", i, X, Y), croppedMat)
		}

		if croppedMat.Empty() {
			croppedMat.Close()
			return result
		}

		if this.globals.ShowImage == 601 {
			pImage := croppedMat.Clone()
			this.pipeUIImage(pImage)
			pImage.Close()
		}

		whiteListCharactes := ""
		for j := 0; j < len(currentCocumentConfigurations.Titles); j++ {
			whiteListCharactes += currentCocumentConfigurations.Titles[j]
		}
		// log.Println("SetWhitelist",whiteListCharactes)
		// this.uniqueCharacters(whiteListCharactes)

		// client.SetWhitelist(whiteListCharactes)

		// imgColor := gocv.NewMat()
		// gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
		client.SetWhitelist(this.uniqueCharacters(whiteListCharactes))
		smaller := this.ResizeMat(croppedMat.Clone(), croppedMat.Cols()/this.globals.TesseractScale, croppedMat.Rows()/this.globals.TesseractScale)

		seterror := client.SetImageFromBytes(fileformatBytes(smaller))
		if seterror != nil {
			// fmt.Println(seterror)
			return result
		}
		out, herr := client.GetBoundingBoxes(3)
		if herr != nil {
			// fmt.Println(herr)
			croppedMat.Close()
			return result
		} else {
			if i == 0 {

			}
			searchFor := ""
			if true {
				for j := 0; j < len(out); j++ {
					searchFor += " " + out[j].Word
				}
			}

			searchFor = this.printableCharacters(searchFor)

			if len(currentOCRChannel) == cap(currentOCRChannel) {
				txt, _ := <-currentOCRChannel
				if false {
					log.Println("OCR", txt)
				}
			}
			currentOCRChannel <- searchFor

			// fmt.Println("searchFor %s %d",searchFor , len(currentCocumentConfigurations.Titles))
			for j := 0; j < len(currentCocumentConfigurations.Titles); j++ {
				distance := levenshtein.ComputeDistance(searchFor, this.printableCharacters(currentCocumentConfigurations.Titles[j]))
				errorRate := float64(distance) /*- float64(len( currentCocumentConfigurations.Titles[j])-len(searchFor)))*/ / float64(len(this.printableCharacters(currentCocumentConfigurations.Titles[j])))

				if strings.Contains(searchFor, this.printableCharacters(currentCocumentConfigurations.Titles[j])) {
					errorRate = 0
				}
				log.Println("ErrorRate:", errorRate, "Search:", searchFor, this.printableCharacters(currentCocumentConfigurations.Titles[j]))
				if errorRate < 0.5 {
					result.Title = currentCocumentConfigurations.Titles[j]

					//title = out[0].Word
					drawContours := gocv.NewPointsVector()
					contour := gocv.NewPointVectorFromPoints([]image.Point{
						out[0].Box.Min,
						image.Point{out[0].Box.Max.X, out[0].Box.Min.Y},
						out[0].Box.Max,
						image.Point{out[0].Box.Min.X, out[0].Box.Max.Y}})
					drawContours.Append(contour)
					gocv.DrawContours(&croppedMat, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)
					result.Point = image.Point{currentCocumentConfigurations.TitleRegion.X, currentCocumentConfigurations.TitleRegion.Y}

					if false {
						fmt.Sprintf("ocr %s %d %d %d", time.Since(start), croppedMat.Cols(), croppedMat.Rows(), os.Getpid())
					}
					result.PageRois = currentCocumentConfigurations.Rois
					croppedMat.Close()
					drawContours.Close()
					smaller.Close()
					return result
				}

			}

		}
		croppedMat.Close()
		smaller.Close()
		//}

	}

	return result

}
