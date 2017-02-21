package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"path/filepath"
	"time"
)

var filenames = []string{
	"Karl1.png",
	"Karl2.png",
	"Karl3.png",
	"Karl4.png",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randomFile() string {
	return RandStringRunes(6) + "-" + filenames[rand.Intn(len(filenames))]
}

func SaveTempFile(b []byte) string {
	file, err := ioutil.TempFile("", "slack_image")
	if err != nil {
		log.Fatalf("error saving file: %s", err)
	}
	if _, err = file.Write(b); err != nil {
		log.Fatalf("error writing file: %s", err)
	}
	if err = file.Close(); err != nil {
		log.Fatalf("error closing file: %s", err)
	}
	return file.Name()
}

func SaveFile(b []byte) string {
	file := randomFile()
	path := filepath.Join(base_path, file)
	ioutil.WriteFile(path, b, 0644)
	return base_url + file
}

func GetFile(file File) []byte {
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	request, err := http.NewRequest(http.MethodGet, file.URLPrivateDownload, nil)
	if err != nil {
		log.Fatalf("error creating request: %s", err)
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("error downloading file\n%v\n%v", file, err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("error downloading file\n%v\n%v", file, err)
	}
	return body
}

func Chrisify(file string) []byte {
	out, err := exec.Command(chrisify, "--haar", haar, "--faces", faces, file).Output()
	if err != nil {
		log.Fatalf("couldn't faceify: %s", err)
	}
	return out
}
