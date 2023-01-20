// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

package apis

import (
	"reflect"
	"testing"

	"d.lambert.fr/encoon/model"
	_ "github.com/lib/pq"
)

func TestGetRowsQueryOutputForGridsApi(t *testing.T) {
	row := model.GetNewRow()
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int8", Owned: true})
	got := getRowsQueryOutputForGridsApi(grid, row)
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
		&row.Revision,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetRowsQueryOutputForGridsApi2(t *testing.T) {
	row := model.GetNewRow()
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "text6", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text7", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text8", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text9", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "text10", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int1", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int2", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int3", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int4", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int5", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int6", Owned: true})
	got := getRowsQueryOutputForGridsApi(grid, row)
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
		&row.Revision,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}

func TestGetRowsQueryOutputForGridsApi3(t *testing.T) {
	row := model.GetNewRow()
	grid := model.GetNewGrid()
	grid.Uuid = model.UuidUsers
	grid.Columns = append(grid.Columns, &model.Column{Name: "int7", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int9", Owned: true})
	grid.Columns = append(grid.Columns, &model.Column{Name: "int10", Owned: true})
	got := getRowsQueryOutputForGridsApi(grid, row)
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
		&row.Revision,
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf(`Got %v instead of %v.`, got, expect)
	}
}
