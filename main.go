package main

import (
	"log"
	"runtime"

	"tualo.de/deep-test/api"
	"tualo.de/deep-test/args"
	"tualo.de/deep-test/camera"
	"tualo.de/deep-test/config"
	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/globals"
	"tualo.de/deep-test/nn"
	"tualo.de/deep-test/postcheck"
	"tualo.de/deep-test/train"
	"tualo.de/deep-test/ui"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	runtime.LockOSThread()

	arguments := args.Parser()

	api := api.Static()

	if *arguments.URL != "" {
		api.RoiConfig()
		api.TitleRegions()
		api.CandidateBarcodes()
		api.BallotpaperSizes()
		//		log.Println(api.RoiConfig())
	}

	if *arguments.EnableNNServiceURL != "" {
		nn.Service().Setup(*arguments.EnableNNServiceURL, *arguments.EnableNNServiceSize, *arguments.EnableNNServiceSize)
	}

	if *arguments.Help {
		arguments.Usage()
	}

	go cv.Monitor()
	if *arguments.MemoryMonitoring {
		go cv.Display()
	}

	if *arguments.Profiling {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	if *arguments.TrainModel {
		trainer := train.Traning{}
		trainer.Run()
	}

	var appID = "io.tualo.bp"
	configData := config.Configuration()
	configData.SetAppID(appID)
	configData.Load()

	g := globals.Globals()
	g.SetDefaults()
	log.Println("globals", g)
	g.ConfigData = configData
	g.Load()

	if *arguments.DBConnection != "" {
		postcheck.Run()
	} else if *arguments.Camera >= 0 {
		camera.Run()
	} else {
		ui.StartAndRun()
	}

}
