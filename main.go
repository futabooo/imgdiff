package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const TimeFormat = "2006-01-02-1504"

func main() {
	http.HandleFunc("/upload", uploadHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	fileBytes, fileType, err := checkFile(file)
	if err != nil {
		log.Fatal(err)
	}

	fileName, err := saveFile(fileBytes, fileType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s Upload success", fileName)
}

func checkFile(f multipart.File) ([]byte, string, error) {
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, "", err
	}

	// check file type, detectcontenttype only needs the first 512 bytes
	fileType := http.DetectContentType(fileBytes)
	switch fileType {
	case "image/jpeg", "image/jpg":
		break
	default:
		return nil, "", err
	}

	return fileBytes, fileType, nil
}

func saveFile(b []byte, t string) (string, error) {
	fileName := getFileName()
	fileEndings, err := mime.ExtensionsByType(t)
	if err != nil {
		return "", err
	}
	newPath := fileName + fileEndings[0]

	newFile, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer newFile.Close()
	if _, err := newFile.Write(b); err != nil {
		return "", err
	}

	return newFile.Name(), nil
}

func getFileName() string {
	t := time.Now()
	return t.Format(TimeFormat)
}
