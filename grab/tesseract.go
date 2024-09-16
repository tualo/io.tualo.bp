package grab

import (
	"fmt"
	"time"
	"image"
	"os"
	"strings"
	"image/color"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"github.com/agnivade/levenshtein"
	structs "io.tualo.bp/structs"
	"log"
)
func fileformatBytes(img gocv.Mat) []byte {
	buffer, err :=gocv.IMEncodeWithParams(gocv.PNGFileExt, img, []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes(	)
}


func (this *GrabcameraClass) uniqueCharacters(str string) string {
    charSet := make(map[rune]bool)

    for _, char := range str {
        charSet[char] = true
    }

	keys := make([]rune, 0, len(charSet))
	res:=""
	for k := range charSet {
		keys = append(keys, k)
		if (k>=65 && k<=90) || (k>=97 && k<=122) || (k>=48 && k<=57) {
			res+=string(k)
		}
	}

	return res
}

func (this *GrabcameraClass) printableCharacters(str string) string {
	res:=""
	for _, char := range str {
		if (char>=65 && char<=90) || (char>=97 && char<=122) || (char>=48 && char<=57) {
			res+=string(char)
		}
	}
	return res
}

func (this *GrabcameraClass) tesseract(img gocv.Mat, currentOCRChannel chan string) (structs.TesseractReturnType) {

	start := time.Now()



	result:=structs.TesseractReturnType{}
	result.Point=image.Point{0,	0}
	result.Title=""
	result.IsCorrect=false
	result.PageRois=[]structs.DocumentConfigurationPageRoi{}
	result.Pagesize=structs.DocumentConfigurationPageSize{}
	result.CircleSize=1
	result.CircleMinDistance=100
	result.Marks=[]structs.CheckMarks{}
	
	client := gosseract.NewClient()
	defer client.Close()

	if this.globals.TesseractPrefix != "" {
		client.SetTessdataPrefix(this.globals.TesseractPrefix)
	}
	client.SetLanguage("deu")

	documentConfigurations := this.documentConfigurations

	if false {
		fmt.Println("tesseract",documentConfigurations, img.Cols(), img.Rows())
	}
	for i := 0; i < len(documentConfigurations); i++ {
		

		result.CircleSize=documentConfigurations[i].CircleSize
		result.CircleMinDistance=documentConfigurations[i].CircleMinDistance
		result.Pagesize=documentConfigurations[i].Pagesize
		


		X := documentConfigurations[i].TitleRegion.X * img.Cols() / result.Pagesize.Width
		Y := documentConfigurations[i].TitleRegion.Y * img.Rows() / result.Pagesize.Height
		W := documentConfigurations[i].TitleRegion.Width * img.Cols() / result.Pagesize.Width
		H := documentConfigurations[i].TitleRegion.Height * img.Rows() / result.Pagesize.Height

		croppedMat := img.Region(image.Rect(X, Y, W+X, H+Y))

		if false {

			gocv.IMWrite(fmt.Sprintf( "tesseract_%d_%d%d.png",i,X,Y), croppedMat)
		}

		if croppedMat.Empty() {
			croppedMat.Close()
			return result
		}

		whiteListCharactes := ""
		for j := 0; j < len(documentConfigurations[i].Titles); j++ {
			whiteListCharactes += documentConfigurations[i].Titles[j]
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
		}else{
			if i==0 {	
				
			}
			searchFor:=""
			if true {
				for j := 0; j < len(out); j++ {
					searchFor += " " +out[j].Word
				}
			}

			searchFor = this.printableCharacters(searchFor)


			if len(currentOCRChannel)==cap(currentOCRChannel) {
				txt,_:=<-currentOCRChannel
				if false {
					log.Println("OCR",txt)
				}
			}
			currentOCRChannel <- searchFor

			// fmt.Println("searchFor %s %d",searchFor , len(documentConfigurations[i].Titles))
			for j := 0; j < len(documentConfigurations[i].Titles); j++ {
				distance := levenshtein.ComputeDistance(searchFor, this.printableCharacters(documentConfigurations[i].Titles[j]))
				errorRate:=float64(distance) /*- float64(len( documentConfigurations[i].Titles[j])-len(searchFor)))*/ /	float64(len( documentConfigurations[i].Titles[j]))

				if true {
					fmt.Printf("OCR ======== \nSearchFor: *%s*\nTitel:*%s*\nTitelNr:%d\nDistance:%d\nprintChars:%s.\nError:%d\n", 
						searchFor, 
						documentConfigurations[i].Titles[j], 
						len( documentConfigurations[i].Titles[j]), 
						distance,
						this.printableCharacters(documentConfigurations[i].Titles[j]),
						errorRate,
					)
			
				}
				if strings.Contains(searchFor, this.printableCharacters(documentConfigurations[i].Titles[j])) {
					fmt.Println("Contains")
					errorRate=0
				}
				if errorRate < 0.3 {
					result.Title=documentConfigurations[i].Titles[j]
					
					//title = out[0].Word
					drawContours := gocv.NewPointsVector()
					contour:= gocv.NewPointVectorFromPoints([]image.Point{
						out[0].Box.Min,
						image.Point{out[0].Box.Max.X, out[0].Box.Min.Y},
						out[0].Box.Max,
						image.Point{out[0].Box.Min.X, out[0].Box.Max.Y} 			})
					drawContours.Append(contour)
					gocv.DrawContours(&croppedMat, drawContours, -1, color.RGBA{0, 255, 0, 0}, 2)
					result.Point = image.Point{documentConfigurations[i].TitleRegion.X, documentConfigurations[i].TitleRegion.Y}

					
					if false {
						fmt.Sprintf("ocr %s %d %d %d",time.Since(start),croppedMat.Cols(),croppedMat.Rows(), os.Getpid() ) 
					}
					result.PageRois=documentConfigurations[i].Rois
					croppedMat.Close()
					drawContours.Close()
					smaller.Close()
					return result
				}

			}

			
		
		}
		croppedMat.Close()
		smaller.Close()
	}

	return result

}