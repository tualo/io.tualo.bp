package train

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
	"gocv.io/x/gocv"
	"tualo.de/deep-test/args"
)

type Traning struct {
	outputFile string
}

func (me *Traning) check(e error) {
	if e != nil {
		panic(e)
	}
}

func (me *Traning) MatFromFile(fileName string) gocv.Mat {
	mat := gocv.IMRead(fileName, gocv.IMReadGrayScale)
	if mat.Empty() {
		fmt.Printf("Failed to read image: %s\n", fileName)
		os.Exit(1)
	}
	return mat
}

func (me *Traning) MatToExample(inputImage *gocv.Mat, result float64) training.Example {
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

func (me *Traning) loadSamples() (training.Examples, error) {
	var markers []string
	var files []string
	var samples training.Examples
	var err error
	path := *args.Parser().TrainModelImagePath
	markers, err = filepath.Glob(path + "/*")

	log.Println(markers, err)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(markers); i++ {
		files, err = filepath.Glob(markers[i] + "/*")
		if err != nil {
			return nil, err
		}
		for f := 0; f < len(files) && f < 50; f++ {
			log.Println(i, files[f])
			// fmt.Printf("Index %d Name %s \n", f, files[f])
			mat := gocv.IMRead(files[f], gocv.IMReadGrayScale)
			if mat.Empty() {
				fmt.Printf("Failed to read image: %s\n", files[f])
				os.Exit(1)
			}
			gocv.Resize(mat, &mat, image.Point{32, 32}, 0, 0, gocv.InterpolationArea)
			bytes := mat.ToBytes()
			//log.Println(bytes)

			var data []float64
			var response []float64
			for b := 0; b < len(bytes); b++ {
				data = append(data, float64(bytes[b])/255)
			}
			// log.Println(data)
			response = append(response, float64(i)*1)
			var sample training.Example
			sample.Input = data
			sample.Response = response
			samples = append(samples, sample)
			mat.Close()
		}
	}
	return samples, nil
}

func (me *Traning) Run() {
	me.outputFile = "models/localnn.json"
	data, e1 := me.loadSamples()
	me.check(e1)

	n := deep.NewNeural(&deep.Config{
		Inputs:     len(data[0].Input),
		Layout:     []int{50, 30, 2},
		Activation: deep.ActivationSigmoid,
		Mode:       deep.ModeMultiLabel,
		Weight:     deep.NewUniform(0.6, 0.1),
		Bias:       true,
	})

	//trainer := training.NewTrainer(training.NewSGD(0.005, 0.5, 1e-6, true), 50)
	//trainer := training.NewBatchTrainer(training.NewSGD(0.005, 0.1, 0, true), 50, 300, 16)
	//trainer := training.NewTrainer(training.NewAdam(0.1, 0, 0, 0), 50)
	trainer := training.NewBatchTrainer(training.NewSGD(0.05, 0.1, 1e-6, true), 50, len(data)/2, 12)
	data, heldout := data.Split(0.5)
	trainer.Train(n, data, heldout, 1000)

	/*
		n := deep.NewNeural(&deep.Config{
			Inputs: len(data[0].Input),
			Layout: []int{50, 10, 1},
			// Activation: deep.ActivationReLU,
			Activation: deep.ActivationSigmoid,

			Mode:   deep.ModeMultiLabel,
			Weight: deep.NewNormal(0.6, 0.1), // slight positive bias helps ReLU
			Bias:   true,
		})
		// params: learning rate, momentum, alpha decay, nesterov
		//optimizer := training.NewSGD(0.05, 0.1, 1e-6, true)
		optimizer := training.NewAdam(0.00001, 0.0009, 0.999, 0.0001)
		//optimizer := training.NewAdam(0.001, 0.9, 0.999, 1e-8)

		// params: optimizer, verbosity (print stats at every 50th iteration)
		trainer := training.NewTrainer(optimizer, 5)
		// trainer := training.NewBatchTrainer(optimizer, 1, 200, 4)

		training, heldout := data.Split(0.75)
		trainer.Train(n, training, heldout, 150) // training, validation, iterations
	*/
	fmt.Println("first", "=>", n.Predict(data[0].Input), data[0].Response)
	fmt.Println("last", "=>", n.Predict(data[len(data)-2].Input), data[len(data)-2].Response)

	bytes, dumpError := n.Marshal()
	me.check(dumpError)

	err := os.WriteFile(me.outputFile, bytes, 0644)
	me.check(err)
}
