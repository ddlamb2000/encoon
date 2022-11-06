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
						<th scope="col">Int01</th>
						<th scope="col">Int02</th>
						<th scope="col">Int03</th>
						<th scope="col">Int04</th>
						<th scope="col">Version</th>
					</tr>
				</thead>
				<tbody className="table-group-divider">
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.props.rowsSelected.includes(row.uuid)}
										rowEdited={this.props.rowsEdited.includes(row.uuid)}
										rowAdded={this.props.rowsAdded.includes(row.uuid)}
										onRowClick={uuid => this.props.onRowClick(uuid)}
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
		const columns = []
		columns.push({col: "text01", value: this.props.row.text01})
		columns.push({col: "text02", value: this.props.row.text02})
		columns.push({col: "text03", value: this.props.row.text03})
		columns.push({col: "text04", value: this.props.row.text04})
		columns.push({col: "int01", value: this.props.row.int01})
		columns.push({col: "int02", value: this.props.row.int02})
		columns.push({col: "int03", value: this.props.row.int03})
		columns.push({col: "int04", value: this.props.row.int04})
		columns.push({col: "versionint01", value: this.props.row.version})
		return (
			<tr className={variant}
				onClick={() => this.props.onRowClick(this.props.row.uuid)}>
				<td scope="row" className="vw-10">
					{!(this.props.rowAdded || this.props.rowSelected) && 
						<button
							type="button"
							className="btn btn-sm mx-0 p-0">
							<a href={this.props.row.path}><i className="bi bi-card-text"></i></a>
						</button>
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
										value={col.value}
										rowAdded={this.props.rowAdded}
										rowSelected={this.props.rowSelected}
										rowEdited={this.props.rowEdited}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										inputRef={this.props.inputRef} />
				)}
			</tr>
		)
	}
}

class GridCell extends React.Component {
	render() {
		return (
			<td>
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && 
					<input type="text"
							className="form-control form-control-sm"
							uuid={this.props.uuid}
							col={this.props.col}
							defaultValue={this.props.value}
							ref={this.props.inputRef}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && <span>{this.props.value}</span>}
			</td>
		)
	}
}
