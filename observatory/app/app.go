package main

//go:generate go get -u github.com/jteeuwen/go-bindata/...
//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix assets/ assets/

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"time"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"

	"github.com/zserge/webview"
)

var stop chan string
var stopped chan string

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			if path == "" {
				path = "index.html"
			}
			if bs, err := Asset(path); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				io.Copy(w, bytes.NewBuffer(bs))
			}
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func startObserving(){
     stop = make(chan string)
     stopped = make(chan string)
     go listen(stop, stopped)
}

func stopObserving(){
     close(stop)
     <-stopped     
}

// empty string response means everything is cool
type Response struct {
        Cmd string `json:"cmd"`
	Error string `json:"error"`
	Data interface{}   `json:"data"`
}

func render(w webview.WebView, resp Response) {
	b, err := json.Marshal(resp)
	if err == nil {
		w.Eval(fmt.Sprintf("rpc.render(%s)", string(b)))
	}
}

func handleRPC(w webview.WebView, data string) {
	cmd := struct {
		Name string `json:"cmd"`
	}{}
	if err := json.Unmarshal([]byte(data), &cmd); err != nil {
		log.Println(err)
		return
	}
	switch cmd.Name {
	case "start_observing":
		startObserving()
		render(w, Response{cmd.Name, "", "Success"})
	case "stop_observing":
		stopObserving()
		render(w, Response{cmd.Name, "", "Success"})
	case "active_uniques":
	        uniques, err := GetCurrentUniques(ACTIVE_SETTINGS.Window())
		var errString string
		if err != nil {
		   errString = err.Error()
		}
		render(w, Response{cmd.Name, errString, uniques})
        }
}

func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Width:  1100,
		Height: 576,
		Title:  "Constellation",
		Resizable: true,
		Debug: true,
		URL:    url,
		ExternalInvokeCallback: handleRPC,
	})
	defer w.Exit()
	w.Run()	


     //fmt.Println(GetAllUniques(time.Minute))

     //fmt.Println(GetReturningUniques(time.Now().Add(-300 * time.Second), time.Now()))
     //fmt.Println(GetStrengthHistogram(time.Now().Add(-300 * time.Second), time.Now()))
}
