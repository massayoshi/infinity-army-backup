package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var execCommand = exec.Command
var outWriter = os.Stdout

func TestCreateFolder(t *testing.T) {
	tempDir := t.TempDir()

	path := filepath.Join(tempDir, "test-folder")
	createFolder(path)

	_, err := os.Stat(path)
	if err != nil {
		t.Errorf("createFolder failed to create directory: %s", err)
	}
}

func TestCreateFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test-file-*.txt")
	if err != nil {
		t.Fatal("Cannot create temporary file", err)
	}
	defer os.Remove(tempFile.Name())

	data := []byte("{ data: true }")
	createFile(tempFile.Name(), data, true)

	fileData, err := os.ReadFile(tempFile.Name())

	if err != nil {
		t.Errorf("createFile failed to create file: %s", err)
	}
	if !bytes.Equal(fileData, data) {
		t.Errorf("createFile created file with incorrect data: expected '%s', got '%s'", data, fileData)
	}

	newData := []byte("")
	createFile(tempFile.Name(), newData, false)

	updatedFileData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Errorf("createFile failed to update file: %s", err)
	}
	if !bytes.Equal(updatedFileData, data) {
		t.Errorf("createFile failed to update file with new data: expected '%s', got '%s'", data, updatedFileData)
	}
}

func TestHttpClient(t *testing.T) {
	client := httpClient()

	expectedTimeout := 10 * time.Second
	if client.Timeout != expectedTimeout {
		t.Errorf("httpClient() returned a client with timeout %v, expected %v", client.Timeout, expectedTimeout)
	}
}

func TestSendRequest(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != getEnvVar("HTTP_REQUEST_HEADER_ORIGIN") {
			t.Error("Origin header not set correctly")
		}
		if r.Header.Get("Referer") != getEnvVar("HTTP_REQUEST_HEADER_REFERER") {
			t.Error("Referer header not set correctly")
		}
		if r.Header.Get("User-Agent") != getEnvVar("HTTP_REQUEST_HEADER_USER_AGENT") {
			t.Error("User-Agent header not set correctly")
		}
		fmt.Fprintln(w, "Hello, world!")
	}))
	defer testServer.Close()

	client := &http.Client{}

	body := sendRequest(client, testServer.URL)

	expectedBody := "Hello, world!\n"
	if string(body) != expectedBody {
		t.Errorf("sendRequest returned unexpected body. Expected %q but got %q", expectedBody, string(body))
	}
}

func TestGetEnvVar(t *testing.T) {
	key := "MY_ENV_VAR"
	expectedValue := "my test value"
	os.Setenv(key, expectedValue)

	value := getEnvVar(key)
	if value != expectedValue {
		t.Errorf("getEnvVar(%q) returned %q, but expected %q", key, value, expectedValue)
	}

	os.Unsetenv(key)
}

func TestPrettyPrint(t *testing.T) {
	input := []byte(`{"data":"value", "array": [1,2,3]}`)

	output := prettyPrint(input)

	expectedOutput := []byte(`{
  "data": "value",
  "array": [
    1,
    2,
    3
  ]
}`)
	if !bytes.Equal(output, expectedOutput) {
		t.Errorf("prettyPrint(%q) returned %q, but expected %q", input, output, expectedOutput)
	}
}
