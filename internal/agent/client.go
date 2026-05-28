package agent

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ExecutionTrace represents a single reasoning/action step by an agent.
type ExecutionTrace struct {
	Timestamp string `json:"timestamp"`
	Project   string `json:"project"`
	Role      string `json:"role"`
	Goal      string `json:"goal"`
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	Status    string `json:"status"` // success | error
	Error     string `json:"error,omitempty"`
}

// AgentClient is a generic OpenAI-compatible LLM client with Telemetry.
type AgentClient struct {
	BaseURL     string
	Model       string
	VisionModel string
	WikiDir     string
	Project     string
	httpClient  *http.Client
}

func NewAgentClient(baseURL, model, visionModel, wikiDir, project string) *AgentClient {
	if baseURL == "" { baseURL = "http://localhost:11434/v1" }
	if model == "" { model = "llama3" }
	if visionModel == "" { visionModel = "llava" }
	if wikiDir == "" { wikiDir = "./wiki" }

	return &AgentClient{
		BaseURL:     baseURL,
		Model:       model,
		VisionModel: visionModel,
		WikiDir:     wikiDir,
		Project:     project,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				IdleConnTimeout:     90 * time.Second,
				MaxIdleConnsPerHost: 20,
			},
		},
	}
}

func (c *AgentClient) logTrace(trace ExecutionTrace) {
	date := time.Now().Format("2006-01-02")
	logDir := filepath.Join(c.WikiDir, "telemetry", date)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "telemetry: mkdir %s: %v\n", logDir, err)
		return
	}

	filename := fmt.Sprintf("%s_%s_%d.json", trace.Role, c.Project, time.Now().UnixNano())
	path := filepath.Join(logDir, filename)

	data, _ := json.MarshalIndent(trace, "", "  ")
	if err := os.WriteFile(path, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "telemetry: write %s: %v\n", path, err)
	}
}

// OpenAI Chat Completion structures
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []ChatMessage   `json:"messages"`
	Stream         bool            `json:"stream"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// Ask sends a chat completion request and records a trace.
func (c *AgentClient) Ask(prompt string, jsonFormat bool) (string, error) {
	trace := ExecutionTrace{
		Timestamp: time.Now().Format(time.RFC3339),
		Project:   c.Project,
		Role:      "assistant", // Generic, specific agents should override or wrap this
		Prompt:    prompt,
		Status:    "success",
	}

	reqBody := ChatRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	if jsonFormat {
		reqBody.ResponseFormat = &ResponseFormat{Type: "json_object"}
	}

	jsonData, _ := json.Marshal(reqBody)
	url := c.BaseURL + "/chat/completions"

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		trace.Status = "error"
		trace.Error = err.Error()
		c.logTrace(trace)
		return "", fmt.Errorf("AI connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		trace.Status = "error"
		trace.Error = string(body)
		c.logTrace(trace)
		return "", fmt.Errorf("AI returned error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		trace.Status = "error"
		trace.Error = err.Error()
		c.logTrace(trace)
		return "", fmt.Errorf("failed to decode AI response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		trace.Status = "error"
		trace.Error = "no choices returned"
		c.logTrace(trace)
		return "", fmt.Errorf("AI returned no choices")
	}

	res := chatResp.Choices[0].Message.Content
	trace.Response = res
	c.logTrace(trace)

	return res, nil
}

// OpenAI Embedding structures
type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// Embed generates embeddings for the given text.
func (c *AgentClient) Embed(text string) ([]float64, error) {
	reqBody := EmbedRequest{
		Model: c.Model,
		Input: text,
	}

	jsonData, _ := json.Marshal(reqBody)
	url := c.BaseURL + "/embeddings"

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("AI embed failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI embed error %d: %s", resp.StatusCode, string(body))
	}

	var embedResp EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("failed to decode AI embed response: %v", err)
	}

	if len(embedResp.Data) == 0 {
		return nil, fmt.Errorf("AI returned no embeddings")
	}

	return embedResp.Data[0].Embedding, nil
}

// DescribeImage uses the vision model to describe an image.
func (c *AgentClient) DescribeImage(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %v", err)
	}

	b64 := base64.StdEncoding.EncodeToString(data)
	imageURL := fmt.Sprintf("data:image/png;base64,%s", b64)

	type ImageURL struct {
		URL string `json:"url"`
	}

	type ImageContent struct {
		Type     string    `json:"type"`
		Text     string    `json:"text,omitempty"`
		ImageURL *ImageURL `json:"image_url,omitempty"`
	}

	type VisionMessage struct {
		Role    string         `json:"role"`
		Content []ImageContent `json:"content"`
	}

	type VisionRequest struct {
		Model    string          `json:"model"`
		Messages []VisionMessage `json:"messages"`
		Stream   bool            `json:"stream"`
	}

	reqBody := VisionRequest{
		Model: c.VisionModel,
		Messages: []VisionMessage{
			{
				Role: "user",
				Content: []ImageContent{
					{Type: "text", Text: "Describe this image in one concise sentence for a knowledge base index."},
					{Type: "image_url", ImageURL: &ImageURL{URL: imageURL}},
				},
			},
		},
		Stream: false,
	}

	jsonData, _ := json.Marshal(reqBody)
	url := c.BaseURL + "/chat/completions"

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("AI vision failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AI vision error %d: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode AI vision response: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("AI vision returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}
