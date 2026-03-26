package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"` // Trigger explicit JSON mode
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func AskOllama(prompt string, expectJSON bool) (string, error) {
	apiURL := os.Getenv("OLLAMA_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:11434/api/generate"
	}

	model := os.Getenv("OLLAMA_MODEL_TEXT")
	if expectJSON {
		model = os.Getenv("OLLAMA_MODEL_JSON")
	}

	reqPayload := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	if expectJSON {
		reqPayload.Format = "json"
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("gagal merakit payload JSON Ollama: %w", err)
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("gagal menghubungi service Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyErr, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API Error HTTP %d: %s", resp.StatusCode, string(bodyErr))
	}

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", fmt.Errorf("gagal decode response Ollama: %w", err)
	}

	return ollamaResp.Response, nil
}
