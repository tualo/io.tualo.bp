package train

import (
	"image"
	"log"
	"os"

	deep "github.com/patrikeh/go-deep"
	"gocv.io/x/gocv"
	"tualo.de/deep-test/args"
)

type Predict struct {
	neural *deep.Neural
}

var predict *Predict

func (me *Predict) check(e error) {
	if e != nil {
		panic(e)
	}
}

func (me *Predict) Predict(inputImage *gocv.Mat) []float64 {
	var data []float64

	mat := inputImage.Clone()
	gocv.Resize(
		mat,
		&mat,
		image.Point{64, 64},
		0,
		0,
		gocv.InterpolationArea,
	)
	bytes := mat.ToBytes()
	for b := 0; b < len(bytes); b++ {
		data = append(data, float64(bytes[b])/255)
	}
	log.Println("len", len(data))
	log.Println("data", data)
	mat.Close()
	return predict.neural.Predict(data)
}

func Static() *Predict {
	var bytes []byte
	var error error
	predict = &Predict{}

	if *args.Parser().EnableLocalNN {
		bytes, error = os.ReadFile("models/dump.json")
		predict.check(error)

		predict.neural, error = deep.Unmarshal(bytes)
		predict.check(error)
	}
	return predict
}
