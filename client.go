package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var filename = flag.String("img", "test-img.jpg", "Path to the image to be uploaded")
var times = flag.Int("times", 10, "Number of times to upload image")

func main() {
	flag.Parse()
	callMany(*times, *filename)
}

func callMany(times int, filename string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Fatalf("error creating form file: [%s]\n", err)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error opening file: [%s]\n", err)
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		log.Fatalf("error copying file contet: [%s]\n", err)
	}
	// content-type of uploaded image
	ct := writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		log.Fatalf("error closing writer: [%s]\n", err)
	}

	forever := make(chan struct{})
	for i := 0; i < times; i++ {
		go func() {
			status, err := call(body.Bytes(), ct)
			if err != nil {
				log.Printf("Request failed with code %d: [%s]\n", status, err)
			}
			log.Printf("Request successfull: [%d]\n", status)
		}()
	}
	<-forever
}

func call(body []byte, ct string) (int, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/upload", bytes.NewReader(body))
	if err != nil {
		return http.StatusBadRequest, err
	}

	req.Header.Set("Content-Type", ct)
	rsp, err := client.Do(req)
	if rsp.StatusCode != http.StatusOK {
		return rsp.StatusCode, err
	}

	return http.StatusOK, nil
}
