package cv

import (
	"fmt"
	"time"

	"gocv.io/x/gocv"
)

type MonitorData struct {
	Count int
	Start time.Time
}

type MonitorItem struct {
	Name  string
	Add   int
	Start time.Time
}

var mat_monitor map[string]*MonitorData
var monitorChannel chan MonitorItem = make(chan MonitorItem, 1000)

func Display() {
	for {
		time.Sleep(10 * time.Second)

		fmt.Println("================used mats================")
		s := fmt.Sprintf("%-30s", "Key")
		v := fmt.Sprintf("%*s", 10, "Count")
		fmt.Println(s, v)
		for key, item := range mat_monitor {
			s = fmt.Sprintf("%-30s", key)
			v = fmt.Sprintf("%*d", 10, item.Count)
			fmt.Println(s, v)
		}
		fmt.Println("=========================================")
	}
}
func Monitor() {
	mat_monitor = make(map[string]*MonitorData)
	for {

		item, ok := <-monitorChannel
		if ok {
			mon, prs := mat_monitor[item.Name]
			if prs {
				mon.Count += item.Add
				if mon.Count < 0 {
					mon.Count = 0
				}

			} else {
				mat_monitor[item.Name] = &MonitorData{
					Count: 0,
					Start: item.Start,
				}
			}
		}
	}
}

var MatCounter int = 0

type Mat struct {
	mat         gocv.Mat
	initialized bool
	id          int
	t           time.Time
	source      string
	inUse       chan bool
}

func (me *Mat) FromMat(source string, m *gocv.Mat) {
	me.mat = *m
	me.init(source)
}

func (me *Mat) Initialized() bool {
	return me.initialized
}

func (me *Mat) Name() string {
	return me.source
}

func (me *Mat) init(source string) {
	me.initialized = true
	me.id = MatCounter
	me.source = source

	me.t = time.Now()
	MatCounter++
	if len(monitorChannel) < cap(monitorChannel) {
		monitorChannel <- MonitorItem{
			Name:  source,
			Start: me.t,
			Add:   1,
		}
	}
	/*
		go func(m *Mat) {
			time.Sleep(50 * time.Second)
			if m.initialized {

				log.Println("i'm still alive", m.id, "from", m.source, "total", MatCounter, m.mat.Closed(), "since me.t", me.t)
				if m.source != "origianlImg" && m.source != "lastPaperMat" && !m.mat.Closed() {
					m.Close()
				}

			}
		}(me)
	*/
}

func NewMat(source string) Mat {
	me := Mat{}
	me.mat = gocv.NewMat()
	me.init(source)
	return me
}

func NewFromMat(source string, m gocv.Mat) Mat {
	me := Mat{
		inUse: make(chan bool, 1),
	}
	me.inUse <- true
	me.mat = m.Clone()
	me.init(source)
	m.Close()
	<-me.inUse
	return me
}

func (me *Mat) Close() {
	// log.Println("closing ", me.id, " from ", me.source)
	me.initialized = false
	me.mat.Close()
	//if me.initialized {
	MatCounter--

	monitorChannel <- MonitorItem{
		Name:  me.source,
		Start: me.t,
		Add:   -1,
	}
	//}
}

func (me *Mat) Clone(source string) Mat {
	return NewFromMat(source, me.mat.Clone())
}

func (me *Mat) Cols() int {
	return me.mat.Cols()
}

func (me *Mat) Rows() int {
	return me.mat.Rows()
}

func (me *Mat) Type() gocv.MatType {
	return me.mat.Type()
}

func (me *Mat) Get() gocv.Mat {
	return me.mat
}
func (me *Mat) GetPointer() *gocv.Mat {
	return &me.mat
}
