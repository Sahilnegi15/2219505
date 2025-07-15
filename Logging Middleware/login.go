package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RegisterRequest struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	MobileNo       string `json:"mobileNo"`
	GithubUsername string `json:"githubUsername"`
	RollNo         string `json:"rollNo"`
	AccessCode     string `json:"accessCode"`
}

type AuthRequest struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	RollNo       string `json:"rollNo"`
	AccessCode   string `json:"accessCode"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

type AuthResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type LogRequest struct {
	Stack   string `json:"stack"`
	Level   string `json:"level"`
	Package string `json:"package"`
	Message string `json:"message"`
}

func main() {

	if err := registerUser(); err != nil {
		fmt.Println("Registration:", err)
	}

	token, err := getAuthToken()
	if err != nil {
		fmt.Println("Auth failed:", err)
		return
	}
	fmt.Println("Token acquired!")

	err = sendLog(token, "backend", "error", "controller", "received string, expected bool")
	if err != nil {
		fmt.Println("Log failed:", err)
	} else {
		fmt.Println("Log sent successfully.")
	}
}

func registerUser() error {
	url := "http://20.244.56.144/evaluation-service/register"

	body := RegisterRequest{
		Email:          "ayushnegi0220@gmail.com",
		Name:           "Sahil Negi",
		MobileNo:       "6398687035",
		GithubUsername: "Sahilnegi15",
		RollNo:         "2219505",
		AccessCode:     "QAhDUr",
	}

	respBody, status, err := postJSON(url, body, "")
	if err != nil {
		return err
	}

	if status == 200 {
		fmt.Println("User registered successfully.")
	} else if status == 400 && strings.Contains(string(respBody), "You can register only once") {
		fmt.Println("User already registered. Proceeding.")
	} else {
		return fmt.Errorf("registration failed: %s", string(respBody))
	}

	return nil
}

func getAuthToken() (string, error) {
	url := "http://20.244.56.144/evaluation-service/auth"

	body := AuthRequest{
		Email:        "ayushnegi0220@gmail.com",
		Name:         "Sahil Negi",
		RollNo:       "2219505",
		AccessCode:   "QAhDUr",
		ClientID:     "5f09afac-5aec-4211-a4d7-f25848e4ba20",
		ClientSecret: "sBmsJSGFYndJdfKD",
	}

	respBody, status, err := postJSON(url, body, "")
	if err != nil {
		return "", err
	}

	if status != 200 {
		return "", fmt.Errorf("auth failed: %s", string(respBody))
	}

	var authResp AuthResponse
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return "", fmt.Errorf("unable to parse auth response: %v", err)
	}

	return strings.TrimSpace(authResp.AccessToken), nil
}

func sendLog(token, stack, level, pkg, message string) error {
	url := "http://20.244.56.144/evaluation-service/logs"

	body := LogRequest{
		Stack:   stack,
		Level:   level,
		Package: pkg,
		Message: message,
	}

	respBody, status, err := postJSON(url, body, token)
	if err != nil {
		return err
	}

	if status != 200 {
		fmt.Printf("Log response (%d): %s\n", status, string(respBody))
		return fmt.Errorf("log failed with status %d", status)
	}

	return nil
}

func postJSON(url string, body interface{}, token string) ([]byte, int, error) {
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	return respBody, resp.StatusCode, nil
}
