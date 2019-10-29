package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	expTime       = 60
	longURL       = "https://www.tracy.com"
	shortLink     = "2"
	shortLinkInfo = `{"url":"https://www.tracy.com","created":"2019-10-29 11:10:46.196556 +0800 CST m = +1033.841130980","expiration_in_minutes":10}`
)

type storageMock struct {
	mock.Mock
}

var app App
var mockR *storageMock

func (s *storageMock) Shorten(url string, exp int64) (string, error) {
	args := s.Called(url, exp)
	return args.String(0), args.Error(1)
}

func (s *storageMock) Unshort(eid string) (string, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func (s *storageMock) ShortLinkInfo(eid string) (interface{}, error) {
	args := s.Called(eid)
	return args.String(0), args.Error(1)
}

func init() {
	app = App{}
	mockR = new(storageMock)
	app.Initialize(&Env{S: mockR})
}

func TestCreateShortLink(t *testing.T) {
	var jsonStr = []byte(`{
		"url":"https://www.tracy.com",
		"expiration_in_minutes":60}`)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal("Should be able to create a request", req)
	}

	req.Header.Set("Content-Type", "application/json")

	mockR.On("Shorten", longURL, int64(expTime)).Return(shortLink, nil).Once()
	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusCreated, rw.Code)
	}

	resp := struct {
		ShortLink string `json:"short_link"`
	}{}

	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("Should decode the response. ")
	}

	if resp.ShortLink != shortLink {
		t.Fatalf("Excepted receive %s. Got %s", shortLink, resp.ShortLink)
	}
}

func TestRedirect(t *testing.T) {
	r := fmt.Sprintf("/%s", shortLink)
	req, err := http.NewRequest("GET", r, nil)
	if err != nil {
		t.Fatal("Should be able to create a request", req)
	}

	mockR.On("Unshort", shortLink).Return(longURL, nil).Once()
	rw := httptest.NewRecorder()
	app.Router.ServeHTTP(rw, req)

	if rw.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Excepted receive %d. Got %d", http.StatusTemporaryRedirect, rw.Code)
	}
}
