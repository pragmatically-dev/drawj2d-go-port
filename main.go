package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	fp "path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	rp "github.com/pragmatically-dev/drawj2d-rm/remarkablepage"
	"gopkg.in/yaml.v2"
)

// Config structure to hold YAML configuration
type Config struct {
	DirToSearch   string `yaml:"dir_to_search"`
	DirToSave     string `yaml:"dir_to_save"`
	FilePrefix    string `yaml:"file_prefix"`
	ServerAddress string `yaml:"server_address"`
}

func postRmDocToWebInterface(filepath, DirToSave string) {
	//fmt.Println("Starting Conversion")
	rp.CannyEdgeDetection(filepath, DirToSave)

	//fmt.Println("File testPNGConversion.rm generated successfully.")
	rmFile := rp.GetFileNameWithoutExtension(filepath)
	rmFile = fmt.Sprintf("%s/%s.rm", DirToSave, rmFile)
	rmDocPath := rp.CreateRmDoc(rmFile, DirToSave)

	go func() {
		url := "http://10.11.99.1"

		file, err := os.Open(rmDocPath)
		if err != nil {
			fmt.Println("Error al abrir el archivo:", err)
			return
		}
		defer file.Close()

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)

		// Create a form file field
		part, err := writer.CreateFormFile("file", fp.Base(file.Name()))
		if err != nil {
			fmt.Println("Error creating the form file field:", err)
			return
		}

		// Copy the file content to the form field
		_, err = io.Copy(part, file)
		if err != nil {
			fmt.Println("Error copying the file content:", err)
			return
		}

		// Close the writer to complete the multipart form
		err = writer.Close()
		if err != nil {
			fmt.Println("Error closing the writer:", err)
			return
		}

		// Create the HTTP request
		req, err := http.NewRequest("POST", url+"/upload", &requestBody)
		if err != nil {
			fmt.Println("Error creating the request:", err)
			return
		}

		// Set the content type
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending the request:", err)
			return
		}
		defer resp.Body.Close()
		defer client.CloseIdleConnections()

		go deleteFile(rmDocPath)
	}()

}

func watchForScreenshots(dirToSearch, filePrefix, dirToSave string) {
	//The watcher will give us two channels, one for Events and other for errors
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	err = watcher.Add(dirToSearch)
	if err != nil {
		log.Fatal(err)
	}

	go func() {

		for {
			// Here we receive data from two channels
			//A select blocks until one of its cases can run, then it executes that case.
			//more info on [https://blog.stackademic.com/go-concurrency-visually-explained-select-statement-b546596c8e6b]
			select {

			case event, ok := <-watcher.Events:

				if !ok {

					return
				}

				if (event.Has(fsnotify.Create)) && (strings.HasPrefix(fp.Base(event.Name), filePrefix)) {

					//DON'T REMOVE: Remarkable screenshot takes almost a sec to save on the filesystm
					rp.DebugPrint("Screenshot found: " + event.Name)
					time.Sleep(1200 * time.Millisecond)

					go postRmDocToWebInterface(event.Name, dirToSave)
					time.Sleep(5 * time.Second)
					if err := deleteFile(event.Name); err != nil {
						log.Printf("Error deleting file: %v\n", err)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	//SYNC GO ROUTINE
	<-done
}

func deleteFile(filepath string) error {
	return os.Remove(filepath)
}

func main() {
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
