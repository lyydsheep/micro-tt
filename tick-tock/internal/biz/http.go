package biz

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"tick-tock/pkg/log"
)

func Do(ctx context.Context, url string, method, body string, header map[string]string, resp any) error {
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(body)))
	if err != nil {
		log.Error(ctx, "new request error: %v", "error", err)
		return err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(ctx, "do request error: %v", "error", err)
		return err
	}
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(ctx, "read response body error: %v", "error", err)
		return err
	}
	return json.Unmarshal(respBody, resp)
}
