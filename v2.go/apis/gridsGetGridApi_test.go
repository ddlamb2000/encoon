// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"reflect"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
)

func TestGetGridForGridsApi(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, _ := database.GetDbByName(dbName)
	gridUri := "_users"
	user := "root"
	grid, err := getGridForGridsApi(context.Background(), db, dbName, user, gridUri, "true")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
		return
	}
	if grid.Uuid != model.UuidUsers {
		t.Errorf(`Grid Uuid is wrong: %v.`, grid.Uuid)
	}
}

func TestGetGridQueryColumnsForGridsApi(t *testing.T) {
	got := getGridQueryForGridsApi()
	expect := "SELECT uuid, text1, text2, text3, text4, enabled, created, createdBy, updated, updatedBy, version FROM rows WHERE gridUuid = $1 AND text1 = $2"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetGridQueryOutputForGridsApi(t *testing.T) {
	var grid model.Grid
	got := getGridQueryOutputForGridsApi(&grid)
	expect := []any{
		&grid.Uuid,
		&grid.Text1,
		&grid.Text2,
		&grid.Text3,
		&grid.Text4,
		&grid.Enabled,
		&grid.Created,
		&grid.CreatedBy,
		&grid.Updated,
		&grid.UpdatedBy,
		&grid.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
