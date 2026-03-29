package main

import (
	"gross/db"
	"gross/templates"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupTest(t *testing.T) {
	t.Helper()
	db.Init(":memory:")
	templates.Load()
}

func TestIndexHandler(t *testing.T) {
	setupTest(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	indexHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Feeds") {
		t.Error("expected page to contain 'Feeds'")
	}
}

func TestAddFeedHandlerInvalidURL(t *testing.T) {
	setupTest(t)

	form := url.Values{}
	form.Set("url", "not-a-real-feed")

	req := httptest.NewRequest(http.MethodPost, "/feeds", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	addFeedHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
