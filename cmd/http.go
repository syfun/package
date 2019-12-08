package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON map[string]interface{}

type Response struct {
	*http.Response
}

func (r *Response) Decode(v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func (r *Response) Json() (JSON, error) {
	var j JSON
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		return nil, err
	}
	return j, nil
}

func (r *Response) Error() string {
	data, err := r.Json()
	check(err)
	return fmt.Sprintln(data)
}

func Post(url string, data interface{}) (*Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(b))
	return &Response{res}, err
}

func Get(url string) (*Response, error) {
	resp, err := http.Get(url)
	return &Response{resp}, err
}
