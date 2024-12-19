package main

// https://github.com/tensorflow/build/tree/master/golang_install_guide
// https://github.com/galeone/tfgo/
import (
	"fmt"
	"image"

	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
	"gocv.io/x/gocv"
)

func adjustInput(inData gocv.Mat) [][][][]float32 {
	width := inData.Cols()
	height := inData.Rows()
	channels := inData.Channels()
	outData := make([][][][]float32, 1)
	outData[0] = make([][][]float32, width)

	// Scale the data and expand dimensions
	for w := 0; w < width; w++ {
		outData[0][w] = make([][]float32, height)
		for h := 0; h < height; h++ {
			outData[0][w][h] = make([]float32, channels)
			val := inData.GetVecbAt(h, w)
			for c := 0; c < channels; c++ {
				outData[0][w][h][c] = (float32(val[c]) - 127.5) * 0.0078125
			}
		}
	}
	return outData
}

func inputTensor(size int) (*tf.Tensor, error) {

	imageData := [][]float32{make([]float32, size)}
	return tf.NewTensor(imageData)
}

func main() {
	model := tg.LoadModel("/Users/thomashoffmann/Documents/Projects/go/create-dataset/", []string{"serve"}, nil)

	numBox := 1
	inputBuf := make([][][][]float32, numBox)
	tmp := gocv.IMRead("/Users/thomashoffmann/Documents/Projects/go/io.tualo.bp/splitroi/test/0.3.jpg", gocv.IMReadGrayScale)
	for i := 0; i < numBox; i++ {
		resized := gocv.NewMat()
		defer resized.Close()
		gocv.Resize(tmp, &resized, image.Pt(24, 24), 0, 0, gocv.InterpolationLinear)
		inputBuf[i] = adjustInput(resized)[0]
	}

	inputBufTensor, _ := tf.NewTensor(inputBuf)

	/*
		root := tg.NewRoot()
		grayImg := image.Read(root, "/Users/thomashoffmann/Documents/Projects/go/io.tualo.bp/splitroi/test/0.3.jpg", 1)
		grayImg = grayImg.Scale(0, 255)
	*/

	/*
		fakeInput, _ := tf.NewTensor("/Users/thomashoffmann/Documents/Projects/go/io.tualo.bp/splitroi/test/0.3.jpg")
		log.Println(fakeInput.Value())
	*/
	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("sequential_1", 0): inputBufTensor,
	})

	predictions := results[0]
	fmt.Println(predictions.Value())
}
