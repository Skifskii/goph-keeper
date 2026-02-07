package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type APIClient struct {
	BaseURL string
	Token   string
}

func NewClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

func (c *APIClient) Register(username, password string) error {
	reqBody := RegisterReq{
		Username: username,
		Password: password,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshall request body: %w", err)
	}

	resp, err := http.Post(c.BaseURL+"/api/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	return nil
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *APIClient) Login(username, password string) error {
	reqBody := LoginReq{
		Username: username,
		Password: password,
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshall request body: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/api/login",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "jwt" {
			c.Token = cookie.Value
			return nil
		}
	}

	return nil
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SecretMeta struct {
	ID         int
	SecretType string
	Metadata   string
}

func (c *APIClient) ListSecrets(limit, offset int) ([]SecretMeta, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/secret/baselist/?limit=%d&offset=%d", c.BaseURL, limit, offset),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: c.Token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list secrets failed: %s", resp.Status)
	}

	var secrets []SecretMeta
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return secrets, nil
}
