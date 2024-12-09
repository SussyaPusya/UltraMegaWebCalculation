package app_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/app"
)

// Test bad request bla bla bla
func TestBadRequestsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	app.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatal("Invalid status code:", res.StatusCode)
	}

}

// test SuccsesRequest pzdc
func TestCalcHandlerSuccsesReq(t *testing.T) {
	TestCases := []struct {
		name     string
		expr     app.JsonReq
		expected string
	}{
		{
			name:     "Simple",
			expr:     app.JsonReq{"2+2"},
			expected: "result: 4.000000",
		},
		{
			name:     "Medium",
			expr:     app.JsonReq{"2*10/6"},
			expected: "result: 3.333333",
		},
		{
			name:     "Hard",
			expr:     app.JsonReq{"((10*5)-2*10/6)"},
			expected: "result: 46.666667",
		},
	}

	for _, tc := range TestCases {
		body, err := json.Marshal(tc.expr)
		if err != nil {
			t.Fatal(err)
		}
		r := bytes.NewReader(body)

		req := httptest.NewRequest(http.MethodGet, "/", r)
		w := httptest.NewRecorder()
		app.CalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		result, err := ioutil.ReadAll(res.Body) // ioutil ругается так как устарел но пофек

		if string(result) != tc.expected {
			t.Fatalf("Error invalid result %s, want result: 4.000000", string(result))
		}

	}

}
