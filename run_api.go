package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

// // get code steps
func Coderrunner(r *http.Request, enc Encoder, parms martini.Params) (int, string) {
	// Otherwise, return all Codes
	id, err := getPostRun(r)
	if err != nil {
		return http.StatusNotAcceptable, Must(enc.Encode(
			NewError(ErrWrongInput, fmt.Sprintf("wrong input"))))
	}
	_, ok := run_map[id.Id]
	if ok == false {
		go run(id)
		log.Println("running.....")
	}
	return http.StatusOK, Must(enc.Encode(NewError(CommitOk, fmt.Sprintf("run id=%s is running", id.Id))))
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
	res2B, _ := json.Marshal(t)
	fmt.Println(string(res2B))

	if err != nil {
		return nil, err
	}
	return &t, nil
}
