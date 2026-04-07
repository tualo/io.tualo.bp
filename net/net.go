package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"log"
	"math"
	"os"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/models"

	"gonum.org/v1/gonum/mat"
)

// neuralNet contains all of the information
// that defines a trained neural network.
type neuralNet struct {
	config  neuralNetConfig
	wHidden *mat.Dense
	bHidden *mat.Dense
	wOut    *mat.Dense
	bOut    *mat.Dense
}

type BB struct {
	Config          neuralNetConfig
	Content_wHidden []byte
	Content_bHidden []byte
	Content_wOut    []byte
	Content_bOut    []byte
}

type Net struct {
	loaded    bool
	neuralNet neuralNet
}

var static = Net{
	loaded: false,
}

func Static() *Net {
	if !static.loaded {
		// static.neuralNet = *Load("models/ln.json")
		static.neuralNet = *LoadStatic()

		imageSize = int(math.Sqrt(float64(static.neuralNet.config.InputNeurons)))

	}
	return &static
}

func (me *Net) MaxIndexInArray(in []float64) int {
	max := 0.0
	res := 0
	for i := 0; i < len(in); i++ {
		if max < in[i] {
			res = i
			max = in[i]
		}
	}
	return res
}

func (me *Net) Predict(img *gocv.Mat) (string, float64) {
	cvmat := img.Clone()
	input := mat.NewDense(1, imageSize*imageSize, me.readInput(&cvmat))
	cvmat.Close()
	predictions2, _ := me.neuralNet.predict(input)
	p := mat.Row(nil, 0, predictions2)
	m := me.MaxIndexInArray(p)
	return me.neuralNet.config.Labels[m], p[m] * 100
}

func (me *Net) readInput(mat *gocv.Mat) []float64 {
	result := make([]float64, imageSize*imageSize)
	gocv.Resize(*mat, mat, image.Point{imageSize, imageSize}, 0, 0, gocv.InterpolationArea)
	bytes := mat.ToBytes()
	for b := 0; b < len(bytes); b++ {
		result[b] = float64(bytes[b]) / 255
	}
	return result
}

func Save(net neuralNet, filename string) {
	// writer = {}
	b := &BB{
		Config: net.config,
	}
	var buf bytes.Buffer
	net.bHidden.MarshalBinaryTo(&buf)
	b.Content_bHidden = buf.Bytes()

	var buf2 bytes.Buffer
	net.wHidden.MarshalBinaryTo(&buf2)
	b.Content_wHidden = buf2.Bytes()

	var buf3 bytes.Buffer
	net.wOut.MarshalBinaryTo(&buf3)
	b.Content_wOut = buf3.Bytes()

	var buf4 bytes.Buffer
	net.bOut.MarshalBinaryTo(&buf4)
	b.Content_bOut = buf4.Bytes()

	json, _ := json.Marshal(b)
	// fmt.Println(string(json))

	err := os.WriteFile(filename, json, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//net.bHidden.UnmarshalBinary
}
func LoadStatic() *neuralNet {
	result := &neuralNet{}
	bb := &BB{}

	json.Unmarshal(models.Model.Content, bb)

	var bh mat.Dense
	var wh mat.Dense
	var bo mat.Dense
	var wo mat.Dense
	bh.UnmarshalBinary(bb.Content_bHidden)
	wh.UnmarshalBinary(bb.Content_wHidden)
	bo.UnmarshalBinary(bb.Content_bOut)
	wo.UnmarshalBinary(bb.Content_wOut)
	result.config = bb.Config
	result.bHidden = &bh
	result.wHidden = &wh
	result.bOut = &bo
	result.wOut = &wo
	return result
}
func Load(filename string) *neuralNet {
	result := &neuralNet{}
	bb := &BB{}
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(b, bb)

	var bh mat.Dense
	var wh mat.Dense
	var bo mat.Dense
	var wo mat.Dense
	bh.UnmarshalBinary(bb.Content_bHidden)
	wh.UnmarshalBinary(bb.Content_wHidden)
	bo.UnmarshalBinary(bb.Content_bOut)
	wo.UnmarshalBinary(bb.Content_wOut)
	result.config = bb.Config
	result.bHidden = &bh
	result.wHidden = &wh
	result.bOut = &bo
	result.wOut = &wo
	return result
}

// neuralNetConfig defines our neural network
// architecture and learning parameters.
type neuralNetConfig struct {
	Labels        []string
	InputNeurons  int
	OutputNeurons int
	HiddenNeurons int
	NumEpochs     int
	LearningRate  float64
}

var imageSize int = 32

/*
var maxFilesPerLabel int = 300

func fileList(pathName string) map[string][]string {
	result := make(map[string][]string)
	var markers []string
	var files []string
	var err error
	pathStringlen := len(pathName) + 1
	markers, err = filepath.Glob(pathName + "/*")
	if err != nil {
		return nil
	}
	for i := 0; i < len(markers); i++ {
		files, err = filepath.Glob(markers[i] + "/*")
		if err != nil {
			return nil
		}
		label := markers[i][pathStringlen:]
		pLen := len(markers[i]) + 1
		for f := 0; f < len(files) && f < maxFilesPerLabel; f++ {
			fName := files[f][pLen:]
			result[label] = append(result[label], fName)

		}
	}
	return result
}

func readImage(fileName string, labelCount int, labelIndex int) []float64 {
	result := make([]float64, imageSize*imageSize+labelCount)
	mat := gocv.IMRead(fileName, gocv.IMReadGrayScale)
	if mat.Empty() {
		fmt.Printf("Failed to read image: %s\n", fileName)
		os.Exit(1)
	}
	gocv.Resize(mat, &mat, image.Point{imageSize, imageSize}, 0, 0, gocv.InterpolationArea)
	bytes := mat.ToBytes()
	for b := 0; b < len(bytes); b++ {
		result[b] = float64(bytes[b]) / 255
	}
	offset := imageSize * imageSize
	for b := 0; b < labelCount; b++ {
		if b == labelIndex {
			result[offset+b] = 1.0
		} else {
			result[offset+b] = 0.0
		}
	}
	return result
}

func readInput(mat *gocv.Mat) []float64 {
	result := make([]float64, imageSize*imageSize)
	gocv.Resize(*mat, mat, image.Point{imageSize, imageSize}, 0, 0, gocv.InterpolationArea)
	bytes := mat.ToBytes()
	for b := 0; b < len(bytes); b++ {
		result[b] = float64(bytes[b]) / 255
	}
	return result
}
*/

/*
func makeImageInputsAndLabels(pathName string) (*mat.Dense, *mat.Dense, int, int, []string, []string) {

	files := fileList(pathName)
	list := []string{}
	labels := []string{}

	labelCount := 0
	filesCount := 0
	for _, item := range files {
		labelCount++
		filesCount += len(item)
	}
	totalImageSize := imageSize * imageSize

	inputsData := make([]float64, filesCount*imageSize*imageSize)
	labelsData := make([]float64, filesCount*labelCount)

	csvList := [][]float64{}
	labelIndex := 0
	for key, item := range files {
		log.Println(labelIndex, key)
		labels = append(labels, key)
		fc := len(item)
		for i := 0; i < fc; i++ {
			imageFile := pathName + "/" + key + "/" + item[i]
			data := readImage(imageFile, labelCount, labelIndex)
			csvList = append(csvList, data)
			list = append(list, key+"/"+item[i])
		}
		labelIndex++
	}


	labelIndex = 0
	inputIndex := 0
	for rowIndex := 0; rowIndex < len(csvList); rowIndex++ {

		for columnIndex := 0; columnIndex < len(csvList[rowIndex]); columnIndex++ {

			// log.Println(inputIndex, csvList[rowIndex][columnIndex])
			// Add to the labelsData if relevant.
			if columnIndex+1 > totalImageSize {
				labelsData[labelIndex] = csvList[rowIndex][columnIndex]
				labelIndex++
				continue
			}

			// Add the float value to the slice of floats.
			inputsData[inputIndex] = csvList[rowIndex][columnIndex]
			inputIndex++
		}
	}

	inputsD := mat.NewDense(filesCount, totalImageSize, inputsData)
	labelsD := mat.NewDense(filesCount, labelCount, labelsData)

	return inputsD, labelsD, totalImageSize, labelCount, labels, list

}
*/
/*
func mainX() {

	network := Load("my_network.txt")

	/*
		inputs, labels, totalImageSize, labelCount, outputLabels, _ := makeImageInputsAndLabels("/Users/thomashoffmann/Documents/Projects/go/deep-test/data")

		config := neuralNetConfig{
			Labels:        outputLabels,
			InputNeurons:  totalImageSize,
			OutputNeurons: labelCount,
			HiddenNeurons: 5,
			NumEpochs:     5000,
			LearningRate:  0.25,
		}

		// Train the neural network.
		network := newNetwork(config)
		if err := network.train(inputs, labels); err != nil {
			log.Fatal(err)
		}
*/

/*
		testInputs, testLabels, _, _, _, names := makeImageInputsAndLabels("/Users/thomashoffmann/Documents/Projects/go/net/test")

		// Make the predictions using the trained model.
		predictions, err := network.predict(testInputs)
		if err != nil {
			log.Fatal(err)
		}

		var truePosNeg int
		numPreds, _ := predictions.Dims()
		for i := 0; i < numPreds; i++ {

			// Get the label.
			labelRow := mat.Row(nil, i, testLabels)
			log.Println(i, labelRow, names[i])
			var prediction int
			for idx, label := range labelRow {
				if label == 1.0 {
					prediction = idx
					break
				}
			}

			// Accumulate the true positive/negative count.
			if predictions.At(i, prediction) == floats.Max(mat.Row(nil, i, predictions)) {
				truePosNeg++
			}
		}

		// Calculate the accuracy (subset accuracy).
		log.Println("numPreds", numPreds)
		log.Println("truePosNeg", truePosNeg)
		accuracy := float64(truePosNeg) / float64(numPreds)

		// Output the Accuracy value to standard out.
		fmt.Printf("\nAccuracy = %0.5f\n\n", accuracy)
		// log.Println(predictions)
	* /

	data := readImage("test/O/12063.0.0.jpg", 0, 0)
	l := len(data)

	/*

		data2 := readImage("test/O/12063.0.0.jpg", 0, 0)
		for b := 0; b < l; b++ {
			data = append(data, data2[b])
		}

		data2 = readImage("test/O/14009.0.0.jpg", 0, 0)
		for b := 0; b < l; b++ {
			data = append(data, data2[b])
		}
	* /
	start := time.Now()
	input := mat.NewDense(1, l, data)
	// tlabel := mat.NewDense(1, 2, []float64{0, 1, 1, 0, 0, 1})

	predictions2, e := network.predict(input)
	log.Println(predictions2, e)

	numPreds2, r := predictions2.Dims()

	log.Println(numPreds2, r)
	p := mat.Row(nil, 0, predictions2)
	m := MaxIndexInArray(p)
	log.Println(">>>>>>", m, network.config.Labels, p[m])

	log.Println(time.Since(start))

	// save(*network, "my_network.txt")
	/*

		for i := 0; i < numPreds2; i++ {

			labelRow := mat.Row(nil, i, tlabel)
			log.Println(i, labelRow)
			var prediction int
			for idx, label := range labelRow {
				// log.Println(i, idx, label)
				if label == 1.0 {
					prediction = idx
					break
				}
			}
			log.Println("prediction", prediction)
		}
	* /
	// log.Println(fileList("/Users/thomashoffmann/Documents/Projects/go/deep-test/data"))
	os.Exit(0)
	/*
		// Form the training matrices.
		// inputs, labels := makeInputsAndLabels("data/train.csv")
		log.Println(labels)

		// Define our network architecture and learning parameters.
		config := neuralNetConfig{
			inputNeurons:  4,
			outputNeurons: 3,
			hiddenNeurons: 3,
			numEpochs:     5000,
			learningRate:  0.3,
		}

		// Train the neural network.
		network := newNetwork(config)
		if err := network.train(inputs, labels); err != nil {
			log.Fatal(err)
		}

		// Form the testing matrices.
		testInputs, testLabels := makeInputsAndLabels("data/test.csv")

		// Make the predictions using the trained model.
		predictions, err := network.predict(testInputs)
		if err != nil {
			log.Fatal(err)
		}

		// Calculate the accuracy of our model.
		var truePosNeg int
		numPreds, _ := predictions.Dims()
		for i := 0; i < numPreds; i++ {

			// Get the label.
			labelRow := mat.Row(nil, i, testLabels)
			var prediction int
			for idx, label := range labelRow {
				if label == 1.0 {
					prediction = idx
					break
				}
			}

			// Accumulate the true positive/negative count.
			if predictions.At(i, prediction) == floats.Max(mat.Row(nil, i, predictions)) {
				truePosNeg++
			}
		}

		// Calculate the accuracy (subset accuracy).
		accuracy := float64(truePosNeg) / float64(numPreds)

		// Output the Accuracy value to standard out.
		fmt.Printf("\nAccuracy = %0.2f\n\n", accuracy)
	* /
}
*/

/*
// NewNetwork initializes a new neural network.
func newNetwork(config neuralNetConfig) *neuralNet {
	return &neuralNet{config: config}
}

// train trains a neural network using backpropagation.
func (nn *neuralNet) train(x, y *mat.Dense) error {

	// Initialize biases/weights.
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)

	wHidden := mat.NewDense(nn.config.InputNeurons, nn.config.HiddenNeurons, nil)
	bHidden := mat.NewDense(1, nn.config.HiddenNeurons, nil)
	wOut := mat.NewDense(nn.config.HiddenNeurons, nn.config.OutputNeurons, nil)
	bOut := mat.NewDense(1, nn.config.OutputNeurons, nil)

	wHiddenRaw := wHidden.RawMatrix().Data
	bHiddenRaw := bHidden.RawMatrix().Data
	wOutRaw := wOut.RawMatrix().Data
	bOutRaw := bOut.RawMatrix().Data

	for _, param := range [][]float64{
		wHiddenRaw,
		bHiddenRaw,
		wOutRaw,
		bOutRaw,
	} {
		for i := range param {
			param[i] = randGen.Float64()
		}
	}

	// Define the output of the neural network.
	output := new(mat.Dense)

	// Use backpropagation to adjust the weights and biases.
	if err := nn.backpropagate(x, y, wHidden, bHidden, wOut, bOut, output); err != nil {
		return err
	}

	// Define our trained neural network.
	nn.wHidden = wHidden
	nn.bHidden = bHidden
	nn.wOut = wOut
	nn.bOut = bOut

	return nil
}

// backpropagate completes the backpropagation method.
func (nn *neuralNet) backpropagate(x, y, wHidden, bHidden, wOut, bOut, output *mat.Dense) error {

	// Loop over the number of epochs utilizing
	// backpropagation to train our model.
	for i := 0; i < nn.config.NumEpochs; i++ {

		// Complete the feed forward process.
		hiddenLayerInput := new(mat.Dense)
		hiddenLayerInput.Mul(x, wHidden)
		addBHidden := func(_, col int, v float64) float64 { return v + bHidden.At(0, col) }
		hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

		hiddenLayerActivations := new(mat.Dense)
		applySigmoid := func(_, _ int, v float64) float64 { return sigmoid(v) }
		hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

		outputLayerInput := new(mat.Dense)
		outputLayerInput.Mul(hiddenLayerActivations, wOut)
		addBOut := func(_, col int, v float64) float64 { return v + bOut.At(0, col) }
		outputLayerInput.Apply(addBOut, outputLayerInput)
		output.Apply(applySigmoid, outputLayerInput)

		// Complete the backpropagation.
		networkError := new(mat.Dense)
		networkError.Sub(y, output)

		slopeOutputLayer := new(mat.Dense)
		applySigmoidPrime := func(_, _ int, v float64) float64 { return sigmoidPrime(v) }
		slopeOutputLayer.Apply(applySigmoidPrime, output)
		slopeHiddenLayer := new(mat.Dense)
		slopeHiddenLayer.Apply(applySigmoidPrime, hiddenLayerActivations)

		dOutput := new(mat.Dense)
		dOutput.MulElem(networkError, slopeOutputLayer)
		errorAtHiddenLayer := new(mat.Dense)
		errorAtHiddenLayer.Mul(dOutput, wOut.T())

		dHiddenLayer := new(mat.Dense)
		dHiddenLayer.MulElem(errorAtHiddenLayer, slopeHiddenLayer)

		// Adjust the parameters.
		wOutAdj := new(mat.Dense)
		wOutAdj.Mul(hiddenLayerActivations.T(), dOutput)
		wOutAdj.Scale(nn.config.LearningRate, wOutAdj)
		wOut.Add(wOut, wOutAdj)

		bOutAdj, err := sumAlongAxis(0, dOutput)
		if err != nil {
			return err
		}
		bOutAdj.Scale(nn.config.LearningRate, bOutAdj)
		bOut.Add(bOut, bOutAdj)

		wHiddenAdj := new(mat.Dense)
		wHiddenAdj.Mul(x.T(), dHiddenLayer)
		wHiddenAdj.Scale(nn.config.LearningRate, wHiddenAdj)
		wHidden.Add(wHidden, wHiddenAdj)

		bHiddenAdj, err := sumAlongAxis(0, dHiddenLayer)
		if err != nil {
			return err
		}
		bHiddenAdj.Scale(nn.config.LearningRate, bHiddenAdj)
		bHidden.Add(bHidden, bHiddenAdj)
	}

	return nil
}
*/

// sigmoid implements the sigmoid function
// for use in activation functions.
func (nn *neuralNet) sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// predict makes a prediction based on a trained
// neural network.
func (nn *neuralNet) predict(x *mat.Dense) (*mat.Dense, error) {

	// Check to make sure that our neuralNet value
	// represents a trained model.
	if nn.wHidden == nil || nn.wOut == nil {
		return nil, errors.New("the supplied weights are empty")
	}
	if nn.bHidden == nil || nn.bOut == nil {
		return nil, errors.New("the supplied biases are empty")
	}

	// Define the output of the neural network.
	output := new(mat.Dense)

	// Complete the feed forward process.
	hiddenLayerInput := new(mat.Dense)
	hiddenLayerInput.Mul(x, nn.wHidden)
	addBHidden := func(_, col int, v float64) float64 { return v + nn.bHidden.At(0, col) }
	hiddenLayerInput.Apply(addBHidden, hiddenLayerInput)

	hiddenLayerActivations := new(mat.Dense)
	applySigmoid := func(_, _ int, v float64) float64 { return nn.sigmoid(v) }
	hiddenLayerActivations.Apply(applySigmoid, hiddenLayerInput)

	outputLayerInput := new(mat.Dense)
	outputLayerInput.Mul(hiddenLayerActivations, nn.wOut)
	addBOut := func(_, col int, v float64) float64 { return v + nn.bOut.At(0, col) }
	outputLayerInput.Apply(addBOut, outputLayerInput)
	output.Apply(applySigmoid, outputLayerInput)

	return output, nil
}

/*
// sigmoid implements the sigmoid function
// for use in activation functions.
func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

// sigmoidPrime implements the derivative
// of the sigmoid function for backpropagation.
func sigmoidPrime(x float64) float64 {
	return sigmoid(x) * (1.0 - sigmoid(x))
}

// sumAlongAxis sums a matrix along a
// particular dimension, preserving the
// other dimension.
func sumAlongAxis(axis int, m *mat.Dense) (*mat.Dense, error) {

	numRows, numCols := m.Dims()

	var output *mat.Dense

	switch axis {
	case 0:
		data := make([]float64, numCols)
		for i := 0; i < numCols; i++ {
			col := mat.Col(nil, i, m)
			data[i] = floats.Sum(col)
		}
		output = mat.NewDense(1, numCols, data)
	case 1:
		data := make([]float64, numRows)
		for i := 0; i < numRows; i++ {
			row := mat.Row(nil, i, m)
			data[i] = floats.Sum(row)
		}
		output = mat.NewDense(numRows, 1, data)
	default:
		return nil, errors.New("invalid axis, must be 0 or 1")
	}

	return output, nil
}
*/
/*
func makeInputsAndLabels(fileName string) (*mat.Dense, *mat.Dense) {
	// Open the dataset file.
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file.
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 7

	// Read in all of the CSV records
	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// inputsData and labelsData will hold all the
	// float values that will eventually be
	// used to form matrices.
	inputsData := make([]float64, 4*len(rawCSVData))
	labelsData := make([]float64, 3*len(rawCSVData))

	// Will track the current index of matrix values.
	var inputsIndex int
	var labelsIndex int

	// Sequentially move the rows into a slice of floats.
	for idx, record := range rawCSVData {

		// Skip the header row.
		if idx == 0 {
			continue
		}

		// Loop over the float columns.
		for i, val := range record {

			// Convert the value to a float.
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Add to the labelsData if relevant.
			if i == 4 || i == 5 || i == 6 {
				labelsData[labelsIndex] = parsedVal
				labelsIndex++
				continue
			}

			// Add the float value to the slice of floats.
			inputsData[inputsIndex] = parsedVal
			inputsIndex++
		}
	}
	inputs := mat.NewDense(len(rawCSVData), 4, inputsData)
	labels := mat.NewDense(len(rawCSVData), 3, labelsData)
	return inputs, labels
}
*/
