package main

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"time"

	fp "path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	rp "github.com/pragmatically-dev/PoC-drawj2d-port-go/remarkablepage"
	_ "go.uber.org/automaxprocs"
)

// Config structure to hold YAML configuration
type Config struct {
	DirToSearch string `yaml:"dir_to_search"`
	FilePrefix  string `yaml:"file_prefix"`
}

var httpClient = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
		WriteBufferSize:    1 << 21,
		ReadBufferSize:     1 << 21,
	},
}

func postRmDocToWebInterface(filepath string) {
	rmData := rp.LaplacianEdgeDetection(filepath)
	rmFile := rp.GetFileNameWithoutExtension(filepath)
	rmDocBuff, rmDocPath := rp.CreateRmDoc(rmFile, rmData)
	if rmDocBuff == nil {
		fmt.Println("Error al crear zip en memoria:")
		return
	}

	go func(rmDocPath string, rmdocbuff *bytes.Buffer) {
		defer deleteFile(filepath)
		defer runtime.GC()

		url := "http://10.11.99.1"
		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		part, err := writer.CreateFormFile("file", fp.Base(rmDocPath))
		if err != nil {
			fmt.Println("Error creating the form file field:", err)
			return
		}

		part.Write(rmdocbuff.Bytes())
	

		err = writer.Close()
		if err != nil {
			fmt.Println("Error closing the writer:", err)
			return
		}

		req, err := http.NewRequest("POST", url+"/upload", &requestBody)
		if err != nil {
			fmt.Println("Error creating the request:", err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := httpClient.Do(req) // Usa el cliente HTTP reutilizable
		if err != nil {
			fmt.Println("Error sending the request:", err)
			return
		}
		defer resp.Body.Close()

	}(rmDocPath, rmDocBuff)
}

func watchForScreenshots(dirToSearch, filePrefix string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(dirToSearch)
	if err != nil {
		log.Fatal(err)
	}

	//Predicates
	doesItContainPrefix := func(event fsnotify.Event) bool {
		return strings.HasPrefix(fp.Base(event.Name), filePrefix)
	}
	isNewFile := func(event fsnotify.Event) bool {
		return event.Has(fsnotify.Create)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Chmod) || event.Has(fsnotify.Create) {
				fmt.Println(event.Name)
			}
			if isNewFile(event) && doesItContainPrefix(event) {
				//DO NOT REMOVE: The remarkable takes 1200ms to save the png to /home/root
				//Wasted hours for trying to fix this: 15
				time.Sleep(1100 * time.Millisecond)
				rp.DebugPrint("Screenshot found: " + event.Name)
				go postRmDocToWebInterface(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func deleteFile(filepath string) error {
	return os.Remove(filepath)
}

func AppStart() {

	config := &Config{
		DirToSearch: "/home/root",
		FilePrefix:  "Screenshot",
	}

	fmt.Println("<--- Looking for new Screenshots --->")
	//go watchForScreenshots("/home/root/.local/share/remarkable/xochitl", config.FilePrefix)
	watchForScreenshots(config.DirToSearch, config.FilePrefix)

}

func main() {
	AppStart()

}
