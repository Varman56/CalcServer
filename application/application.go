package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Varman56/CalcServer.git/pkg/calculator"
)

type Request struct {
	Expression string `json:"expression"`
}

type AnswerOk struct {
	Result float64 `json:"result"`
}

type AnswerBad struct {
	Error string `json:"error"`
}

var (
	ErrInvalidInput = errors.New("invalid json request")
	ErrServer       = errors.New("internal server error")
	ErrPartsWrtie   = errors.New("wrtied only part of data")
)

var errorsToCheck = []error{
	calculator.ErrInvalidExpression,
	calculator.ErrDivisionByZero,
	calculator.ErrIncorrectBracketSequence,
	calculator.ErrMultipleOperands,
	calculator.ErrConvertingToFloat64,
	calculator.ErrUndefinedOperand,
}

func TryMarshalError(e error) ([]byte, int) {
	res := AnswerBad{Error: e.Error()}
	jsonBytes, err_dec := json.Marshal(res)
	if err_dec != nil {
		ans := AnswerBad{Error: ErrServer.Error()}
		errBytes, _ := json.Marshal(ans)
		return errBytes, http.StatusInternalServerError
	}
	return jsonBytes, http.StatusBadRequest
}

func TryMarshalData(num float64) ([]byte, int) {
	res := AnswerOk{Result: num}
	jsonBytes, err_dec := json.Marshal(res)
	if err_dec != nil {
		ans := AnswerBad{Error: ErrServer.Error()}
		errBytes, _ := json.Marshal(ans)
		return errBytes, http.StatusInternalServerError
	}
	return jsonBytes, -1
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		jsonBytes, status := TryMarshalError(ErrInvalidInput)
		http.Error(w, string(jsonBytes), status)
		return
	}

	result, err := calculator.Calc(request.Expression)
	if err != nil {
		for _, errToCheck := range errorsToCheck {
			if errors.Is(err, errToCheck) {
				jsonBytes, status := TryMarshalError(err)
				http.Error(w, string(jsonBytes), status)
				return
			}
		}
		jsonBytes, status := TryMarshalError(ErrServer)
		http.Error(w, string(jsonBytes), status)
		return
	}
	jsonBytes, status := TryMarshalData(result)
	if status != -1 {
		http.Error(w, string(jsonBytes), status)
		return
	}
	n, err := w.Write(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n != len(jsonBytes) {
		http.Error(w, ErrPartsWrtie.Error(), http.StatusInternalServerError)
		return
	}
}

func RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":8080", nil)
}
