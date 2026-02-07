package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

// APIClient is a lightweight HTTP client used by the CLI/UI to call the
// server API. It stores the base URL and a JWT token for authenticated
// requests.
type APIClient struct {
	BaseURL string
	Token   string
}

// RegisterReq is the JSON payload used to create a new user account.
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SecretMeta contains minimal metadata about a secret as returned
// by the listing endpoint.
type SecretMeta struct {
	ID         int    `json:"id"`
	SecretType string `json:"secret_type"`
	Metadata   string `json:"metadata"`
}

// ListSecrets retrieves a paginated list of secret metadata for the
// authenticated user.
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

type Secret struct {
	SecretType string          `json:"secret_type"`
	Metadata   string          `json:"metadata"`
	Payload    json.RawMessage `json:"payload"`
}

// NewClient creates an APIClient configured to communicate with the
// specified base URL.
func NewClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

// Register sends a registration request for a new user.
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

// Login authenticates with the server and stores the returned JWT in
// the client's Token field.
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

// GetSecret fetches a decrypted secret by its ID.
func (c *APIClient) GetSecret(secretID int) (Secret, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/api/secret/%d", c.BaseURL, secretID),
		nil,
	)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: c.Token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Secret{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Secret{}, fmt.Errorf("list secrets failed: %s", resp.Status)
	}

	var secret Secret
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return Secret{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return secret, nil
}

// CreateSecret sends a request to create a new secret and returns the
// assigned secret ID on success.
func (c *APIClient) CreateSecret(secret Secret) (int, error) {
	data, err := json.Marshal(secret)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal secret: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/api/secret",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// передаём JWT в cookie
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: c.Token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	if !slices.Contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return 0, fmt.Errorf("create secret failed: %s", resp.Status)
	}

	// читаем ответ { "ID": 8 }
	var respBody struct {
		ID int `json:"ID"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	return respBody.ID, nil
}

// UpdateSecret updates an existing secret identified by secretID.
func (c *APIClient) UpdateSecret(secretID int, secret Secret) error {
	data, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/secret/%d", c.BaseURL, secretID),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// передаём JWT в cookie
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: c.Token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send PUT request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update secret failed: %s", resp.Status)
	}

	return nil
}

// DeleteSecret removes a secret identified by secretID.
func (c *APIClient) DeleteSecret(secretID int) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/secret/%d", c.BaseURL, secretID),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// передаём JWT в cookie
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: c.Token,
	})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send DELETE request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusForbidden:
		return fmt.Errorf("access denied")
	case http.StatusNotFound:
		return fmt.Errorf("secret not found")
	default:
		return fmt.Errorf("delete secret failed: %s", resp.Status)
	}
}
