// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

package apis

import (
	"context"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
	"d.lambert.fr/encoon/utils"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func postGridEmbeddings(ct context.Context, p ApiParameters, grid *model.Grid, reponse GridResponse) error {
	r, cancel, err := createContextAndApiRequest(ct, p)
	defer cancel()
	if err != nil {
		return r.logAndReturnError("Can't connect to Gemini: %v", err)
	}
	apiKey := readFileContent(p.DbName, r.p.UserName, configuration.GetConfiguration().AI.ApiKeyFile)
	embeddingModel := configuration.GetConfiguration().AI.EmbeddingModel
	if apiKey != "" && embeddingModel != "" {
		r.log("Access Gemini")
		clientAI, err := genai.NewClient(ct, option.WithAPIKey(apiKey))
		if err != nil {
			return r.logAndReturnError("Can't connect to Gemini: %v", err)
		}
		defer clientAI.Close()
		reponse.generateEmbedding(r, grid, clientAI)
	} else {
		r.log("Gemini not configured")
	}
	return nil
}

func (response *GridResponse) generateEmbedding(r ApiRequest, grid *model.Grid, clientAI *genai.Client) error {
	for _, row := range response.Rows {
		if row.Revision > row.RevisionEmbedding {
			row.EmbeddingString = "{UUIDREFERENCE:" + r.p.DbName + "/" + row.GridUuid + "/" + row.Uuid + "}\n"
			row.EmbeddingString = row.EmbeddingString + response.Grid.DisplayString + "\n"
			for _, col := range response.Grid.Columns {
				if *col.TypeUuid == model.UuidTextColumnType || *col.TypeUuid == model.UuidIntColumnType {
					colAsString := row.GetValueAsString(*col.Name)
					if colAsString != "" {
						row.EmbeddingString = row.EmbeddingString + *col.Label + ": " + colAsString + "\n"
					}
				}
			}
			for _, ref := range row.References {
				row.EmbeddingString = row.EmbeddingString + ref.Label + ": "
				for indexRefRow, refRow := range ref.Rows {
					if indexRefRow > 0 {
						row.EmbeddingString = row.EmbeddingString + " ; "
					}
					row.EmbeddingString = row.EmbeddingString + refRow.DisplayString
				}
				row.EmbeddingString = row.EmbeddingString + "\n"
			}
			if err := generateEmbeddingWithModel(r, grid, row, clientAI); err != nil {
				return err
			}
		}
	}
	return nil
}

func generateEmbeddingWithModel(r ApiRequest, grid *model.Grid, row model.Row, clientAI *genai.Client) error {
	row.RevisionEmbedding = row.Revision
	em := clientAI.EmbeddingModel(configuration.GetConfiguration().AI.EmbeddingModel)
	em.TaskType = genai.TaskTypeRetrievalDocument
	res, err := em.EmbedContent(r.ctx, genai.Text(row.EmbeddingString))
	if err != nil {
		return r.logAndReturnError("Issue when generating embedding with Gemini for grid %s and row %s with embedding %s: %v", grid.Uuid, row.Uuid, row.EmbeddingString, err)
	}
	row.Embedding = res.Embedding.Values
	row.TokenCount = int64(len(res.Embedding.Values))
	r.log("Embedding generated for grid %s row %s: %d tokens", row.GridUuid, row.Uuid, row.TokenCount)
	return persistEmbedding(r, grid, row)
}

func persistEmbedding(r ApiRequest, grid *model.Grid, row model.Row) error {
	r.trace("persistEmbedding()")
	if err := r.beginTransaction(); err != nil {
		return err
	}
	if err := persistEmbeddingRow(r, grid, row); err != nil {
		_ = r.rollbackTransaction()
		return err
	}
	if err := r.commitTransaction(); err != nil {
		_ = r.rollbackTransaction()
		return err
	}
	r.trace("persistEmbedding() - Done")
	return nil
}

func persistEmbeddingRow(r ApiRequest, grid *model.Grid, row model.Row) error {
	query := getUpdateStatementForEmbeddingRow(grid)
	parms := getUpdateValuesForEmbeddingRow(row)
	r.trace("persistEmbeddingRow(%s, %s) - query=%s ; parms=%s", row.GridUuid, row.Uuid, query, parms)
	if err := r.execContext(query, parms...); err != nil {
		return r.logAndReturnError("Update row error: %v.", err)
	}
	r.log("Row [%s] updated with embedding.", row.Uuid)
	return nil
}

// function is available for mocking
var getUpdateStatementForEmbeddingRow = func(grid *model.Grid) string {
	return "UPDATE " +
		grid.GetTableName() +
		" SET embedding = $3, " +
		"revisionEmbedding = $4, " +
		"tokenCount = $5, " +
		"embeddingText = $6" +
		" WHERE gridUuid = $1 " +
		"AND uuid = $2"
}

func getUpdateValuesForEmbeddingRow(row model.Row) []any {
	values := make([]any, 0)
	values = append(values, row.GridUuid)
	values = append(values, row.Uuid)
	values = append(values, utils.VectorToString(row.Embedding))
	values = append(values, row.RevisionEmbedding)
	values = append(values, row.TokenCount)
	values = append(values, row.EmbeddingString)
	return values
}
