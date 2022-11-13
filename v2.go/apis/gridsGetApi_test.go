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
	"d.lambert.fr/encoon/utils"
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
	configuration.LoadConfiguration("../configuration.yml")
	dbName := "test"
	db, _ := database.GetDbByName(dbName)
	gridUri := "_users"
	user := "root"
	grid, err := getGridForGridsApi(context.Background(), db, dbName, user, gridUri, "")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
	if grid.Uuid != utils.UuidUsers {
		t.Errorf(`Grid Uuid is wrong: %v.`, grid.Uuid)
	}
}

func TestGetRowsForGridsApi(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
	dbName := "test"
	user := "root"
	db, _ := database.GetDbByName(dbName)
	_, err := getRowsForGridsApi(context.Background(), db, dbName, user, utils.UuidUsers, "", "")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
}

func TestGetRowsForGridsApi2(t *testing.T) {
	configuration.LoadConfiguration("../configuration.yml")
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
		&row.Text01,
		&row.Text02,
		&row.Text03,
		&row.Text04,
		&row.Int01,
		&row.Int02,
		&row.Int03,
		&row.Int04,
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
	expect := "SELECT uuid, text01, text02, text03, text04, enabled, createdBy, updatedBy, version"
	if got != expect {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetGridQueryOutputForGridsApi(t *testing.T) {
	var grid model.Grid
	got := getGridQueryOutputForGridsApi(&grid)
	expect := []any{
		&grid.Uuid,
		&grid.Text01,
		&grid.Text02,
		&grid.Text03,
		&grid.Text04,
		&grid.Enabled,
		&grid.CreateBy,
		&grid.UpdatedBy,
		&grid.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}