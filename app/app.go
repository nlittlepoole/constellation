package main

//go:generate go get -u github.com/jteeuwen/go-bindata/...
//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix assets/ assets/

import (
	"github.com/sirupsen/logrus"
	"github.com/zserge/webview"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var stop chan string
var stopped chan string
var log = logrus.New()

var Binding Observatory = Observatory{mutex: &sync.Mutex{}}

const INDEX string = "assets/index.html"

type Observatory struct {
	TimeLine Timeseries
	mutex    *sync.Mutex
}

func (o *Observatory) updateTimeLine(series Timeseries) {
	o.mutex.Lock()
	o.TimeLine = series
	o.mutex.Unlock()
}

func UpdateTimeLine(observatory *Observatory) {
	for {
		series, err := GetAllUniques(time.Hour)
		if err != nil {
			log.Error(err)
		}
		observatory.updateTimeLine(series)
		time.Sleep(time.Second * 30)
	}
}

func startObserving() {
	stop = make(chan string)
	stopped = make(chan string)
	go listen(stop, stopped)
}

func stopObserving() {
	close(stop)
	<-stopped
}

func main() {
	logFilePath := filepath.Join(getCachePath(), "observatory.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Info("Logging to ", logFilePath)
		log.Out = logFile
	} else {
		log.Info("Failed to log to file, using stderr")
	}
	defer logFile.Close()

	go UpdateTimeLine(&Binding)

	html := `data:text/html,` + url.PathEscape(string(MustAsset(INDEX)))
	w := webview.New(webview.Settings{
		Width:     1100,
		Height:    576,
		Title:     "Observatory by Constellation",
		Resizable: true,
		Debug:     true,
		URL:       html,
	})
	defer w.Exit()

	w.Dispatch(func() {
		w.Bind("observatoryBinding", &Binding)
		loadUIFramework(w)
	})
	w.Run()

}
