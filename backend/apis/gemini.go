// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"fmt"
	"os"
	"strings"

	"d.lambert.fr/encoon/configuration"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const modelName = "gemini-2.0-flash-exp"

var client *genai.Client = nil

func readFileContent(filePath string) string {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		configuration.LogError("", "", "Can't read file %s: %v", filePath, err)
		return ""
	}
	return string(dat)
}

func getResponse(resp *genai.GenerateContentResponse) string {
	var response strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				response.WriteString(fmt.Sprintf("%v", part))
			}
		}
	}
	return response.String()
}

func answerPrompt(request ApiParameters) responseContent {
	apiKey := readFileContent(configuration.GetConfiguration().AI.ApiKeyFile)
	if apiKey != "" {
		ctx := context.Background()
		if client == nil {
			var err error = nil
			client, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
			if err != nil {
				configuration.LogError("", "", "Can't connect to Gemini: %v", err)
				return responseContent{
					Status:      FailedStatus,
					Action:      request.Action,
					ActionText:  request.ActionText,
					TextMessage: err.Error(),
				}
			}
		}
		configuration.Log("", "", "Access Gemini for generating content")
		model := client.GenerativeModel(modelName)
		resp, err := model.GenerateContent(ctx, genai.Text(request.ActionText))
		if err != nil {
			configuration.LogError("", "", "Gemini can't generated content: %v", err)
			return responseContent{
				Status:      FailedStatus,
				Action:      request.Action,
				ActionText:  request.ActionText,
				TextMessage: err.Error(),
			}
		}
		configuration.Log("", "", "Gemini content generated")
		return responseContent{
			Status:      SuccessStatus,
			Action:      request.Action,
			ActionText:  request.ActionText,
			TextMessage: getResponse(resp),
		}
	} else {
		return responseContent{
			Status:      FailedStatus,
			Action:      request.Action,
			ActionText:  request.ActionText,
			TextMessage: "Gemini not configured",
		}
	}
}
