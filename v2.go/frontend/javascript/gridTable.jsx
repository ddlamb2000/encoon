// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive">
				<thead>
					<tr>
						<th scope="col" style={{width: "24px"}}></th>
						<th scope="col">Text01</th>
						<th scope="col">Text02</th>
						<th scope="col">Text03</th>
						<th scope="col">Text04</th>
						<th scope="col" style={{width: "64px"}}>Int01</th>
						<th scope="col" style={{width: "64px"}}>Int02</th>
						<th scope="col" style={{width: "64px"}}>Int03</th>
						<th scope="col" style={{width: "64px"}}>Int04</th>
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
		columns.push({col: "text01", value: this.props.row.text01, type: "text", readonly: false})
		columns.push({col: "text02", value: this.props.row.text02, type: "text", readonly: false})
		columns.push({col: "text03", value: this.props.row.text03, type: "text", readonly: false})
		columns.push({col: "text04", value: this.props.row.text04, type: "text", readonly: false})
		columns.push({col: "int01", value: this.props.row.int01, type: "number", readonly: false})
		columns.push({col: "int02", value: this.props.row.int02, type: "number", readonly: false})
		columns.push({col: "int03", value: this.props.row.int03, type: "number", readonly: false})
		columns.push({col: "int04", value: this.props.row.int04, type: "number", readonly: false})
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
