package kv

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tPhume/gokv/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	happyTestBody     = store.Value{"value": "Hello, I am a value"}
	newHappyTestBody  = store.Value{"newValue": "Hello, I am the new and better value"}
	badFormatTestBody = "this is wrong :("
	router            *gin.Engine
)

func setUp() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	DefaultRestWithEngine(router)
}

func TestHappyPath(t *testing.T) {
	setUp()

	// insert
	body, _ := json.Marshal(happyTestBody)
	req, _ := http.NewRequest("POST", "/store/v1/test", bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody := make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "test inserted", resBody["message"])

	// search
	req, _ = http.NewRequest("GET", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, map[string]string(happyTestBody), resBody)

	// update
	body, _ = json.Marshal(newHappyTestBody)
	req, _ = http.NewRequest("PATCH", "/store/v1/test", bytes.NewBuffer(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test updated", resBody["message"])

	// search new value
	req, _ = http.NewRequest("GET", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, map[string]string(newHappyTestBody), resBody)

	// delete the value
	req, _ = http.NewRequest("DELETE", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusOK, w.Code)

	// deleted value cannot be found
	req, _ = http.NewRequest("GET", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBadPath(t *testing.T) {
	setUp()

	// insert no body
	req, _ := http.NewRequest("POST", "/store/v1/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody := make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorValueEmpty, resBody["message"])

	// insert bad json
	req, _ = http.NewRequest("POST", "/store/v1/test", bytes.NewBufferString(badFormatTestBody))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorBadJSON, resBody["message"])

	// update key not found
	body, _ := json.Marshal(happyTestBody)
	req, _ = http.NewRequest("PATCH", "/store/v1/test", bytes.NewBuffer(body))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, errorKeyNotFound, resBody["message"])

	// update no body
	req, _ = http.NewRequest("PATCH", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorValueEmpty, resBody["message"])

	// insert bad json
	req, _ = http.NewRequest("PATCH", "/store/v1/test", bytes.NewBufferString(badFormatTestBody))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, errorBadJSON, resBody["message"])

	// search key not found
	req, _ = http.NewRequest("GET", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, errorKeyNotFound, resBody["message"])

	// delete key not found
	req, _ = http.NewRequest("DELETE", "/store/v1/test", nil)

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resBody = make(map[string]string)
	_ = json.Unmarshal(w.Body.Bytes(), &resBody)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, errorKeyNotFound, resBody["message"])
}
