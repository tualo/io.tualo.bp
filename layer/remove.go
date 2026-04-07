package layer

import (
	"time"

	"tualo.de/deep-test/structs"
)

func (me *Layer) Remove(title string) {
	me.removeChannel <- structs.LayerRemoveStruct{
		Title:     title,
		Timestamp: time.Now(),
	}
}
