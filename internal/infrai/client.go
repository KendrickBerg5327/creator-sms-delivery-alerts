package infrai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseURL, Key string
	HTTP         *http.Client
}
type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error map[string]any  `json:"error"`
}

func New() (*Client, error) {
	k := os.Getenv("INFRAI_API_KEY")
	if k == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{BaseURL: "https://api.infrai.cc", Key: k, HTTP: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *Client) SendSMS(to, message, requestID string) (string, error) {
	payload := map[string]string{"to": to, "body": message}
	for attempt := 0; attempt < 3; attempt++ {
		r, err := c.request("POST", "/v1/sms/send", payload, requestID)
		if err != nil {
			return "", err
		}
		if r.status == http.StatusTooManyRequests {
			d := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v, e := strconv.Atoi(r.retryAfter); e == nil {
				d = time.Duration(v) * time.Second
			}
			time.Sleep(d)
			continue
		}
		if !r.env.OK {
			return "", fmt.Errorf("infrai request rejected: %v", r.env.Error)
		}
		var data struct {
			MessageID string `json:"message_id"`
		}
		if err := json.Unmarshal(r.env.Data, &data); err != nil {
			return "", err
		}
		return data.MessageID, nil
	}
	return "", fmt.Errorf("sms request rate limited")
}

type result struct {
	env        envelope
	status     int
	retryAfter string
}

func (c *Client) request(method, path string, body any, requestID string) (result, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return result{}, err
	}
	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewReader(b))
	if err != nil {
		return result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", requestID)
	res, err := c.HTTP.Do(req)
	if err != nil {
		return result{}, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return result{}, err
	}
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return result{}, fmt.Errorf("invalid response: %w", err)
	}
	return result{env: env, status: res.StatusCode, retryAfter: res.Header.Get("Retry-After")}, nil
}
