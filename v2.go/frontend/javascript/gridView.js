// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
  render() {
    const columns = getColumnValuesForRow(this.props.columns, this.props.row, true);
    const columnsUsage = getColumnValuesForRow(this.props.columnsUsage, this.props.row, false);
    const audits = this.props.row ? this.props.row.audits : [];
    return /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("h4", {
      className: "card-title"
    }, this.props.row.displayString), /*#__PURE__*/React.createElement("div", {
      className: "card-subtitle mb-2 text-muted"
    }, this.props.grid.displayString, /*#__PURE__*/React.createElement("a", {
      href: "#",
      onClick: () => this.props.navigateToGrid(this.props.grid.uuid, "")
    }, /*#__PURE__*/React.createElement("i", {
      className: "bi bi-box-arrow-up-right mx-1"
    }))), /*#__PURE__*/React.createElement("table", {
      className: "table table-hover table-sm"
    }, /*#__PURE__*/React.createElement("thead", {
      className: "table-light"
    }), /*#__PURE__*/React.createElement("tbody", null, columns && columns.map(column => /*#__PURE__*/React.createElement("tr", {
      key: column.uuid
    }, /*#__PURE__*/React.createElement("td", null, column.label, !column.owned && /*#__PURE__*/React.createElement("small", null, "\xA0[", column.grid.displayString, "]"), /*#__PURE__*/React.createElement("small", null, "\xA0", /*#__PURE__*/React.createElement("em", null, column.name))), /*#__PURE__*/React.createElement(GridCell, {
      uuid: this.props.row.uuid,
      columnUuid: column.uuid,
      owned: column.owned,
      columnName: column.name,
      type: column.type,
      typeUuid: column.typeUuid,
      value: column.value,
      values: column.values,
      readonly: column.readonly,
      rowAdded: this.props.rowAdded,
      rowSelected: this.props.rowSelected,
      rowEdited: this.props.rowEdited,
      referencedValuesAdded: this.props.referencedValuesAdded,
      referencedValuesRemoved: this.props.referencedValuesRemoved,
      onSelectRowClick: uuid => this.props.onSelectRowClick(uuid),
      onEditRowClick: uuid => this.props.onEditRowClick(uuid),
      inputRef: this.props.inputRef,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
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
    })))))), this.props.grid.uuid == UuidGrids && /*#__PURE__*/React.createElement("div", null, /*#__PURE__*/React.createElement("h5", {
      className: "card-subtitle text-muted"
    }, "Data"), /*#__PURE__*/React.createElement(Grid, {
      token: this.props.token,
      dbName: this.props.dbName,
      gridUuid: this.props.row.uuid,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    })), columnsUsage && /*#__PURE__*/React.createElement("div", null, columnsUsage.length > 0 && /*#__PURE__*/React.createElement("h5", {
      className: "card-subtitle text-muted"
    }, "Usages"), columnsUsage.map(column => /*#__PURE__*/React.createElement(Grid, {
      token: this.props.token,
      dbName: this.props.dbName,
      key: column.uuid,
      gridUuid: column.grid.uuid,
      filterColumnName: column.name,
      filterColumnLabel: column.label,
      filterColumnValue: this.props.row.uuid,
      filterColumnDisplayString: this.props.row.displayString,
      navigateToGrid: (gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)
    }))));
  }
}