package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/getlantern/systray"
	"github.com/kuangcp/logger"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//go:embed sync.png
var iconImg string

var (
	sideList []string // 对端列表 格式 host:port
	port     int
	version  bool
	initSide string
)

var (
	lastMod = time.Now()
)

func init() {
	flag.IntVar(&port, "p", 8000, "port")
	flag.BoolVar(&version, "v", false, "version")
	flag.StringVar(&initSide, "s", "", "init side url")
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("1.0.0")
		return
	}

	initSideBind()

	go webServer()

	systray.Run(OnReady, OnExit)
}

func syncFile() []string {
	var result []string
	dir, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return result
	}

	for _, info := range dir {
		if info.IsDir() {
			//fmt.Println("dir", info)
			continue
		}
		if lastMod.Before(info.ModTime()) {
			logger.Info("need sync", info.Name(), info.ModTime())
			result = append(result, info.Name())
		}
	}
	if len(result) != 0 {
		lastMod = time.Now()
	}
	return result
}

func webServer() {
	http.HandleFunc("/sync", func(writer http.ResponseWriter, request *http.Request) {
		name := request.URL.Query().Get("name")
		unescape, err := url.QueryUnescape(name)
		if err != nil {
			logger.Error(err)
			return
		}

		open, err := os.Create(unescape)
		if err != nil {
			logger.Error(err)
			return
		}

		var buf = make([]byte, 4096)
		for {
			read, err := request.Body.Read(buf)
			if read != 0 {
				open.Write(buf[:read])
			}
			if err != nil {
				break
			}
		}

		open.Close()
	})
	http.HandleFunc("/register", func(writer http.ResponseWriter, request *http.Request) {
		client := request.Header.Get("self")
		logger.Info("add new", client)
		sideList = append(sideList, client)
		writer.Write([]byte("OK"))
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func initSideBind() {
	if initSide == "" {
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://"+initSide+"/register", nil)
	req.Header.Set("self", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
	sideList = append(sideList, initSide)
}