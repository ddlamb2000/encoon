// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"d.lambert.fr/encoon/configuration"
	"d.lambert.fr/encoon/database"
	"d.lambert.fr/encoon/model"
	_ "github.com/lib/pq"
)

func TestGetRowsWhereQueryForGridsApi(t *testing.T) {
	var tests = []struct {
		testName, uuid, expect string
	}{
		{"withUuid", "1234", " WHERE uuid = $2 AND griduuid = $1 "},
		{"withoutUuid", "", " WHERE griduuid = $1 "},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := getRowsWhereQueryForGridsApi(tt.uuid)
			if got != tt.expect {
				t.Errorf(`Got %v instead of %v.`, got, tt.expect)
			}
		})
	}
}

func TestGetRowsQueryParametersForGridsApi(t *testing.T) {
	var tests = []struct {
		testName, griduuid, uuid string
		expect                   []any
	}{
		{"withUuid", "1", "1234", []any{"1", "1234"}},
		{"withoutUuid", "2", "", []any{"2"}},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got := getRowsQueryParametersForGridsApi(tt.griduuid, tt.uuid)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf(`Got %v instead of %v.`, got, tt.expect)
			}
		})
	}
}

func TestGetGridForGridsApi(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	db, _ := database.GetDbByName(dbName)
	gridUri := "_users"
	user := "root"
	grid, err := getGridForGridsApi(context.Background(), db, dbName, user, gridUri, "")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
	if grid.Uuid != model.UuidUsers {
		t.Errorf(`Grid Uuid is wrong: %v.`, grid.Uuid)
	}
}

func TestGetRowsForGridsApi(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	user := "root"
	db, _ := database.GetDbByName(dbName)
	_, err := getRowsForGridsApi(context.Background(), db, dbName, user, model.UuidUsers, "", "")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
}

func TestGetRowsForGridsApi2(t *testing.T) {
	configuration.LoadConfiguration("../testData/validConfiguration1.yml")
	dbName := "test"
	user := "root"
	db, _ := database.GetDbByName(dbName)
	_, err := getRowsForGridsApi(context.Background(), db, dbName, user, "xxx", "", "")
	if err == nil {
		t.Errorf(`expect error.`)
	}
	if !strings.Contains(fmt.Sprintf("%v", err), "Error when querying rows:") {
		t.Errorf(`Wrong error: %v.`, err)
	}
}

func TestGetRowsQueryOutputForGridsApi(t *testing.T) {
	var row model.Row
	got := getRowsQueryOutputForGridsApi(&row)
	expect := []any{
		&row.Uuid,
		&row.Text1,
		&row.Text2,
		&row.Text3,
		&row.Text4,
		&row.Int1,
		&row.Int2,
		&row.Int3,
		&row.Int4,
		&row.Enabled,
		&row.CreateBy,
		&row.UpdatedBy,
		&row.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetGridQueryColumnsForGridsApi(t *testing.T) {
	got := getGridQueryColumnsForGridsApi()
	expect := "SELECT uuid, text1, text2, text3, text4, enabled, createdBy, updatedBy, version"
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
		&grid.CreateBy,
		&grid.UpdatedBy,
		&grid.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
