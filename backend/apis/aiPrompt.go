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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func readFileContent(dbName string, userName string, filePath string) string {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		configuration.LogError(dbName, userName, "Can't read file %s: %v", filePath, err)
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

func answerPrompt(dbName, userUuid, userName, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey string, request ApiParameters) responseContent {
	apiKey := readFileContent(dbName, userName, configuration.GetConfiguration().AI.ApiKeyFile)
	if apiKey != "" {
		configuration.Log(dbName, userName, "Access Gemini")
		ctx := context.Background()
		clientAI, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			configuration.LogError(dbName, userName, "Can't connect to Gemini: %v", err)
			return responseContent{
				Status:      FailedStatus,
				Action:      request.Action,
				ActionText:  request.ActionText,
				TextMessage: err.Error(),
			}
		}
		defer clientAI.Close()
		configuration.Log(dbName, userName, "Request Gemini")
		model := clientAI.GenerativeModel(configuration.GetConfiguration().AI.Model)
		model.SetTemperature(configuration.GetConfiguration().AI.Temperature)
		model.SetTopP(configuration.GetConfiguration().AI.TopP)
		model.SetTopK(configuration.GetConfiguration().AI.TopK)
		model.SetMaxOutputTokens(configuration.GetConfiguration().AI.MaxOutputTokens)
		model.ResponseMIMEType = "text/plain"
		model.SystemInstruction = genai.NewUserContent(genai.Text(configuration.GetConfiguration().AI.SystemInstruction))
		iter := model.GenerateContentStream(ctx, genai.Text(request.ActionText))
		return handleGeminiResponse(dbName, userUuid, userName, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey, request, iter)
	} else {
		return responseContent{
			Status:      FailedStatus,
			Action:      request.Action,
			ActionText:  request.ActionText,
			TextMessage: "Gemini not configured",
		}
	}
}

func handleGeminiResponse(dbName, userUuid, userName, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey string, request ApiParameters, iter *genai.GenerateContentResponseIterator) responseContent {
	for responseNumber := 1; ; responseNumber++ {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			configuration.LogError(dbName, userName, "Error when retrieving content stream: %v", err)
			return responseContent{
				Status:      FailedStatus,
				Action:      request.Action,
				ActionText:  request.ActionText,
				TextMessage: err.Error(),
			}
		}
		response := responseContent{
			Status:         SuccessStatus,
			Action:         request.Action,
			ActionText:     request.ActionText,
			ResponseNumber: int64(responseNumber),
			TextMessage:    getResponse(resp),
		}
		WriteMessage(dbName, userUuid, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey, response)
	}
	configuration.Log(dbName, userName, "Gemini stream closed")
	return responseContent{
		Status:     SuccessStatus,
		Action:     request.Action,
		ActionText: request.ActionText,
	}
}
