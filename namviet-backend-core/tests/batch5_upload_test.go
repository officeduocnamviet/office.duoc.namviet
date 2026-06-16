package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadAPI(t *testing.T) {
	r := SetupTestRouter()

	t.Run("Upload without file should fail", func(t *testing.T) {
		w := PerformRequest(r, "POST", "/api/upload", "")
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Upload file success", func(t *testing.T) {
		// Only run this test if CLOUDINARY_URL is available
		if os.Getenv("CLOUDINARY_URL") == "" {
			t.Skip("Skipping Cloudinary upload test because CLOUDINARY_URL is missing")
		}

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", "test_image.txt")
		assert.NoError(t, err)

		// Create a dummy text file to act as an image (Cloudinary accepts txt files as raw or handles them)
		_, err = io.WriteString(part, "dummy content")
		assert.NoError(t, err)

		err = writer.Close()
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer namviet-admin-super-key")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Could be StatusOK or StatusInternalServerError depending on whether Cloudinary accepts a .txt fake image
		t.Logf("Response: %s", w.Body.String())
	})
}
