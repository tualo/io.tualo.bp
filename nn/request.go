package nn

import (
	"encoding/base64"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gocv.io/x/gocv"
)

func (me *NeuronalNetworkService) RequestImageFloat(mat *gocv.Mat) (float64, error) {
	gocv.Resize(*mat, mat, image.Point{me.Width, me.Height}, 0, 0, gocv.InterpolationArea)
	image_bytes, _ := gocv.IMEncode(gocv.JPEGFileExt, *mat)
	image_base64 := base64.StdEncoding.EncodeToString(image_bytes.GetBytes())
	image_bytes.Close()
	return me.RequestFloat(image_base64)
}

func (me *NeuronalNetworkService) RequestFloat(data string) (float64, error) {
	client := &http.Client{}
	var req *http.Request
	var resp *http.Response
	var e1 error
	var body []byte
	req, _ = http.NewRequest(http.MethodPut, me.GetURL(), strings.NewReader(data))
	resp, e1 = client.Do(req)
	if e1 != nil {
		log.Println("Length", len(data))
		log.Panicln(e1)
		os.Exit(1)
	}

	body, _ = io.ReadAll(resp.Body)

	f, _ := strconv.ParseFloat(string(body), 64)

	return f, nil
}
