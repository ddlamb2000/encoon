// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

package apis

import (
	"reflect"
	"testing"

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

func TestGetRowsQueryOutputForGridsApi(t *testing.T) {
	row := model.Row{}
	grid := model.Grid{}
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8"})
	got := getRowsQueryOutputForGridsApi(&grid, &row)
	expect := []any{
		&row.Uuid,
		&row.GridUuid,
		&row.Text2,
		&row.Text5,
		&row.Int8,
		&row.Enabled,
		&row.Created,
		&row.CreatedBy,
		&row.Updated,
		&row.UpdatedBy,
		&row.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetRowsQueryOutputForGridsApi2(t *testing.T) {
	row := model.Row{}
	grid := model.Grid{}
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text6"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text7"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text8"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text9"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text10"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int1"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int2"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int3"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int4"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int5"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int6"})
	got := getRowsQueryOutputForGridsApi(&grid, &row)
	expect := []any{
		&row.Uuid,
		&row.GridUuid,
		&row.Text6,
		&row.Text7,
		&row.Text8,
		&row.Text9,
		&row.Text10,
		&row.Int1,
		&row.Int2,
		&row.Int3,
		&row.Int4,
		&row.Int5,
		&row.Int6,
		&row.Enabled,
		&row.Created,
		&row.CreatedBy,
		&row.Updated,
		&row.UpdatedBy,
		&row.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetRowsQueryOutputForGridsApi3(t *testing.T) {
	row := model.Row{}
	grid := model.Grid{}
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "int7"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int9"})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int10"})
	got := getRowsQueryOutputForGridsApi(&grid, &row)
	expect := []any{
		&row.Uuid,
		&row.GridUuid,
		&row.Int7,
		&row.Int9,
		&row.Int10,
		&row.Enabled,
		&row.Created,
		&row.CreatedBy,
		&row.Updated,
		&row.UpdatedBy,
		&row.Version,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
