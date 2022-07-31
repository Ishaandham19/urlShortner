package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBaseURL(t *testing.T) {
	resp := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal(err)
	}

	http.DefaultServeMux.ServeHTTP(resp, req)
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		if strings.Contains(string(p), "Error") {
			t.Errorf("header response shouldn't return error: %s", p)
		} else if !strings.Contains(string(p), `API live and kicking`) {
			t.Errorf("header response doesn't match:\n%s", p)
		}
	}
}
