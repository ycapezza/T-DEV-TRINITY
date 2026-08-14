package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	baseURL      string
	accessToken  string
}

func NewClient(clientID, clientSecret string, sandbox bool) *Client {
	baseURL := "https://api-m.paypal.com"
	if sandbox {
		baseURL = "https://api-m.sandbox.paypal.com"
	}

	return &Client{
		httpClient:   &http.Client{},
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      baseURL,
	}
}

func (c *Client) GetAccessToken() error {
	data := strings.NewReader("grant_type=client_credentials")
	req, err := http.NewRequest("POST", c.baseURL+"/v1/oauth2/token", data)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.clientID, c.clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get access token, status: %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return err
	}

	c.accessToken = tokenResp.AccessToken
	return nil
}

func (c *Client) CreateOrder(total float64, invoiceID uint, successURL, cancelURL string) (string, string, error) {
	if c.accessToken == "" {
		if err := c.GetAccessToken(); err != nil {
			return "", "", err
		}
	}

	orderData := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": fmt.Sprintf("%d", invoiceID),
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", total),
				},
			},
		},
		"application_context": map[string]interface{}{
			"return_url": successURL,
			"cancel_url": cancelURL,
		},
	}

	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v2/checkout/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("failed to create order, status: %d", resp.StatusCode)
	}

	var orderResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return "", "", err
	}

	orderID := orderResp["id"].(string)

	var approveURL string
	for _, link := range orderResp["links"].([]interface{}) {
		linkMap := link.(map[string]interface{})
		if linkMap["rel"].(string) == "approve" {
			approveURL = linkMap["href"].(string)
			break
		}
	}

	return orderID, approveURL, nil
}

func (c *Client) CaptureOrder(orderID string) (string, error) {
	if c.accessToken == "" {
		if err := c.GetAccessToken(); err != nil {
			return "", err
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/checkout/orders/%s/capture", c.baseURL, orderID), bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to capture order, status: %d", resp.StatusCode)
	}

	var captureResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&captureResp); err != nil {
		return "", err
	}

	captureID := captureResp["purchase_units"].([]interface{})[0].(map[string]interface{})["payments"].(map[string]interface{})["captures"].([]interface{})[0].(map[string]interface{})["id"].(string)

	return captureID, nil
}

func (c *Client) VerifyWebhookSignature(headers map[string]string, body []byte, webhookID string) (bool, error) {
	if c.accessToken == "" {
		if err := c.GetAccessToken(); err != nil {
			return false, err
		}
	}

	verifyData := map[string]interface{}{
		"auth_algo":         headers["PAYPAL-AUTH-ALGO"],
		"cert_url":          headers["PAYPAL-CERT-URL"],
		"transmission_id":   headers["PAYPAL-TRANSMISSION-ID"],
		"transmission_sig":  headers["PAYPAL-TRANSMISSION-SIG"],
		"transmission_time": headers["PAYPAL-TRANSMISSION-TIME"],
		"webhook_id":        webhookID,
		"webhook_event":     json.RawMessage(body),
	}

	jsonData, err := json.Marshal(verifyData)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", c.baseURL+"/v1/notifications/verify-webhook-signature", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var verifyResp struct {
		VerificationStatus string `json:"verification_status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return false, err
	}

	return verifyResp.VerificationStatus == "SUCCESS", nil
}