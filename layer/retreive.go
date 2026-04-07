package layer

import (
	"log"
	"sync"
)

var mutex sync.Mutex

func (me *Layer) retreive() {
	isRunning := false
	for me.enabled {

		if !isRunning && len(me.removeChannel) > 0 {
			item := <-me.removeChannel
			log.Println("retreive", "removeChannel", item.Title)

			_, prs := me.drawLayers[item.Title]
			if prs {
				me.drawLayers[item.Title].DoNotDraw = true
			}
		}
		if !isRunning && len(me.channel) > 0 {
			isRunning = true

			mutex.Lock()
			item := <-me.channel
			element, prs := me.drawLayers[item.Title]

			if prs {
				//if false {

				log.Println("element.Mat.Close", element.Mat.Name())

				element.Mat.Close()

				//}
				me.drawLayers[item.Title] = &item
			} else {
				me.drawLayers[item.Title] = &item
			}
			mutex.Unlock()
			isRunning = false

		}
	}
}
