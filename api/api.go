package api

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"tualo.de/deep-test/args"
	"tualo.de/deep-test/structs"
)

type API struct {
	// cookies []http.Cookie
	jar     *cookiejar.Jar
	timeout time.Duration
	// Timeout = time.Duration(10 * time.Second)

	systemURL string
	lastPing  PingResponse

	allRois            map[int]*structs.ROIS
	AllTitles          map[int]string
	AllTitlesInv       map[string]int
	AllOCRTitles       map[int]string
	AllOCRTitlesInv    map[string]int
	titleRegionsConfig []TitleRegionsConfig

	RegionsOfInterestList []structs.ROI
	CandidateBarcodeList  []CandidateBarcode
	BallotpaperSizesList  []BallotpaperSizes
}

var api *API

func Static() *API {
	if api == nil {

		arguments := args.Parser()
		api = &API{
			timeout:   time.Duration(10 * time.Second),
			systemURL: *arguments.URL,
		}
		api.initJar()

		api.AllTitles = make(map[int]string)
		api.AllTitlesInv = make(map[string]int)
		api.AllOCRTitles = make(map[int]string)
		api.AllOCRTitlesInv = make(map[string]int)

		api.allRois = make(map[int]*structs.ROIS)

		if *arguments.Username != "" && *arguments.Password != "" {
			res, err := api.Login(*arguments.Username, *arguments.Password)
			if err != nil {
				os.Exit(1)
			} else {
				api.Switch("bwbriefwahl_muenchen")
				if res.Success {
					api.RoiConfig()
					api.TitleRegions()
					api.CandidateBarcodes()
				} else {
					fmt.Println(res.Msg)
					os.Exit(1)
				}
			}

		}
	}
	return api
}

func (me *API) SetUrl(url string) {
	me.systemURL = url
}

func (me *API) dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, me.timeout)
}

func (me *API) initJar() {
	if me.jar == nil {
		jar, err := cookiejar.New(&cookiejar.Options{})
		if err != nil {
			log.Fatal(err)
		}
		me.jar = jar
	}
}

func (me *API) Get(url string) (string, error) {
	me.initJar()

	transport := http.Transport{
		Dial: me.dialTimeout,
	}
	client := http.Client{
		Transport: &transport,
		Jar:       me.jar,
	}
	var resp *http.Response
	var err error
	var body []byte
	resp, err = client.Get(url)
	if err != nil {
		return "", err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (me *API) PostForm(url string, data url.Values) (string, error) {
	me.initJar()
	transport := http.Transport{
		Dial: me.dialTimeout,
	}
	client := http.Client{
		Transport: &transport,
		Jar:       me.jar,
	}
	var resp *http.Response
	var err error
	var body []byte
	resp, err = client.PostForm(url, data)
	if err != nil {
		log.Println("PostForm ERROR", err)
		return "", err
	}

	body, err = io.ReadAll(resp.Body)
	// fmt.Println("POSTFORM RESULT", string(body))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (me *API) Post(url string, data string) (string, error) {
	me.initJar()
	// fmt.Println("POST", url, data)
	transport := http.Transport{
		Dial: me.dialTimeout,
	}
	client := http.Client{
		Transport: &transport,
		Jar:       me.jar,
	}
	var resp *http.Response
	var err error
	var body []byte
	resp, err = client.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		log.Println("Post ERROR", err)
		return "", err
	}

	body, err = io.ReadAll(resp.Body)
	// fmt.Println("POST RESULT", string(body))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (me *API) GetTitleRegionsConfig() []TitleRegionsConfig {
	return me.titleRegionsConfig
}

func (me *API) GetBallotpaper(id int) *BallotpaperSizes {

	for i := 0; i < len(me.BallotpaperSizesList); i++ {
		if me.BallotpaperSizesList[i].Ballotpaper == id {
			return &me.BallotpaperSizesList[i]
		}
	}
	return nil
}
