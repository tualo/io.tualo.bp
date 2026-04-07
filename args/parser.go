package args

import (
	"flag"
)

type ParserObject struct {
	initialized bool

	Help             *bool
	Profiling        *bool
	DBConnection     *string
	DBFilter         *string
	EnableCVWindow   *bool
	MemoryMonitoring *bool

	SaveBestMarkerImage     *bool
	SaveBestMarkerImagePath *string

	EnableNNServiceQuery *bool
	EnableNNServiceURL   *string
	EnableNNServiceSize  *int

	EnableLocalNN     *bool
	EnableLocalNNFile *string

	Camera         *int
	InputImageFile *string

	URL      *string
	Username *string
	Password *string

	TrainModel          *bool
	TrainModelImagePath *string

	BlurryThreshold *int
}

var parser = &ParserObject{
	initialized: false,
}

func (me *ParserObject) Usage() {
	flag.Usage()
}

func Parser() *ParserObject {
	if !parser.initialized {

		// flag.Usage = parser.Usage
		parser.Init()
		parser.initialized = true
	}
	return parser
}
func (me *ParserObject) Init() {
	me.Help = flag.Bool("h", false, "show usage")

	me.Profiling = flag.Bool("p", false, "enable profiling on port 6060")
	me.MemoryMonitoring = flag.Bool("m", false, "enable memory monitoring")

	me.EnableCVWindow = flag.Bool("cv", false, "enable opencv window")
	me.SaveBestMarkerImage = flag.Bool("savebestmarkers", false, "save best found markers")
	me.SaveBestMarkerImagePath = flag.String("savebestmarkerspath", "data", "save best found markers")

	me.EnableNNServiceQuery = flag.Bool("nnservice", false, "enable neuronal network service queries")
	me.EnableNNServiceURL = flag.String("nnserviceurl", "", "neuronal network service url")
	me.EnableNNServiceSize = flag.Int("nnservicesize", 32, "neuronal network service input size")

	me.DBConnection = flag.String("dbconnection", "", "username:@tcp(127.0.0.1:3306)/database")
	me.DBFilter = flag.String("filtersql", "pagination_id like '12345'", "filter sql string")

	me.Camera = flag.Int("camera", -1, "camera number")

	me.InputImageFile = flag.String("input_image", "", "input image file")

	me.URL = flag.String("url", "", "backend url")
	me.Username = flag.String("username", "", "backend username")
	me.Password = flag.String("password", "", "backend password")

	me.TrainModel = flag.Bool("t", false, "train a model")
	me.TrainModelImagePath = flag.String("train_data_dir", "data", "train data dir")

	me.EnableLocalNN = flag.Bool("nn", false, "use the local nn model")
	me.EnableLocalNNFile = flag.String("nnfile", "", "use the local nn model")

	me.BlurryThreshold = flag.Int("blurryThreshold", 255, "images below will not be used, -1 to disable the check")

	flag.Parse()
}
