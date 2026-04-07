package shadow

import (
	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

type Shadow struct {
}

func (me *Shadow) Remove(img cv.Mat) cv.Mat {
	channels := gocv.Split(img.Get())

	// log.Println("T1", img.Type().String())
	var mv []gocv.Mat
	for i := 0; i < len(channels); i++ {
		dilated_img := gocv.NewMat()
		// log.Println("C", i, channels[i].Type().String())
		bg_img := gocv.NewMatFromScalar(gocv.NewScalar(255.0, 255.0, 255.0, 0.0), channels[i].Type())
		diff_img := gocv.NewMat()
		norm_img := gocv.NewMat()
		o := gocv.Ones(77, 77, channels[i].Type())
		gocv.Dilate(channels[i], &dilated_img, o)
		gocv.MedianBlur(dilated_img, &bg_img, 91)
		gocv.AbsDiff(channels[i], bg_img, &diff_img)
		gocv.Normalize(diff_img, &norm_img, 0.0, 255.0, gocv.NormMinMax)
		mv = append(mv, norm_img)

		dilated_img.Close()
		diff_img.Close()
		bg_img.Close()
		o.Close()

	}
	res := gocv.NewMat()
	gocv.Merge(mv, &res)
	// log.Println("R", res.Type().String())

	// gocv.Invert(res, &res, gocv.SolveDecompositionLu)
	return cv.NewFromMat("shadow", res)
}
