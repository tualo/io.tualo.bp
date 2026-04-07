package layer2

func (me *Layer2) Remove(title string) {

	channel, exists := me.layers[title]
	if exists {
		if len(channel) > 0 {
			oldLayer, ok := <-channel
			if ok {
				oldLayer.Mat.Close()

			}

		}
	}
}
