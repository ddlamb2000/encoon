// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
  render() {
    return /*#__PURE__*/React.createElement("table", {
      className: "table table-hover table-sm table-responsive align-middle"
    }, /*#__PURE__*/React.createElement("thead", null, /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("th", {
      scope: "col",
      style: {
        width: "24px"
      }
    }), this.props.columns && this.props.columns.map(column => /*#__PURE__*/React.createElement(GridRowHeader, {
      key: column.uuid,
      column: column,
      filterColumnOwned: this.props.filterColumnOwned,
      filterColumnName: this.props.filterColumnName,
      filterColumnGridUuid: this.props.filterColumnGridUuid,
      grid: this.props.grid
    })), /*#__PURE__*/React.createElement("th", {
      className: "text-end",
      scope: "col"
    }, "Revision"))), /*#__PURE__*/React.createElement("tbody", {
      className: "table-group-divider"
    }, this.props.rows.map(row => /*#__PURE__*/React.createElement(GridRow, {
      key: row.uuid,
      row: row,
      rowSelected: this.props.rowsSelected.includes(row.uuid),
      rowEdited: this.props.rowsEdited.includes(row.uuid),
      rowAdded: this.props.rowsAdded.includes(row.uuid),
      referencedValuesAdded: this.props.referencedValuesAdded.filter(ref => ref.fromUuid == row.uuid),
      referencedValuesRemoved: this.props.referencedValuesRemoved.filter(ref => ref.fromUuid == row.uuid),
      columns: this.props.columns,
      onSelectRowClick: uuid => this.props.onSelectRowClick(uuid),
      onEditRowClick: uuid => this.props.onEditRowClick(uuid),
      onDeleteRowClick: uuid => this.props.onDeleteRowClick(uuid),
      onAddReferencedValueClick: reference => this.props.onAddReferencedValueClick(reference),
      onRemoveReferencedValueClick: reference => this.props.onRemoveReferencedValueClick(reference),
      inputRef: this.props.inputRef,
      dbName: this.props.dbName,
      token: this.props.token,
      grid: this.props.grid,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    }))));
  }
}
class GridRowHeader extends React.Component {
  render() {
    const {
      column,
      filterColumnOwned,
      filterColumnName,
      filterColumnGridUuid
    } = this.props;
    return /*#__PURE__*/React.createElement("th", {
      scope: "col"
    }, this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && /*#__PURE__*/React.createElement("mark", null, column.label), !this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && column.label, !column.owned && /*#__PURE__*/React.createElement("small", null, /*#__PURE__*/React.createElement("br", null), "[", column.grid.displayString, "]"), trace && /*#__PURE__*/React.createElement("small", null, /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("em", null, column.gridUuid)), trace && /*#__PURE__*/React.createElement("small", null, /*#__PURE__*/React.createElement("br", null), /*#__PURE__*/React.createElement("em", null, column.name)));
  }
  matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) {
    const ownership = column.owned && filterColumnOwned == 'true' || !column.owned && filterColumnOwned != 'true';
    return ownership && column.name == filterColumnName && column.gridUuid == filterColumnGridUuid;
  }
}
class GridRow extends React.Component {
  render() {
    const {
      row,
      rowAdded,
      rowSelected,
      rowEdited,
      referencedValuesAdded,
      referencedValuesRemoved
    } = this.props;
    const variantEdited = this.props.rowEdited ? "table-warning" : "";
    const columns = getColumnValuesForRow(this.props.columns, row);
    return /*#__PURE__*/React.createElement("tr", {
      className: variantEdited
    }, /*#__PURE__*/React.createElement("td", {
      scope: "row",
      className: "vw-10"
    }, !(rowAdded || rowSelected) && /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: () => this.props.navigateToGrid(row.gridUuid, row.uuid)
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-card-text"
    })), row.canEditRow && (rowAdded || rowSelected) && /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "btn text-danger btn-sm mx-0 p-0",
      onClick: () => this.props.onDeleteRowClick(row.uuid)
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-dash-circle"
    }))), columns.map(column => /*#__PURE__*/React.createElement(GridCell, {
      uuid: row.uuid,
      key: column.uuid,
      columnUuid: column.uuid,
      owned: column.owned,
      columnName: column.name,
      columnLabel: column.label,
      type: column.type,
      typeUuid: column.typeUuid,
      value: column.value,
      values: column.values,
      gridPromptUuid: column.gridPromptUuid,
      readonly: column.readonly,
      bidirectional: false,
      canEditRow: row.canEditRow,
      rowAdded: rowAdded,
      rowSelected: rowSelected,
      rowEdited: rowEdited,
      referencedValuesAdded: referencedValuesAdded.filter(ref => ref.columnUuid == column.uuid),
      referencedValuesRemoved: referencedValuesRemoved.filter(ref => ref.columnUuid == column.uuid),
      onSelectRowClick: uuid => this.props.onSelectRowClick(uuid),
      onEditRowClick: uuid => this.props.onEditRowClick(uuid),
      onAddReferencedValueClick: reference => this.props.onAddReferencedValueClick(reference),
      onRemoveReferencedValueClick: reference => this.props.onRemoveReferencedValueClick(reference),
      inputRef: this.props.inputRef,
      dbName: this.props.dbName,
      token: this.props.token,
      grid: this.props.grid,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    })), /*#__PURE__*/React.createElement("td", {
      className: "text-end"
    }, this.props.row.revision));
  }
}
class GridCell extends React.Component {
  render() {
    const variantReadOnly = this.props.readonly ? "form-control-plaintext" : "";
    const checkedBoolean = this.props.value && this.props.value == "true" ? true : false;
    const variantMonospace = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : "";
    const embedded = this.props.bidirectional && this.props.owned;
    return /*#__PURE__*/React.createElement("td", {
      onClick: () => this.props.onSelectRowClick(this.props.canEditRow ? this.props.uuid : '')
    }, (this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" && !embedded && /*#__PURE__*/React.createElement(GridCellInput, {
      type: this.props.type,
      variantReadOnly: variantReadOnly,
      uuid: this.props.uuid,
      columnName: this.props.columnName,
      readOnly: this.props.readonly,
      checkedBoolean: checkedBoolean,
      value: this.props.value,
      inputRef: this.props.inputRef,
      onEditRowClick: uuid => this.props.onEditRowClick(uuid)
    }), !(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" && !embedded && /*#__PURE__*/React.createElement("span", {
      className: variantMonospace
    }, getCellValue(this.props.type, this.props.value)), (this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" && !embedded && /*#__PURE__*/React.createElement(GridCellDropDown, {
      uuid: this.props.uuid,
      columnUuid: this.props.columnUuid,
      columnName: this.props.columnName,
      owned: this.props.owned,
      values: this.props.values,
      dbName: this.props.dbName,
      token: this.props.token,
      gridPromptUuid: this.props.gridPromptUuid,
      referencedValuesAdded: this.props.referencedValuesAdded,
      referencedValuesRemoved: this.props.referencedValuesRemoved,
      onAddReferencedValueClick: reference => this.props.onAddReferencedValueClick(reference),
      onRemoveReferencedValueClick: reference => this.props.onRemoveReferencedValueClick(reference)
    }), !(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" && !embedded && /*#__PURE__*/React.createElement(GridCellReferences, {
      uuid: this.props.uuid,
      values: this.props.values,
      referencedValuesAdded: this.props.referencedValuesAdded,
      referencedValuesRemoved: this.props.referencedValuesRemoved,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    }), this.props.type == "reference" && embedded && /*#__PURE__*/React.createElement(Grid, {
      token: this.props.token,
      dbName: this.props.dbName,
      gridUuid: this.props.gridPromptUuid,
      filterColumnOwned: this.props.owned ? 'false' : 'true',
      filterColumnName: this.props.columnName,
      filterColumnLabel: this.props.columnLabel,
      filterColumnGridUuid: this.props.grid.uuid,
      filterColumnValue: this.props.uuid,
      filterColumnDisplayString: this.props.displayString,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    }));
  }
}
class GridCellInput extends React.Component {
  render() {
    const className = this.props.type == "checkbox" ? "form-check-input" : "form-control";
    return /*#__PURE__*/React.createElement("input", {
      type: this.props.type,
      className: className + " form-control-sm rounded-2 shadow gap-2 p-1 " + this.props.variantReadOnly,
      name: this.props.uuid,
      uuid: this.props.uuid,
      column: this.props.columnName,
      readOnly: this.props.readonly,
      defaultChecked: this.props.checkedBoolean,
      defaultValue: this.props.value,
      ref: this.props.inputRef,
      onInput: () => this.props.onEditRowClick(this.props.uuid)
    });
  }
}
class GridCellReferences extends React.Component {
  render() {
    const {
      values,
      referencedValuesAdded,
      referencedValuesRemoved
    } = this.props;
    const referencedValuesIncluded = values.concat(referencedValuesAdded).filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).filter((value, index, self) => index == self.findIndex(t => t.uuid == value.uuid));
    if (referencedValuesIncluded.length > 0) {
      return /*#__PURE__*/React.createElement("ul", {
        className: "list-unstyled mb-0"
      }, values.map(value => /*#__PURE__*/React.createElement("li", {
        key: value.uuid
      }, /*#__PURE__*/React.createElement("a", {
        className: "gap-2 p-0",
        href: "#",
        onClick: () => this.props.navigateToGrid(value.gridUuid, value.uuid)
      }, value.displayString))));
    }
  }
}
class GridCellDropDown extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: "",
      isLoaded: false,
      isLoading: false,
      rows: []
    };
  }
  render() {
    const {
      isLoading,
      error,
      rows
    } = this.state;
    const {
      gridPromptUuid,
      values,
      referencedValuesAdded,
      referencedValuesRemoved
    } = this.props;
    const referencedValuesIncluded = error ? [] : values.concat(referencedValuesAdded).filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).filter((value, index, self) => index == self.findIndex(t => t.uuid == value.uuid));
    const referencedValuesNotIncluded = error ? [] : rows.concat(referencedValuesRemoved).filter(ref => !referencedValuesAdded.map(ref => ref.uuid).includes(ref.uuid)).filter(ref => !referencedValuesIncluded.map(ref => ref.uuid).includes(ref.uuid)).filter((value, index, self) => index == self.findIndex(t => t.uuid == value.uuid));
    const countRows = referencedValuesNotIncluded ? referencedValuesNotIncluded.length : 0;
    if (trace) {
      console.log("[GridCellDropDown.render()] this.props.columnName=", this.props.columnName);
      console.log("[GridCellDropDown.render()] values=", values);
      console.log("[GridCellDropDown.render()] referencedValuesAdded=", referencedValuesAdded);
      console.log("[GridCellDropDown.render()] referencedValuesRemoved=", referencedValuesRemoved);
      console.log("[GridCellDropDown.render()] referencedValuesIncluded=", referencedValuesIncluded);
      console.log("[GridCellDropDown.render()] referencedValuesNotIncluded=", referencedValuesNotIncluded);
    }
    return /*#__PURE__*/React.createElement("ul", {
      className: "list-unstyled mb-0"
    }, referencedValuesIncluded.map(ref => /*#__PURE__*/React.createElement("li", {
      key: ref.uuid
    }, /*#__PURE__*/React.createElement("span", null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "btn text-danger btn-sm mx-0 p-0",
      onClick: () => this.props.onRemoveReferencedValueClick({
        fromUuid: this.props.uuid,
        columnUuid: this.props.columnUuid,
        owned: this.props.owned,
        columnName: this.props.columnName,
        toGridUuid: this.props.gridPromptUuid,
        uuid: ref.uuid,
        displayString: ref.displayString
      })
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-box-arrow-down pe-1"
    })), ref.displayString))), gridPromptUuid && /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement("input", {
      type: "search",
      className: "form-control form-control-sm rounded-2 shadow gap-2 p-1",
      autoComplete: "false",
      placeholder: "Search...",
      onInput: e => this.loadDropDownData(gridPromptUuid, e.target.value)
    })), isLoading && /*#__PURE__*/React.createElement("li", null, /*#__PURE__*/React.createElement(Spinner, null)), error && !isLoading && /*#__PURE__*/React.createElement("li", {
      className: "alert alert-danger",
      role: "alert"
    }, error), referencedValuesNotIncluded && countRows > 0 && referencedValuesNotIncluded.map(ref => /*#__PURE__*/React.createElement("li", {
      key: ref.uuid
    }, /*#__PURE__*/React.createElement("span", null, /*#__PURE__*/React.createElement("button", {
      type: "button",
      className: "btn text-success btn-sm mx-0 p-0",
      onClick: () => this.props.onAddReferencedValueClick({
        fromUuid: this.props.uuid,
        columnUuid: this.props.columnUuid,
        owned: this.props.owned,
        columnName: this.props.columnName,
        toGridUuid: this.props.gridPromptUuid,
        uuid: ref.uuid,
        displayString: ref.displayString
      })
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-box-arrow-up pe-1"
    })), ref.displayString))));
  }
  loadDropDownData(gridPromptUuid, value) {
    this.setState({
      isLoading: true
    });
    if (value.length > 0) {
      const uri = `/${this.props.dbName}/api/v1/${gridPromptUuid}`;
      fetch(uri, {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + this.props.token
        }
      }).then(response => {
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.indexOf("application/json") !== -1) {
          return response.json().then(result => {
            if (result.response != undefined) {
              this.setState({
                isLoading: false,
                isLoaded: true,
                rows: result.response.rows,
                error: result.response.error
              });
            } else {
              this.setState({
                isLoading: false,
                isLoaded: true,
                error: result.error
              });
            }
          }, error => {
            this.setState({
              isLoading: false,
              isLoaded: false,
              rows: [],
              error: error.message
            });
          });
        } else {
          this.setState({
            isLoading: false,
            isLoaded: false,
            rows: [],
            error: `[${response.status}] Internal server issue.`
          });
        }
      });
    } else {
      this.setState({
        isLoading: false,
        isLoaded: false,
        rows: [],
        error: ""
      });
    }
  }
}