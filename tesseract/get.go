package tesseract

import (
	"image"
	"image/color"
	"log"
	"math"
	"sort"
	"strings"

	"github.com/agnivade/levenshtein"
	"gocv.io/x/gocv"
	"tualo.de/deep-test/api"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/layer"
	"tualo.de/deep-test/structs"
)

func (me *Tesseract) DetectBallotpaperType(img *cv.Mat) int {
	cvx := img.Clone("DetectBallotpaperType")
	result := -1

	// find best paper format
	configs := api.Static().GetTitleRegionsConfig()

	indexSort := []structs.SortIndexDistance{}
	for i := 0; i < len(configs); i++ {
		itm := structs.SortIndexDistance{}

		log.Println("configs  ", i, configs[i], " ", configs[i].Page_width, configs[i].Page_height)

		itm.Distance = math.Abs(float64(img.Cols())/float64(img.Rows()) - float64(configs[i].Page_width)/float64(configs[i].Page_height))
		itm.Index = i
		// chekc if it is not nan
		if math.IsNaN(itm.Distance) {
			log.Println("Skipping NaN distance for index", i)
			continue
		}
		indexSort = append(indexSort, itm)
	}
	sort.SliceStable(indexSort, func(i, j int) bool {
		return indexSort[i].Distance < indexSort[j].Distance
	})

	mapHash := make(map[float64]map[int][]string)
	for i := 0; i < len(indexSort); i++ {
		itm := indexSort[i]
		area := configs[itm.Index].Area
		_, prs := mapHash[itm.Distance]
		if !prs {
			mapHash[itm.Distance] = make(map[int][]string)
		}
		log.Println("itm.Distance  ", itm.Distance)

		_, prs2 := mapHash[itm.Distance][area]
		if !prs2 {

			mapHash[itm.Distance][area] = []string{}
		}

		mapHash[itm.Distance][area] = append(mapHash[itm.Distance][area], configs[itm.Index].OCRTitle)
	}

	newmat := cv.NewFromMat("Tesseract", gocv.Zeros(img.Rows(), img.Cols(), img.Type()))

	for _, areas := range mapHash {
		for area, titles := range areas {
			// log.Println("area", area)
			if len(titles) == 0 {
				continue
			}

			configIndex, exists := api.Static().AllOCRTitlesInv[titles[0]]
			if !exists || configIndex >= len(configs) {
				continue
			}

			config := configs[configIndex]

			if area == config.Area {
				scale := structs.Scale{
					PixelScaleX: float64(cvx.Cols()) / float64(config.Page_width),
					PixelScaleY: float64(cvx.Rows()) / float64(config.Page_height),
				}

				X := float64(config.Roi_x) * scale.PixelScaleX
				Y := float64(config.Roi_y) * scale.PixelScaleY
				W := float64(config.Roi_width) * scale.PixelScaleX
				H := float64(config.Roi_height) * scale.PixelScaleY

				gocv.Rectangle(
					newmat.GetPointer(),
					image.Rect(int(X), int(Y), int(X+W), int(Y+H)),
					color.RGBA{255, 0, 255, 100},
					15,
				)

				croppedMat := cv.NewFromMat("cropped tess", img.GetPointer().Region(image.Rect(int(X), int(Y), int(W+X), int(H+Y))))
				seterror := me.Client.SetImageFromBytes(me.fileformatBytes(&croppedMat))
				croppedMat.Close()
				if seterror != nil {
					cvx.Close()
					newmat.Close()
					return -1
				}
				whiteListCharactes := config.OCRTitle
				me.Client.SetWhitelist(me.uniqueCharacters(whiteListCharactes))

				out, herr := me.Client.GetBoundingBoxes(3)
				log.Println("search", config.OCRTitle)
				if herr != nil {
					cvx.Close()
					newmat.Close()
					return -1
				} else {
					searchFor := ""
					if true {
						for j := 0; j < len(out); j++ {
							searchFor += " " + out[j].Word
						}
					}

					searchFor = me.printableCharacters(searchFor)

					for j := 0; j < len(titles); j++ {
						log.Println("Comparing searchFor:", searchFor, " with title:", titles[j])
						distance := levenshtein.ComputeDistance(searchFor, me.printableCharacters(titles[j]))
						errorRate := float64(distance) / float64(len(me.printableCharacters(titles[j])))

						if strings.Contains(searchFor, me.printableCharacters(titles[j])) {
							errorRate = 0
						}

						log.Println(
							"ErrorRate:",
							errorRate,
							"Search:",
							searchFor,
							me.printableCharacters(titles[j]),
							"T",
						)
						log.Println(
							titles[j],
							distance,
							api.Static().AllOCRTitlesInv,
							api.Static().AllOCRTitlesInv[titles[j]],
						)

						if errorRate < 2.95 {

							gocv.Rectangle(
								newmat.GetPointer(),
								image.Rect(int(X), int(Y), int(X+W), int(Y+H)),
								color.RGBA{0, 255, 0, 100},
								30,
							)

							if resultIndex, exists := api.Static().AllOCRTitlesInv[titles[j]]; exists {
								result = resultIndex
							}
							log.Println("found result", result)
							cvx.Close()
							newmat.Close()
							return result
						}
					}
				}
			}

		}

	}
	layer.Static().Set("tess rois", newmat, false, "inverseMat", 1.9, 1, 0.0)
	newmat.Close()

	cvx.Close()
	return result
}
