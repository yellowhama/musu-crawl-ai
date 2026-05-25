package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type OllamaClient struct {
	BaseURL     string
	Model       string
	VisionModel string
}

func NewOllamaClient(model string) *OllamaClient {
	return &OllamaClient{
		BaseURL:     "http://localhost:11434/api/generate",
		Model:       model,
		VisionModel: "llava", // Default vision model
	}
}

type OllamaRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Stream bool     `json:"stream"`
	Format string   `json:"format,omitempty"`
	Images []string `json:"images,omitempty"` // Base64 encoded images
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type EmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbedResponse struct {
	Embedding []float64 `json:"embedding"`
}

func (c *OllamaClient) Ask(prompt string, jsonFormat bool) (string, error) {
	if c.BaseURL == "" {
		return "", fmt.Errorf("ollama not configured")
	}

	reqBody := OllamaRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}
	if jsonFormat {
		reqBody.Format = "json"
	}

	jsonData, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(c.BaseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ollama connection failed: %v. if you are an AI agent, you should perform this reasoning task yourself using the provided data", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama returned error %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode ollama response: %v", err)
	}

	return ollamaResp.Response, nil
}

func (c *OllamaClient) Embed(text string) ([]float64, error) {
	reqBody := EmbedRequest{
		Model:  c.Model,
		Prompt: text,
	}

	jsonData, _ := json.Marshal(reqBody)
	url := "http://localhost:11434/api/embeddings"
	
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama embed failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama embed error %d: %s", resp.StatusCode, string(body))
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode ollama embed response: %v", err)
	}

	return embedResp.Embedding, nil
}

func (c *OllamaClient) DescribeImage(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)

	reqBody := OllamaRequest{
		Model:  c.VisionModel,
		Prompt: "Describe this image in one concise sentence for a knowledge base index.",
		Stream: false,
		Images: []string{b64},
	}

	jsonData, _ := json.Marshal(reqBody)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(c.BaseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("ollama vision failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama vision error %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode ollama vision response: %v", err)
	}

	return ollamaResp.Response, nil
}
