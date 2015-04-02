package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"strconv"
	"time"
)

// // get code steps
func Coderrunner(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	// Otherwise, return all Codes
	id, err := getPostRun(r)
	if err != nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("wrong input"))))
	}

	if run_map[id.Id] == nil {
		//handler with command
	}
	return http.StatusOK, fmt.Sprintf("run id=%s is running", id.Id)
}

// Parse the request body, load into an Code structure.
func getPostRun(r *http.Request) (*Run, error) {
	decoder := json.NewDecoder(r.Body)
	var t Run
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	log.Println(t)

	if err != nil {
		return nil, err
	}
	return &t, nil
}
