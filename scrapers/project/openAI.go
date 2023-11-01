package project

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// A general function used to make API request and get the response with the given `prompt`
func openAIRequest(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	apiURL := "https://api.openai.com/v1/chat/completions"

	messages := make([]any, 0)
	messages = append(messages, map[string]string{
		"role":    "user",
		"content": prompt,
	})

	data := map[string]any{
		"model":    "gpt-3.5-turbo",
		"messages": messages,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshalling the request data failed:\n\t%v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("preparing API request to OpenAI failed:\n\t%v", err)
	}

	// Add the headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("making API request to OpenAI failed:\n\t%v", err)
	}
	defer resp.Body.Close()

	// Read the response
	var result map[string]any
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("unmarshalling the response data failed:\n\t%v", err)
	}

	if completion, ok := result["choices"].([]any); ok && len(completion) > 0 {
		if message, ok := completion[0].(map[string]any)["message"]; ok {
			if content, ok := message.(map[string]any)["content"].(string); ok {
				return content, nil
			}
		}
	}

	return "", fmt.Errorf("failed to get a valid response")
}

// Used to get the description from OpenAI for a partiucular organization
func getAIDescription(orgName string) (string, error) {
	prompt := fmt.Sprintf("Give me a short description of about 80-100 words about the popular LFX organization %v which introduces the same to new contributors. Answer in a single paragraph and return nothing else that the textual description.", orgName)

	desc, err := openAIRequest(prompt)
	if err != nil {
		return "", fmt.Errorf("openAI request failed for org %v:\n\t%v", orgName, err)
	}

	return desc, err
}
