package train

import (
	"fmt"
	"os"

	"gocv.io/x/gocv"
)

func MatFromFile(fileName string) gocv.Mat {
	mat := gocv.IMRead(fileName, gocv.IMReadGrayScale)
	if mat.Empty() {
		fmt.Printf("Failed to read image: %s\n", fileName)
		os.Exit(1)
	}
	return mat
}
