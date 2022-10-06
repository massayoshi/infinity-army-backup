package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func createFolder(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func createFile(fileName string, data []byte, update bool) {
	if _, err := os.Stat(fileName); !errors.Is(err, os.ErrNotExist) && !update {
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatal("Cannot write to file", err)
	}
	file.Sync()
}

func getHTTPResponse(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseData
}
