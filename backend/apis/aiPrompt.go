// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"
	"fmt"
	"os"
	"strings"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
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
	request.DbName = dbName
	request.UserUuid = userUuid
	request.UserName = userName
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
		return answerPromptWithClient(ctx, dbName, userUuid, userName, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey, request, clientAI)
	} else {
		return responseContent{
			Status:      FailedStatus,
			Action:      request.Action,
			ActionText:  request.ActionText,
			TextMessage: "Gemini not configured",
		}
	}
}

func answerPromptWithClient(ctx context.Context, dbName, userUuid, userName, user, gridUuid, contextUuid, requestInitiatedOn, requestReceivedOn, messageKey string, request ApiParameters, clientAI *genai.Client) responseContent {
	configuration.Log(dbName, userName, "Request Gemini")
	embeddingModel := configuration.GetConfiguration().AI.EmbeddingModel
	if embeddingModel != "" {
		em := clientAI.EmbeddingModel(configuration.GetConfiguration().AI.EmbeddingModel)
		em.TaskType = genai.TaskTypeRetrievalQuery
		res, err := em.EmbedContent(ctx, genai.Text(request.ActionText))
		if err != nil {
			return responseContent{
				Status:      FailedStatus,
				Action:      request.Action,
				ActionText:  request.ActionText,
				TextMessage: fmt.Sprint("Issue when generating embedding with Gemini: %v", err),
			}
		}
		rowSet, rowSetCount, err := getEmbeddingMatchingPrompt(ctx, request, res.Embedding.Values)
		if err != nil {
			return responseContent{
				Status:      FailedStatus,
				Action:      request.Action,
				ActionText:  request.ActionText,
				TextMessage: fmt.Sprint("Issue when querying vectors in database: %v", err),
			}
		}
		model := clientAI.GenerativeModel(configuration.GetConfiguration().AI.Model)
		configureModel(model)
		prompt := getPromptWithContext(rowSetCount, rowSet, request.ActionText)
		configuration.Log(dbName, userName, "Prompt submitted to Gemini: %s", prompt)
		iter := model.GenerateContentStream(ctx, genai.Text(prompt))
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

func getPromptWithContext(rowSetCount int, rowSet []model.Row, question string) string {
	if rowSetCount == 0 {
		return question
	} else {
		prompt := "Answer the following question based on the provided context:\n"
		prompt = prompt + "\n\n"
		prompt = prompt + "Question: " + question + "\n\n"
		prompt = prompt + "Context:\n"
		for _, row := range rowSet {
			prompt = prompt + row.EmbeddingString + "\n"
		}
		return prompt
	}
}

func configureModel(model *genai.GenerativeModel) {
	model.SetTemperature(configuration.GetConfiguration().AI.Temperature)
	model.SetTopP(configuration.GetConfiguration().AI.TopP)
	model.SetTopK(configuration.GetConfiguration().AI.TopK)
	model.SetMaxOutputTokens(configuration.GetConfiguration().AI.MaxOutputTokens)
	model.ResponseMIMEType = "text/plain"
	model.SystemInstruction = genai.NewUserContent(genai.Text(configuration.GetConfiguration().AI.SystemInstruction) + " and show URI_REFERENCE")
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

func getEmbeddingMatchingPrompt(ct context.Context, p ApiParameters, vector []float32) ([]model.Row, int, error) {
	r, cancel, err := createContextAndApiRequest(ct, p)
	defer cancel()
	t := r.startTiming()
	defer r.stopTiming("getEmbeddingMatchingPrompt()", t)
	if err != nil {
		return nil, 0, err
	}
	r.trace("getEmbeddingMatchingPrompt()")
	rowSet, rowSetCount, err := getRowSetForEmbeddingMatchingPrompt(r, vector)
	if err != nil {
		return nil, 0, err
	}
	r.trace("getGridsRows() - Done")
	return rowSet, rowSetCount, nil
}

func getRowSetForEmbeddingMatchingPrompt(r ApiRequest, vector []float32) ([]model.Row, int, error) {
	r.trace("getRowSetForEmbeddingMatchingPrompt()")
	t := r.startTiming()
	defer r.stopTiming("getRowSetForEmbeddingMatchingPrompt()", t)
	query := getRowsQueryForEmbeddingMatchingPrompt()
	parms := getRowsQueryParametersForEmbeddingMatchingPrompt(vector)
	r.trace("getRowSetForEmbeddingMatchingPrompt() - query=%s ; parms=%s", query, parms)
	set, err := r.db.QueryContext(r.ctx, query, parms...)
	if err != nil {
		return nil, 0, r.logAndReturnError("Error when querying rows: %v.", err)
	}
	defer set.Close()
	rows := make([]model.Row, 0)
	for set.Next() {
		row := model.GetNewRow()
		if err := set.Scan(getRowsQueryOutputForEmbeddingMatchingPrompt(row)...); err != nil {
			return nil, 0, r.logAndReturnError("Error when scanning rows: %v.", err)
		}
		r.trace("getRowSetForEmbeddingMatchingPrompt(%s, %s, %v) - row=%v", row)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, *row)
	}
	return rows, len(rows), nil
}

var getRowsQueryForEmbeddingMatchingPrompt = func() string {
	return "SELECT rows.embeddingText " +
		"FROM rows " +
		"WHERE enabled = true " +
		"AND embedding <-> $1 < $2"
}

// function is available for mocking
var getRowsQueryParametersForEmbeddingMatchingPrompt = func(vector []float32) []any {
	parameters := make([]any, 0)
	parameters = append(parameters, utils.VectorToString(vector))
	parameters = append(parameters, configuration.GetConfiguration().AI.VectorDistance)
	return parameters
}

var getRowsQueryOutputForEmbeddingMatchingPrompt = func(row *model.Row) []any {
	output := make([]any, 0)
	output = append(output, &row.EmbeddingString)
	return output
}
