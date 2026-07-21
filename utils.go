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
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// maxConcurrentRequests bounds how many HTTP requests are in flight at once
// across the whole program, keeping load on the upstream APIs polite.
const maxConcurrentRequests = 8

// reqSem is a global semaphore acquired for the duration of every outbound
// request in sendRequest, so total concurrency is capped regardless of how
// many goroutines the army/wiki/logo loops spawn.
var reqSem = make(chan struct{}, maxConcurrentRequests)

// folderMu guards createFolder so concurrent goroutines creating the same
// folder (e.g. a shared wiki category) don't race on stat/mkdir.
var folderMu sync.Mutex

func createFolder(path string) {
	folderMu.Lock()
	defer folderMu.Unlock()
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
}

// sharedClient is created once and reused for every request so the connection
// pool (keep-alives) is shared across all goroutines.
var (
	clientOnce   sync.Once
	sharedClient *http.Client
)

func httpClient() *http.Client {
	clientOnce.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:        maxConcurrentRequests * 2,
			MaxIdleConnsPerHost: maxConcurrentRequests,
			IdleConnTimeout:     90 * time.Second,
		}
		sharedClient = &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		}
	})
	return sharedClient
}

func sendRequest(client *http.Client, endpoint string) []byte {
	reqSem <- struct{}{}
	defer func() { <-reqSem }()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
	}

	req.Header.Add("origin", getEnvVar("HTTP_REQUEST_HEADER_ORIGIN"))
	req.Header.Add("referer", getEnvVar("HTTP_REQUEST_HEADER_REFERER"))
	req.Header.Add("user-agent", getEnvVar("HTTP_REQUEST_HEADER_USER_AGENT"))
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "cross-site")

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

var envOnce sync.Once

func getEnvVar(value string) string {
	envOnce.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Some error occured. Err: %s", err)
		}
	})

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
