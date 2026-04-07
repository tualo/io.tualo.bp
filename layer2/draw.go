package layer2

import (
	"image"
	"log"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

func (me *Layer2) Draw(width int, height int) (cv.Mat, int) {
	var lastMat cv.Mat
	var invMat cv.Mat

	imgFound := 0

	main, exists := me.layers["main"]
	if !exists {
		return lastMat, 0
	}

	if len(main) > 0 {

		inv, existsInv := me.layers["inverseMat"]
		if existsInv {
			existsInv = false
			if len(inv) > 0 {
				existsInv = true
				invLayer := <-inv
				inv <- invLayer
				invMat = invLayer.Mat.Clone("drawinvm")
			}
		}

		imgFound++
		mainLayer := <-main

		lastMat = mainLayer.Mat.Clone("mainLayerDraw")

		for key, channel := range me.layers {
			if key != "main" {
				if key != "inverseMat" {
					//log.Println("Draw2", key, len(channel))
					if len(channel) > 0 {
						item := <-channel
						channel <- item
						useImage := item.Mat.Clone("drawUseImage")

						if !item.DoNotDraw {

							if item.BasedOn == "" {
								gocv.AddWeighted(
									useImage.Get(),
									item.Weighted_Alpha,
									lastMat.Get(),
									item.Weighted_Beta,
									item.Weighted_Gamma,
									lastMat.GetPointer(),
								)
							} else if existsInv {

								//log.Println(">", item.Title, item.BasedOn)
								playground := gocv.Zeros(lastMat.GetPointer().Rows(), lastMat.GetPointer().Cols(), lastMat.GetPointer().Type())
								gocv.WarpPerspective(useImage.Get(),
									&playground,
									invMat.Get(),
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

						//item.Mat.Close()
						useImage.Close()
					}
				}
			}
		}

		if existsInv {
			invMat.Close()
		}

		mainLayer.Mat.Close()

		f := float64(lastMat.Cols()) / float64(width)
		r := float64(lastMat.Rows()) / float64(height)
		if r > f {
			f = r
		}

		log.Println(width, height, r, f)
		gocv.Resize(
			lastMat.Get(),
			lastMat.GetPointer(),
			image.Point{
				int(float64(lastMat.Cols())/f) * 2,
				int(float64(lastMat.Rows())/f) * 2,
			},
			0,
			0,
			gocv.InterpolationCubic,
		)

		return lastMat, imgFound

	}
	return lastMat, 0

}
