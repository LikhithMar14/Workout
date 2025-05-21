package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)


type Envelope map[string]interface{}

func WriteJSON(w http.ResponseWriter, status int, data Envelope)error{
	//This will return a byte slice to the js if there are no errors
	js, err := json.MarshalIndent(data,""," ")
	if err != nil {
		return err
	}
	js = append(js,'\n')
	w.WriteHeader(status)
	w.Header().Set("Content-Type","application/json")
	w.Write(js)
	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		return 0, errors.New("invalid id parameter")
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter type")
	}

	return id, nil
}