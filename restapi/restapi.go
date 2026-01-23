package restapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type ApiPack struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Token  string `json:"token"`
	Body   any    `json:"body"`
}

func Call(payload ApiPack) ([]byte, error) {
	client := &http.Client{Timeout: time.Second * 10}
	body, _ := json.Marshal(payload.Body)
	req, err := http.NewRequest(
		payload.Method,
		payload.Url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	authToken := "Bearer " + payload.Token
	req.Header.Set("authorization", authToken)
	req.Header.Set("content-type", "application/json; charset=utf-8")
	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(rsp.Status)
	}
	resBody, err := io.ReadAll(rsp.Body)
	return resBody, err
}
