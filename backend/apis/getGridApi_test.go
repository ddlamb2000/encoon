// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"context"
	"reflect"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/model"
)

func TestGetGridForGridsApi(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	p := ApiParameters{
		DbName:   "test",
		GridUuid: model.UuidUsers,
		UserName: "root",
	}
	r, cancel, _ := createContextAndApiRequest(context.Background(), p)
	defer cancel()
	grid, err := getGridForGridsApi(r, p.GridUuid)
	if err != nil {
		t.Errorf(`Error: %v.`, err)
		return
	}
	if grid.Uuid != model.UuidUsers {
		t.Errorf(`Grid Uuid is wrong: %v.`, grid.Uuid)
	}
}

func TestGetGridQueryOutputForGridsApi(t *testing.T) {
	grid := model.GetNewGrid("")
	got := getGridQueryOutputForGridsApi(grid)
	expect := []any{
		&grid.Uuid,
		&grid.GridUuid,
		&grid.Text1,
		&grid.Text2,
		&grid.Text3,
		&grid.Enabled,
		&grid.Created,
		&grid.CreatedBy,
		&grid.Updated,
		&grid.UpdatedBy,
		&grid.Revision,
		&grid.OwnerUuid,
		&grid.DefaultAccessUuid,
		&grid.ViewAccessUuid,
		&grid.EditAccessUuid,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
