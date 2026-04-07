package train

import (
	"image"

	"github.com/patrikeh/go-deep/training"
	"gocv.io/x/gocv"
)

func MatToExample(inputImage *gocv.Mat, result float64) training.Example {
	mat := inputImage.Clone()
	gocv.Resize(mat, &mat, image.Point{64, 64}, 0, 0, gocv.InterpolationArea)
	bytes := mat.ToBytes()
	var data []float64
	var response []float64
	for b := 0; b < len(bytes); b++ {
		data = append(data, float64(bytes[b])/255)
	}
	// log.Println(data)
	response = append(response, result)
	var sample training.Example
	sample.Input = data
	sample.Response = response
	return sample
}
