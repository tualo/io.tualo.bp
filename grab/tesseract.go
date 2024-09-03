package grab

import (
	"fmt"
	"time"
	"image"
	"os"
	"image/color"
	"github.com/otiai10/gosseract/v2"
	"gocv.io/x/gocv"
	"github.com/agnivade/levenshtein"
	structs "io.tualo.bp/structs"
)
func fileformatBytes(img gocv.Mat) []byte {
	buffer, err :=gocv.IMEncodeWithParams(gocv.PNGFileExt, img, []int{gocv.IMWriteJpegQuality, 100})
	if err != nil {
		return nil
	}
	return buffer.GetBytes(	)
}



func (this *GrabcameraClass) tesseract(img gocv.Mat) (structs.TesseractReturnType) {

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

		if croppedMat.Empty() {
			croppedMat.Close()
			return result
		}



		// imgColor := gocv.NewMat()
		// gocv.CvtColor(croppedMat, &imgColor, gocv.ColorGrayToBGR)
		// client.SetWhitelist("EinzelhandelEnergie")
		smaller := this.ResizeMat(croppedMat.Clone(), croppedMat.Cols()/this.globals.TesseractScale, croppedMat.Rows()/this.globals.TesseractScale)

		seterror := client.SetImageFromBytes(fileformatBytes(smaller))
		if seterror != nil {
			fmt.Println(seterror)
			return result
		}
		out, herr := client.GetBoundingBoxes(3)
		if herr != nil {
			fmt.Println(herr)
			croppedMat.Close()
			return result
		}else{
			if i==0 {	
				if false {
					gocv.IMWrite("tesseract.png", croppedMat)
				}
			}
			searchFor:=""
			if true {
				for j := 0; j < len(out); j++ {
					searchFor += " " +out[j].Word
				}
			}
			fmt.Println("searchFor %s %d",searchFor , len(documentConfigurations[i].Titles))
			for j := 0; j < len(documentConfigurations[i].Titles); j++ {
				distance := levenshtein.ComputeDistance(searchFor, documentConfigurations[i].Titles[j])
				errorRate:=float64(distance) /*- float64(len( documentConfigurations[i].Titles[j])-len(searchFor)))*/ /	float64(len( documentConfigurations[i].Titles[j]))

				if true {
					fmt.Printf("The distance between:  *%s*  *%s*  is %d %d. \n", 
					searchFor, 
					documentConfigurations[i].Titles[j], 
					len( documentConfigurations[i].Titles[j]), 
					distance,

				)
				/*
				fmt.Printf("Len diff %.2f\nRatio %.2f\n",
					float64(len( documentConfigurations[i].Titles[j])-len(searchFor)),
					(float64(distance) - float64(len( documentConfigurations[i].Titles[j])-len(searchFor))) /	float64(len( documentConfigurations[i].Titles[j]))			,
				)*/
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