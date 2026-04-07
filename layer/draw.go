package layer

import (
	"image"
	"log"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

func (me *Layer) Draw() (cv.Mat, int) {
	var lastMat cv.Mat
	imgFound := 0

	mutex.Lock()

	item, present := me.drawLayers["main"]
	if !present {
		mutex.Unlock()
		return lastMat, 0
	}
	lastMat = item.Mat

	imgFound++

	for key, item := range me.drawLayers {
		if key != "main" {
			if !item.DoNotDraw {

				if item.BasedOn == "" {
					gocv.AddWeighted(
						item.Mat.Get(),
						item.Weighted_Alpha,
						lastMat.Get(),
						item.Weighted_Beta,
						item.Weighted_Gamma,
						lastMat.GetPointer(),
					)
				} else {

					log.Println(">", item.Title, item.BasedOn)
					playground := gocv.Zeros(lastMat.GetPointer().Rows(), lastMat.GetPointer().Cols(), lastMat.GetPointer().Type())
					gocv.WarpPerspective(item.Mat.Get(),
						&playground,
						me.drawLayers[item.BasedOn].Mat.Get(),
						image.Point{
							lastMat.GetPointer().Cols(),
							lastMat.GetPointer().Rows(),
						})

					gocv.AddWeighted(
						playground,
						item.Weighted_Alpha,
						lastMat.Get(),
						item.Weighted_Beta,
						item.Weighted_Gamma,
						lastMat.GetPointer(),
					)
					playground.Close()
				}

				imgFound++
			}
		}
	}
	log.Println("Draw", "C")

	mutex.Unlock()
	return lastMat, imgFound
}
