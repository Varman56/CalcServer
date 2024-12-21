package application

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCalcHandlerBadRequestCase(t *testing.T) {
	type AnswerBad struct {
		Error string `json:"error"`
	}
	testCasesBad := []struct {
		name           string
		data           map[string]string
		expectedResult map[string]string
		wantBadRequest int
	}{
		{
			name:           "division by zero",
			data:           map[string]string{"expression": "24/0"},
			expectedResult: map[string]string{"error": "division by zero"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
		{
			name:           "incorrect count of brackets",
			data:           map[string]string{"expression": "((2+3)"},
			expectedResult: map[string]string{"error": "incorrect count of brackets"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
		{
			name:           "multiple operands in a row",
			data:           map[string]string{"expression": "2++3"},
			expectedResult: map[string]string{"error": "multiple operands in a row"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
		{
			name:           "invalid expression",
			data:           map[string]string{"expression": ""},
			expectedResult: map[string]string{"error": "invalid expression"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
		{
			name:           "failure to convert to float64",
			data:           map[string]string{"expression": "2+2..2"},
			expectedResult: map[string]string{"error": "failure to convert to float64"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
		{
			name:           "undefined operand",
			data:           map[string]string{"expression": "2&3"},
			expectedResult: map[string]string{"error": "undefined operand"},
			wantBadRequest: http.StatusUnprocessableEntity,
		},
	}
	for _, testCase := range testCasesBad {
		jsonValue, _ := json.Marshal(testCase.data)
		req := httptest.NewRequest(http.MethodPost, "localhost:8080/api/v1/calculate", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		CalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("Test: %s\nhandler returned wrong status code: got %v want %v", testCase.name, res.StatusCode, http.StatusUnprocessableEntity)
		} else {
			var reqResult AnswerBad
			err := json.Unmarshal(w.Body.Bytes(), &reqResult)
			if err != nil {
				t.Fatalf("Test: %s\npanic while unmarshal answer: %v", testCase.name, w.Body.String())
			}
			if reqResult.Error != testCase.expectedResult["error"] {
				jsonWantString, _ := json.Marshal(testCase.expectedResult)
				t.Fatalf("Test: %s\nhandler returned wrong answer: got %v want %v", testCase.name, w.Body.String(), string(jsonWantString))
			}
		}
	}
}

func TestCalcHandlerSuccessCase(t *testing.T) {
	type AnswerOk struct {
		Result float64 `json:"result"`
	}
	testCasesBad := []struct {
		name           string
		data           map[string]string
		expectedResult map[string]string
		wantBadRequest int
	}{
		{
			name:           "simple",
			data:           map[string]string{"expression": "1+1"},
			expectedResult: map[string]string{"result": "2"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "priority_brackets",
			data:           map[string]string{"expression": "(2+2)*2"},
			expectedResult: map[string]string{"result": "8"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "priority",
			data:           map[string]string{"expression": "2+2*2"},
			expectedResult: map[string]string{"result": "6"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "/",
			data:           map[string]string{"expression": "1/2"},
			expectedResult: map[string]string{"result": "0.5"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "brackets 1",
			data:           map[string]string{"expression": "(0)-(1)"},
			expectedResult: map[string]string{"result": "-1"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "brackets 2",
			data:           map[string]string{"expression": "(-1)-(1)"},
			expectedResult: map[string]string{"result": "-2"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "brackets 3",
			data:           map[string]string{"expression": "(-1-(-1))"},
			expectedResult: map[string]string{"result": "0"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t1",
			data:           map[string]string{"expression": "(-1)+(0)"},
			expectedResult: map[string]string{"result": "-1"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t2",
			data:           map[string]string{"expression": "-(-(-1)+0)"},
			expectedResult: map[string]string{"result": "-1"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t3",
			data:           map[string]string{"expression": "(-(-0)+1)"},
			expectedResult: map[string]string{"result": "1"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t4",
			data:           map[string]string{"expression": "((3/2)*2)*(7/7*5)"},
			expectedResult: map[string]string{"result": "15"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t5",
			data:           map[string]string{"expression": "((3/2)*2)*(7/7-7)"},
			expectedResult: map[string]string{"result": "-18"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "t6",
			data:           map[string]string{"expression": "((22.2/2)*3)*(-7)"},
			expectedResult: map[string]string{"result": "-233.1"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "big expression",
			data:           map[string]string{"expression": "15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))*3-(2+(1+1))*(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))"},
			expectedResult: map[string]string{"result": "-30.0721649485"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "added mult",
			data:           map[string]string{"expression": "(7)(5)"},
			expectedResult: map[string]string{"result": "35"},
			wantBadRequest: http.StatusOK,
		},
		{
			name:           "big expression without additional mult",
			data:           map[string]string{"expression": "15/(7-(1+1))*3-(2+(1+1))*15/(7-(200+1))3-(2+(1+1))(15/(7-(1+1))*3-(2+(1+1))+15/(7-(1+1))*3-(2+(1+1)))"},
			expectedResult: map[string]string{"result": "-30.0721649485"},
			wantBadRequest: http.StatusOK,
		},
	}
	const EPS = 1e-9
	for _, testCase := range testCasesBad {
		jsonValue, _ := json.Marshal(testCase.data)
		req := httptest.NewRequest(http.MethodPost, "localhost:8080/api/v1/calculate", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		CalcHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("Test: %s\nhandler returned wrong status code: got %v want %v", testCase.name, res.StatusCode, http.StatusOK)
		} else {
			var reqResult AnswerOk
			err := json.Unmarshal(w.Body.Bytes(), &reqResult)
			if err != nil {
				t.Fatalf("Test: %s\npanic while unmarshal answer: %v", testCase.name, w.Body.String())
			}
			val, _ := strconv.ParseFloat(testCase.expectedResult["result"], 64)
			if math.Abs(reqResult.Result-val) > EPS {
				t.Fatalf("Test: %s\n%f should be equal %f", testCase.name, reqResult.Result, val)
			}
		}
	}
}
