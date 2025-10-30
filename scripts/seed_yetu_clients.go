package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type gqlRequest struct {
	Query string `json:"query"`
}

type authPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func doGraphQL(url string, query string, token string) (map[string]any, error) {
	body, _ := json.Marshal(gqlRequest{Query: query})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, fmt.Errorf("non-2xx status: %d", resp.StatusCode)
	}
	if errs, ok := out["errors"]; ok {
		return out, fmt.Errorf("graphql error: %v", errs)
	}
	return out, nil
}

func loginAndGetToken(apiURL, email, password string) (string, error) {
	q := fmt.Sprintf(`mutation { userLogin(input: { email: "%.200s", password: "%.200s" }) { accessToken } }`, email, password)
	res, err := doGraphQL(apiURL, q, "")
	if err != nil {
		return "", err
	}
	data, _ := res["data"].(map[string]any)
	if data == nil {
		return "", fmt.Errorf("no data in response")
	}
	up, _ := data["userLogin"].(map[string]any)
	if up == nil {
		return "", fmt.Errorf("userLogin missing in response")
	}
	at, _ := up["accessToken"].(string)
	if at == "" {
		return "", fmt.Errorf("accessToken empty")
	}
	return at, nil
}

func createClient(apiURL, token, name, password string) (map[string]any, error) {
	q := fmt.Sprintf(`mutation { clientCreate(input: { name: "%.200s", password: "%.200s" }) { id name clientId } }`, name, password)
	return doGraphQL(apiURL, q, token)
}

func main() {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/query"
	}

	// Admin credentials (override with env if needed)
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@mlm.com"
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123456"
	}

	// Clients password (override with env)
	clientsPassword := os.Getenv("CLIENTS_PASSWORD")
	if clientsPassword == "" {
		clientsPassword = "yetu@2025"
	}

	fmt.Printf("API: %s\n", apiURL)
	fmt.Println("Authenticating as admin...")
	token, err := loginAndGetToken(apiURL, adminEmail, adminPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "login failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("OK: got access token")

	for i := 1; i <= 7; i++ {
		name := fmt.Sprintf("yetu%d", i)
		fmt.Printf("Creating client: %s...\n", name)
		res, err := createClient(apiURL, token, name, clientsPassword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed: %v\n", err)
			continue
		}
		enc, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(enc))
	}

	fmt.Println("Done.")
}
