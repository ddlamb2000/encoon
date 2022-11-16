// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive">
				<thead>
					<tr>
						<th scope="col" style={{width: "24px"}}></th>
						<th scope="col">Text1</th>
						<th scope="col">Text2</th>
						<th scope="col">Text3</th>
						<th scope="col">Text4</th>
						<th scope="col" style={{width: "64px"}}>Int1</th>
						<th scope="col" style={{width: "64px"}}>Int2</th>
						<th scope="col" style={{width: "64px"}}>Int3</th>
						<th scope="col" style={{width: "64px"}}>Int4</th>
						<th scope="col" style={{width: "64px"}}>Version</th>
					</tr>
				</thead>
				<tbody className="table-group-divider">
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.props.rowsSelected.includes(row.uuid)}
										rowEdited={this.props.rowsEdited.includes(row.uuid)}
										rowAdded={this.props.rowsAdded.includes(row.uuid)}
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
		const columns = this.getColumns()
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
										col={col.col}
										type={col.type}
										value={col.value}
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

	getColumns() {
		const columns = []
		columns.push({col: "text1", value: this.props.row.text1, type: "text", readonly: false})
		columns.push({col: "text2", value: this.props.row.text2, type: "text", readonly: false})
		columns.push({col: "text3", value: this.props.row.text3, type: "text", readonly: false})
		columns.push({col: "text4", value: this.props.row.text4, type: "text", readonly: false})
		columns.push({col: "int1", value: this.props.row.int1, type: "number", readonly: false})
		columns.push({col: "int2", value: this.props.row.int2, type: "number", readonly: false})
		columns.push({col: "int3", value: this.props.row.int3, type: "number", readonly: false})
		columns.push({col: "int4", value: this.props.row.int4, type: "number", readonly: false})
		columns.push({col: "version", value: this.props.row.version, type: "number", readonly: true})
		return columns
	}
}

class GridCell extends React.Component {
	render() {
		const variant = this.props.readonly ? " form-control form-control-sm form-control-plaintext" : "form-control form-control-sm"
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
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && <span>{this.props.value}</span>}
			</td>
		)
	}
}
