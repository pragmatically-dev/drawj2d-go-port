package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"time"

	fp "path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	rp "github.com/pragmatically-dev/drawj2d-rm/remarkablepage"
	_ "go.uber.org/automaxprocs"
	"gopkg.in/yaml.v2"
)

// Config structure to hold YAML configuration
type Config struct {
	DirToSearch   string `yaml:"dir_to_search"`
	DirToSave     string `yaml:"dir_to_save"`
	FilePrefix    string `yaml:"file_prefix"`
	ServerAddress string `yaml:"server_address"`
}

var httpClient = &http.Client{}

func postRmDocToWebInterface(filepath, dirToSave string) {
	rmData := rp.LaplacianEdgeDetection(filepath, dirToSave)
	rmFile := rp.GetFileNameWithoutExtension(filepath)
	rmFile = fmt.Sprintf("%s/%s.rm", dirToSave, rmFile)
	rmDocBuff, rmDocPath := rp.CreateRmDoc(rmFile, dirToSave, rmData)
	if rmDocBuff == nil {
		fmt.Println("Error al crear zip en memoria:")
		return
	}

	go func(rmDocPath string, rmdocbuff *bytes.Buffer) {
		defer deleteFile(rmDocPath)
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

		_, err = io.Copy(part, rmdocbuff)
		if err != nil {
			fmt.Println("Error copying the file content:", err)
			return
		}

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

func watchForScreenshots(dirToSearch, filePrefix, dirToSave string) {
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

			if isNewFile(event) && doesItContainPrefix(event) {
				//DO NOT REMOVE: The remarkable takes 1200ms to save the png to /home/root
				//Wasted hours for trying to fix this: 15
				time.Sleep(1200 * time.Millisecond)
				rp.DebugPrint("Screenshot found: " + event.Name)
				go postRmDocToWebInterface(event.Name, dirToSave)
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
	configFile, err := os.ReadFile("/home/root/.config/png2rm/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Println("<--- Looking for new Screenshots --->")
	watchForScreenshots(config.DirToSearch, config.FilePrefix, config.DirToSave)
}

func main() {
	AppStart()

}
