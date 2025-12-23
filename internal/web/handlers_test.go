package web

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func createMultiPartRequest(t *testing.T, files map[string]string) *http.Request {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for filename, content := range files {
		part, err := writer.CreateFormFile("files", filename)
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}

		_, err = part.Write([]byte(content))
		if err != nil {
			t.Fatalf("failed to write form data for %s: %v", filename, err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/tools/concurrent-upper", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestConcurrentUpper_Success(t *testing.T) {
	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Change to project root (two levels up from internal/web)
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd) // Restore working directory after test

	app, err := NewApplication(
		log.New(os.Stdout, "TEST INFO: ", 0), // See logs!
		log.New(os.Stdout, "TEST ERROR: ", 0),
	)

	if err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"test1.txt": "hello world",
		"test2.txt": "goodbye world",
	}

	req := createMultiPartRequest(t, files)
	recorder := httptest.NewRecorder()

	app.concurrentUpper(recorder, req)

	responseBody := recorder.Body.String()

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	if !strings.Contains(responseBody, "HELLO WORLD") {
		t.Errorf("expected uppercase 'HELLO WORLD' in response")
	}

	if !strings.Contains(responseBody, "GOODBYE WORLD") {
		t.Error("expected uppercase 'GOODBYE WORLD' in response")
	}

	if strings.Contains(responseBody, "Error:") {
		t.Errorf("unexpected error in response")
	}
}

func TestConcurrentUpper_NoFiles(t *testing.T) {
	// Save current working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Change to project root (two levels up from internal/web)
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd) // Restore working directory after test

	app, err := NewApplication(
		log.New(os.Stdout, "TEST INFO: ", 0), // See logs!
		log.New(os.Stdout, "TEST ERROR: ", 0),
	)
	if err != nil {
		t.Fatal(err)
	}

	files := map[string]string{}

	req := createMultiPartRequest(t, files)
	recorder := httptest.NewRecorder()

	app.concurrentUpper(recorder, req)

	responseBody := recorder.Body.String()

	if recorder.Code != 400 {
		t.Errorf("expected status of 400, got %d", recorder.Code)
	}

	if !strings.Contains(responseBody, "Error: ") {
		t.Errorf("expected error in response")
	}
}

// func TestConcurrentUpper_FileTooLarge(t *testing.T) {
// 	// Save current working directory
// 	originalWd, err := os.Getwd()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Change to project root (two levels up from internal/web)
// 	if err := os.Chdir("../.."); err != nil {
// 		t.Fatal(err)
// 	}
// 	defer os.Chdir(originalWd) // Restore working directory after test

// 	app, err := NewApplication(
// 		log.New(os.Stdout, "TEST INFO: ", 0), // See logs!
// 		log.New(os.Stdout, "TEST ERROR: ", 0),
// 	)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
