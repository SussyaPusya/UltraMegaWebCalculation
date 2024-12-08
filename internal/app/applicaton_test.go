package app_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SussyaPusya/UltraMegaWebCalculation/internal/app"
)

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

func TestCalcHandlerSuccsesReq(t *testing.T) {
	TestCases := []struct {
		name string
		expr app.JsonReq
	}{
		{
			name: "Simple",
			expr: app.JsonReq{"2+2"},
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

		if res.Body.Close().Error() != "result: 4.00000" {
			t.Fatal(res.Body)
		}
	}

}
