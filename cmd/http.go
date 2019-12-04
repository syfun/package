package cmd

import (
	"bytes"
	"encoding/json"
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
