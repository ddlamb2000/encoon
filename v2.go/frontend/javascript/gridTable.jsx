// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive">
				<thead>
					<tr>
						<th scope="col" style={{width: "24px"}}></th>
						{this.props.columns && this.props.columns.map(
							col => <th scope="col" key={col.name}>{col.label}<br/><small>{col.name}</small></th>
						)}
					</tr>
				</thead>
				<tbody className="table-group-divider">
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.props.rowsSelected.includes(row.uuid)}
										rowEdited={this.props.rowsEdited.includes(row.uuid)}
										rowAdded={this.props.rowsAdded.includes(row.uuid)}
										columns={this.props.columns}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										onDeleteRowClick={uuid => this.props.onDeleteRowClick(uuid)}
										inputRef={this.props.inputRef} />
					)}
				</tbody>
			</table>
		)
	}
}

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowEdited ? "table-warning" : ""
		const columns = getColumnValueForRow(this.props.columns, this.props.row)
		return (
			<tr className={variant}>
				<td scope="row" className="vw-10">
					{!(this.props.rowAdded || this.props.rowSelected) && 
						<a href={this.props.row.path}><i className="bi bi-card-text"></i></a>
					}
					{(this.props.rowAdded || this.props.rowSelected) && 
						<button
							type="button"
							className="btn text-danger btn-sm mx-0 p-0"
							onClick={() => this.props.onDeleteRowClick(this.props.row.uuid)}>
							<i className="bi bi-dash-circle"></i>
						</button>
					}
				</td>
				{columns.map(
					col => <GridCell uuid={this.props.row.uuid}
										key={col.col}
										col={col}
										type={col.type}
										typeUuid={col.typeUuid}
										value={col.value}
										values={col.values}
										readonly={col.readonly}
										rowAdded={this.props.rowAdded}
										rowSelected={this.props.rowSelected}
										rowEdited={this.props.rowEdited}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										inputRef={this.props.inputRef} />
				)}
			</tr>
		)
	}
}

class GridCell extends React.Component {
	render() {
		const variant = this.props.readonly ? " form-control form-control-sm form-control-plaintext" : "form-control form-control-sm "
		const variantSize = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : ""
		return (
			<td onClick={() => this.props.onSelectRowClick(this.props.uuid)}>
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && 
					<input type={this.props.type}
							className={variant}
							uuid={this.props.uuid}
							col={this.props.col}
							readOnly={this.props.readonly}
							defaultValue={this.props.value}
							ref={this.props.inputRef}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" &&
					<span className={variantSize}>{getCellValue(this.props.type, this.props.value)}</span>
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" &&
					<GridCellReferences values={this.props.values}/>
				}
			</td>
		)
	}
}

class GridCellReferences extends React.Component {
	render() {
		return (
			<span>
				{this.props.values.map(value => <a className="pe-2" key={value.uuid} href={value.path}>{value.label}</a>)}
			</span>
		)
	}
}
