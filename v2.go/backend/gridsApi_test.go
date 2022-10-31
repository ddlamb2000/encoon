// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package backend

import (
	"reflect"
	"testing"

	"d.lambert.fr/encoon/utils"
	_ "github.com/lib/pq"
)

func TestGetRowsWhereQueryForGridsApi(t *testing.T) {
	var tests = []struct {
		testName, uuid, expect string
	}{
		{"withUuid", "1234", "WHERE uuid = $2 AND griduuid = $1"},
		{"withoutUuid", "", "WHERE griduuid = $1"},
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
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	dbName := "test"
	db := getDbByName(dbName)
	gridUri := "users"
	grid, err := getGridForGridsApi(db, gridUri)
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
	if grid.Uuid != utils.UuidUsers {
		t.Errorf(`Grid Uuid is wrong: %v.`, grid.Uuid)
	}
}

func TestGetRowsForGridsApi(t *testing.T) {
	utils.LoadConfiguration("../configurations/")
	ConnectDbServers(utils.DatabaseConfigurations)
	dbName := "test"
	db := getDbByName(dbName)
	_, err := getRowsForGridsApi(db, utils.UuidUsers, "")
	if err != nil {
		t.Errorf(`Error: %v.`, err)
	}
}

func TestGetRowsQueryOutputForGridsApi(t *testing.T) {
	var row Row
	got := getRowsQueryOutputForGridsApi(&row)
	expect := []any{
		&row.Uuid,
		&row.Version,
		&row.Uri,
		&row.Text01,
		&row.Text02,
		&row.Text03,
		&row.Text04,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
