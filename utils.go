package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

	if fileName == "" {
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

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func sendRequest(client *http.Client, endpoint string) []byte {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Add("Origin", "https://infinitytheuniverse.com")
	req.Header.Add("Referer", "https://infinitytheuniverse.com")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return body
}

func getEnvVar(value string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	return os.Getenv(value)
}

func showFinalMessage() {
	cmdStruct := exec.Command("git", "status")
	out, err := cmdStruct.Output()
	if err != nil {
		fmt.Println(err)
	}
	if strings.Contains(string(out), "nothing to commit") {
		fmt.Println("Nothing has been changed")
		os.Exit(0)
	} else {
		fmt.Println("Changes have been made")
		fmt.Println(string(out))
	}
}

func prettyPrint(data []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	return out.Bytes()
}
