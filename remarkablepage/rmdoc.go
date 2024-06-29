package remarkablepage

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ReMarkableAPIrmdoc representa la estructura para empaquetar archivos .rm en un .rmdoc
type ReMarkableAPIrmdoc struct {
	Content          string `json:"content"`
	NotebookMetadata string `json:"notebookmetadata"`
	Metadata0rm      string `json:"metadata0rm"`
	Rmdata           []byte `json:"-"`
	Time             int64  `json:"time"`
	dirToSave        string
	internalBuffer   *bytes.Buffer
}

// NewReMarkableAPIrmdoc crea una nueva instancia de ReMarkableAPIrmdoc
func NewReMarkableAPIrmdoc(zipfile, dirToSave string, rmdata []byte) *ReMarkableAPIrmdoc {
	rmdoc := &ReMarkableAPIrmdoc{
		Rmdata:    rmdata,
		Time:      time.Now().Unix(),
		dirToSave: dirToSave,
	}
	rmdoc.process(zipfile)
	return rmdoc
}

func (rmdoc *ReMarkableAPIrmdoc) process(zipfile string) {
	notebookID := uuid.NewString()
	pageID := uuid.NewString()
	visibleName := filepath.Base(zipfile)
	if strings.HasSuffix(visibleName, ".rmdoc") {
		visibleName = visibleName[:len(visibleName)-len(".rmdoc")]
	}
	if strings.HasPrefix(visibleName, "out-") && len(visibleName) > 4 {
		visibleName = visibleName[len("out-"):]
	}

	rmdoc.Content = rmdoc.createContent(pageID)
	rmdoc.NotebookMetadata = rmdoc.createNotebookMetadata(visibleName)

	rmdoc.writeZip(zipfile, notebookID, pageID)
}

func (rmdoc *ReMarkableAPIrmdoc) writeZip(zipfile, notebookID, pageID string) {
	f := new(bytes.Buffer)

	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()

	contentData := []byte(rmdoc.Content)
	notebookMetadataData := []byte(rmdoc.NotebookMetadata)

	// Escribir contenido
	contentFile, err := zipWriter.Create(notebookID + ".content")
	if err != nil {
		log.Fatalf("Error creating zip entry for content: %v", err)
	}
	_, err = contentFile.Write(contentData)
	if err != nil {
		log.Fatalf("Error writing content to zip entry: %v", err)
	}

	// Escribir metadatos de cuaderno
	metadataFile, err := zipWriter.Create(notebookID + ".metadata")
	if err != nil {
		log.Fatalf("Error creating zip entry for metadata: %v", err)
	}
	_, err = metadataFile.Write(notebookMetadataData)
	if err != nil {
		log.Fatalf("Error writing metadata to zip entry: %v", err)
	}

	// Escribir archivo .rm
	rmFile, err := zipWriter.Create(notebookID + "/" + pageID + ".rm")
	if err != nil {
		log.Fatalf("Error creating zip entry for .rm file: %v", err)
	}
	_, err = rmFile.Write(rmdoc.Rmdata)
	if err != nil {
		log.Fatalf("Error writing .rm file to zip entry: %v", err)
	}

	rmdoc.internalBuffer = f
}

func (rmdoc *ReMarkableAPIrmdoc) createContent(pageID string) string {
	// Crear contenido JSON
	content := struct {
		CPages struct {
			LastOpened struct {
				Timestamp string `json:"timestamp"`
				Value     string `json:"value"`
			} `json:"lastOpened"`
			Original struct {
				Timestamp string `json:"timestamp"`
				Value     int    `json:"value"`
			} `json:"original"`
			Pages []struct {
				ID  string `json:"id"`
				Idx struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"idx"`
				Template struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"template"`
			} `json:"pages"`
			UUIDs []struct {
				First  string `json:"first"`
				Second int    `json:"second"`
			} `json:"uuids"`
		} `json:"cPages"`
		CoverPageNumber       int    `json:"coverPageNumber"`
		CustomZoomCenterX     int    `json:"customZoomCenterX"`
		CustomZoomCenterY     int    `json:"customZoomCenterY"`
		CustomZoomOrientation string `json:"customZoomOrientation"`
		CustomZoomPageHeight  int    `json:"customZoomPageHeight"`
		CustomZoomPageWidth   int    `json:"customZoomPageWidth"`
		CustomZoomScale       int    `json:"customZoomScale"`
		DocumentMetadata      struct {
		} `json:"documentMetadata"`
		ExtraMetadata struct {
			LastBallpointv2Color string `json:"LastBallpointv2Color"`
			LastBallpointv2Size  string `json:"LastBallpointv2Size"`
			LastEraserColor      string `json:"LastEraserColor"`
			LastEraserSize       string `json:"LastEraserSize"`
			LastEraserTool       string `json:"LastEraserTool"`
			LastPen              string `json:"LastPen"`
			LastTool             string `json:"LastTool"`
		} `json:"extraMetadata"`
		FileType      string `json:"fileType"`
		FontName      string `json:"fontName"`
		FormatVersion int    `json:"formatVersion"`
		LineHeight    int    `json:"lineHeight"`
		Margins       int    `json:"margins"`
		Orientation   string `json:"orientation"`
		PageCount     int    `json:"pageCount"`
		PageTags      []struct {
		} `json:"pageTags"`
		SizeInBytes string `json:"sizeInBytes"`
		Tags        []struct {
		} `json:"tags"`
		TextAlignment string `json:"textAlignment"`
		TextScale     int    `json:"textScale"`
		ZoomMode      string `json:"zoomMode"`
	}{
		CPages: struct {
			LastOpened struct {
				Timestamp string `json:"timestamp"`
				Value     string `json:"value"`
			} `json:"lastOpened"`
			Original struct {
				Timestamp string `json:"timestamp"`
				Value     int    `json:"value"`
			} `json:"original"`
			Pages []struct {
				ID  string `json:"id"`
				Idx struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"idx"`
				Template struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"template"`
			} `json:"pages"`
			UUIDs []struct {
				First  string `json:"first"`
				Second int    `json:"second"`
			} `json:"uuids"`
		}{
			LastOpened: struct {
				Timestamp string `json:"timestamp"`
				Value     string `json:"value"`
			}{
				Timestamp: "1:1",
				Value:     pageID,
			},
			Original: struct {
				Timestamp string `json:"timestamp"`
				Value     int    `json:"value"`
			}{
				Timestamp: "1:1",
				Value:     -1,
			},
			Pages: []struct {
				ID  string `json:"id"`
				Idx struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"idx"`
				Template struct {
					Timestamp string `json:"timestamp"`
					Value     string `json:"value"`
				} `json:"template"`
			}{
				{
					ID: pageID,
					Idx: struct {
						Timestamp string `json:"timestamp"`
						Value     string `json:"value"`
					}{
						Timestamp: "1:2",
						Value:     "ba",
					},
					Template: struct {
						Timestamp string `json:"timestamp"`
						Value     string `json:"value"`
					}{
						Timestamp: "1:2",
						Value:     "Blank",
					},
				},
			},
			UUIDs: []struct {
				First  string `json:"first"`
				Second int    `json:"second"`
			}{
				{
					First:  "25248a5b-7602-5a83-b6b8-885ee4e4f813",
					Second: 1,
				},
			},
		},
		CoverPageNumber:       -1,
		CustomZoomCenterX:     0,
		CustomZoomCenterY:     936,
		CustomZoomOrientation: "portrait",
		CustomZoomPageHeight:  1872,
		CustomZoomPageWidth:   1404,
		CustomZoomScale:       1,
		DocumentMetadata:      struct{}{},
		ExtraMetadata: struct {
			LastBallpointv2Color string `json:"LastBallpointv2Color"`
			LastBallpointv2Size  string `json:"LastBallpointv2Size"`
			LastEraserColor      string `json:"LastEraserColor"`
			LastEraserSize       string `json:"LastEraserSize"`
			LastEraserTool       string `json:"LastEraserTool"`
			LastPen              string `json:"LastPen"`
			LastTool             string `json:"LastTool"`
		}{
			LastBallpointv2Color: "Black",
			LastBallpointv2Size:  "2",
			LastEraserColor:      "Black",
			LastEraserSize:       "2",
			LastEraserTool:       "Eraser",
			LastPen:              "Ballpointv2",
			LastTool:             "Ballpointv2",
		},
		FileType:      "notebook",
		FontName:      "",
		FormatVersion: 2,
		LineHeight:    -1,
		Margins:       125,
		Orientation:   "portrait",
		PageCount:     1,
		PageTags:      []struct{}{},
		SizeInBytes:   "0",
		Tags:          []struct{}{},
		TextAlignment: "justify",
		TextScale:     1,
		ZoomMode:      "bestFit",
	}

	contentJSON, err := json.MarshalIndent(content, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling content JSON: %v", err)
	}

	return string(contentJSON)
}

func (rmdoc *ReMarkableAPIrmdoc) createNotebookMetadata(visibleName string) string {
	notebookMetadata := struct {
		CreatedTime    int64  `json:"createdTime"`
		LastModified   int64  `json:"lastModified"`
		LastOpened     int64  `json:"lastOpened"`
		LastOpenedPage int    `json:"lastOpenedPage"`
		Parent         string `json:"parent"`
		Pinned         bool   `json:"pinned"`
		Type           string `json:"type"`
		VisibleName    string `json:"visibleName"`
	}{
		CreatedTime:    rmdoc.Time,
		LastModified:   rmdoc.Time,
		LastOpened:     rmdoc.Time,
		LastOpenedPage: 0,
		Parent:         "",
		Pinned:         true,
		Type:           "DocumentType",
		VisibleName:    visibleName,
	}

	notebookMetadataJSON, err := json.MarshalIndent(notebookMetadata, "", "    ")
	if err != nil {
		log.Fatalf("Error marshaling notebook metadata JSON: %v", err)
	}

	return string(notebookMetadataJSON)
}

func CreateRmDoc(rmName, dirToSave string, rmData []byte) (*bytes.Buffer, string) {
	var zipName string

	if strings.HasSuffix(rmName, ".rm") {
		zipName = rmName[:len(rmName)-len(".rm")] + ".rmdoc"
	} else {
		zipName = rmName + ".rmdoc"
	}

	rmfiledata := rmData

	rmdoc := NewReMarkableAPIrmdoc(zipName, dirToSave, rmfiledata)
	DebugPrint("File " + zipName + " created successfully.")

	return rmdoc.internalBuffer, zipName
}
