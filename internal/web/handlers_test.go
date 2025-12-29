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

func TestConcurrentUpper_FileTooLarge(t *testing.T) {
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

	// Create a file that exceeds the 10MB limit
	largeContent := strings.Repeat("a", 11*1024*1024) // 11MB
	files := map[string]string{
		"large.txt": largeContent,
	}

	req := createMultiPartRequest(t, files)
	recorder := httptest.NewRecorder()

	app.concurrentUpper(recorder, req)

	responseBody := recorder.Body.String()

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", recorder.Code)
	}

	if !strings.Contains(responseBody, "File too large") {
		t.Errorf("expected 'File too large' error in response, got: %s", responseBody)
	}
}

func createMultiPartRequestForHash(t *testing.T, files map[string]string) *http.Request {
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

	req := httptest.NewRequest("POST", "/tools/concurrent-hash", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func TestConcurrentHash_Success(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir("../.."); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(originalWd)

	app, err := NewApplication(
		log.New(os.Stdout, "TEST INFO: ", 0),
		log.New(os.Stdout, "TEST ERROR: ", 0),
	)
	if err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"test1.txt": "hello",
		"test2.txt": "world",
	}

	req := createMultiPartRequestForHash(t, files)
	recorder := httptest.NewRecorder()

	app.concurrentHash(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	body := recorder.Body.String()

	// SHA-256 of "hello"
	if !strings.Contains(body, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") {
		t.Error("expected hash for 'hello' in response")
	}

	// SHA-256 of "world"
	if !strings.Contains(body, "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7") {
		t.Error("expected hash for 'world' in response")
	}
}
