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
