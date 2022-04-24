package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", w.Body.String())
}

func TestRootHeaderIpRoute(t *testing.T) {
	router := setRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("X-Real-Ip", "127.0.0.1")
	router.ServeHTTP(w, req)

	m := "{\"country_code\":\"\",\"country_name_cn\":\"\",\"country_name_en\":\"\",\"country_name_jp\":\"\",\"ip\":\"127.0.0.1\"}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, m, w.Body.String())
}

func TestRootHeaderNoIpRoute(t *testing.T) {
	router := setRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	// req.Header.Add("X-Real-Ip", "127.0.0.1")
	router.ServeHTTP(w, req)

	m := "{}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, m, w.Body.String())
}

func TestSearchHeaderIpRoute(t *testing.T) {
	router := setRouter()

	w := httptest.NewRecorder()
	params := "?ip=" + url.QueryEscape("127.0.0.1")
	path := "/search" + params
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)

	m := "{\"country_code\":\"\",\"country_name_cn\":\"\",\"country_name_en\":\"\",\"country_name_jp\":\"\",\"ip\":\"127.0.0.1\"}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, m, w.Body.String())
}

func TestSearchHeaderNoIpRoute(t *testing.T) {
	router := setRouter()

	w := httptest.NewRecorder()
	path := "/search"
	req, _ := http.NewRequest("GET", path, nil)
	router.ServeHTTP(w, req)

	m := "{}"
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, m, w.Body.String())
}
