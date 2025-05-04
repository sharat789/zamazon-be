package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TokenUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	UserRole string `json:"user_role"`
}

type AuthClient struct {
	BaseURL string
}

func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		BaseURL: baseURL,
	}
}

func (c *AuthClient) CreateHashPassword(password string) (string, error) {
	requestBody, err := json.Marshal(map[string]string{
		"password": password,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/hash-password", c.BaseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to hash password: %d", resp.StatusCode)
	}

	var response struct {
		HashedPassword string `json:"hashed_password"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.HashedPassword, nil
}

func (c *AuthClient) VerifyPassword(plainPassword, hashedPassword string) error {
	requestBody, err := json.Marshal(map[string]string{
		"plain_password":  plainPassword,
		"hashed_password": hashedPassword,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/verify-password", c.BaseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("password verification failed")
	}

	return nil
}

func (c *AuthClient) GenerateToken(id uint, email, role string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"id":    id,
		"email": email,
		"role":  role,
	})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/generate-token", c.BaseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate token: %d", resp.StatusCode)
	}

	var response struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *AuthClient) VerifyToken(token string) (*TokenUser, error) {
	requestBody, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/verify-token", c.BaseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("token verification failed")
	}

	var response struct {
		User TokenUser `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.User, nil
}

func (c *AuthClient) GenerateCode() (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/auth/generate-code", c.BaseURL))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to generate code")
	}

	var response struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Code, nil
}

func (c *AuthClient) AuthorizeByRole(token, role string) (*TokenUser, error) {
	requestBody, err := json.Marshal(map[string]string{
		"token": token,
		"role":  role,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/auth/authorize-by-role", c.BaseURL),
		"application/json",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("authorization failed")
	}

	var response struct {
		User TokenUser `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.User, nil
}
