// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
  render() {
    const {
      row,
      rowsAdded,
      rowsSelected,
      rowsEdited,
      referencedValuesAdded,
      referencedValuesRemoved
    } = this.props;
    const columns = getColumnValuesForRow(this.props.columns, this.props.row, true);
    const columnsUsage = getColumnValuesForRow(this.props.columnsUsage, this.props.row, false);
    const audits = this.props.row ? this.props.row.audits : [];
    const activeGrid = this.props.grid.uuid == UuidGrids ? "show active" : "";
    const activeDefinition = this.props.grid.uuid != UuidGrids ? "show active" : "";
    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("h5", {
      className: "card-title"
    }, this.props.row.displayString), this.props.grid.uuid == UuidGrids && /*#__PURE__*/React.createElement("div", {
      className: "card-subtitle mb-2 text-muted"
    }, this.props.row.text2, " ", this.props.row.text3 && /*#__PURE__*/React.createElement("small", null, /*#__PURE__*/React.createElement("i", {
      className: `bi bi-${this.props.row.text3} mx-1`
    }))), this.props.grid.uuid != UuidGrids && /*#__PURE__*/React.createElement("div", {
      className: "card-subtitle mb-2 text-muted"
    }, this.props.grid.displayString, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: () => this.props.navigateToGrid(UuidGrids, this.props.grid.uuid)
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-box-arrow-up-right mx-1"
    }))), /*#__PURE__*/React.createElement("nav", {
      className: "nav nav-tabs",
      id: "tabs",
      role: "tablist"
    }, this.props.grid.uuid == UuidGrids && /*#__PURE__*/React.createElement("a", {
      className: "nav-link " + activeGrid,
      id: "rows-tab",
      "data-bs-toggle": "tab",
      "data-bs-target": "#rows-tab-pane",
      type: "button",
      role: "tab",
      "aria-controls": "rows-tab-pane",
      "aria-selected": "false"
    }, "Data rows"), /*#__PURE__*/React.createElement("a", {
      className: "nav-link " + activeDefinition,
      id: "definition-tab",
      "data-bs-toggle": "tab",
      "data-bs-target": "#definition-tab-pane",
      type: "button",
      role: "tab",
      "aria-controls": "definition-tab-pane",
      "aria-selected": "true"
    }, "Definition"), columnsUsage && /*#__PURE__*/React.createElement("a", {
      className: "nav-link",
      id: "usages-tab",
      "data-bs-toggle": "tab",
      "data-bs-target": "#usages-tab-pane",
      type: "button",
      role: "tab",
      "aria-controls": "usages-tab-pane",
      "aria-selected": "false"
    }, "Referenced by")), /*#__PURE__*/React.createElement("div", {
      className: "tab-content"
    }, /*#__PURE__*/React.createElement("div", {
      className: "tab-pane fade " + activeDefinition,
      id: "definition-tab-pane",
      role: "tabpanel",
      "aria-labelledby": "definition-tab",
      tabIndex: "0"
    }, /*#__PURE__*/React.createElement("table", {
      className: "table table-hover table-sm"
    }, /*#__PURE__*/React.createElement("thead", {
      className: "table-light"
    }), /*#__PURE__*/React.createElement("tbody", null, columns && columns.map(column => /*#__PURE__*/React.createElement("tr", {
      key: column.uuid
    }, /*#__PURE__*/React.createElement("td", null, column.label, !column.owned && /*#__PURE__*/React.createElement("small", null, "\xA0[", column.grid.displayString, "]"), trace && /*#__PURE__*/React.createElement("small", null, "\xA0", /*#__PURE__*/React.createElement("em", null, column.name))), /*#__PURE__*/React.createElement(GridCell, {
      uuid: row.uuid,
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
      bidirectional: column.bidirectional,
      canEditRow: row.canEditRow,
      grid: this.props.grid,
      displayString: row.displayString,
      rowSelected: rowsSelected.includes(row.uuid),
      rowEdited: rowsEdited.includes(row.uuid),
      rowAdded: rowsAdded.includes(row.uuid),
      referencedValuesAdded: referencedValuesAdded.filter(ref => ref.columnUuid == column.uuid),
      referencedValuesRemoved: referencedValuesRemoved.filter(ref => ref.columnUuid == column.uuid),
      onSelectRowClick: uuid => this.props.onSelectRowClick(uuid),
      onEditRowClick: uuid => this.props.onEditRowClick(uuid),
      onAddReferencedValueClick: reference => this.props.onAddReferencedValueClick(reference),
      onRemoveReferencedValueClick: reference => this.props.onRemoveReferencedValueClick(reference),
      inputRef: this.props.inputRef,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid),
      dbName: this.props.dbName,
      token: this.props.token,
      loadParentData: () => this.props.loadParentData()
    }))), audits && /*#__PURE__*/React.createElement("tr", null, /*#__PURE__*/React.createElement("td", null, "Audit"), /*#__PURE__*/React.createElement("td", null, /*#__PURE__*/React.createElement("ul", {
      className: "list-unstyled mb-0"
    }, audits.map(audit => {
      return /*#__PURE__*/React.createElement("li", {
        key: audit.uuid
      }, audit.actionName, " on ", /*#__PURE__*/React.createElement(DateTime, {
        dateTime: audit.created
      }), ", by ", /*#__PURE__*/React.createElement("a", {
        href: "#",
        onClick: () => this.props.navigateToGrid(UuidUsers, audit.createdBy)
      }, audit.createdByName));
    }))))))), this.props.grid.uuid == UuidGrids && /*#__PURE__*/React.createElement("div", {
      className: "tab-pane fade " + activeGrid,
      id: "rows-tab-pane",
      role: "tabpanel",
      "aria-labelledby": "rows-tab",
      tabIndex: "1"
    }, /*#__PURE__*/React.createElement(Grid, {
      token: this.props.token,
      dbName: this.props.dbName,
      gridUuid: this.props.row.uuid,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid),
      innerGrid: true
    })), columnsUsage && /*#__PURE__*/React.createElement("div", {
      className: "tab-pane fade",
      id: "usages-tab-pane",
      role: "tabpanel",
      "aria-labelledby": "usages-tab",
      tabIndex: "2"
    }, columnsUsage && columnsUsage.map(column => /*#__PURE__*/React.createElement(Grid, {
      token: this.props.token,
      dbName: this.props.dbName,
      key: column.uuid,
      gridUuid: column.grid.uuid,
      filterColumnOwned: "true",
      filterColumnName: column.name,
      filterColumnLabel: column.label,
      filterColumnGridUuid: column.gridUuid,
      filterColumnValue: this.props.row.uuid,
      filterColumnDisplayString: this.props.row.displayString,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    })))));
  }
}