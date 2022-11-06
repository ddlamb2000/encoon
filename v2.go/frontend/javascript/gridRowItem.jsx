// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowEdited ? "table-warning" : ""
		const columns = []
		columns.push({col: "uri", value: this.props.row.uri})
		columns.push({col: "text01", value: this.props.row.text01})
		columns.push({col: "text02", value: this.props.row.text02})
		columns.push({col: "text03", value: this.props.row.text03})
		columns.push({col: "text04", value: this.props.row.text04})
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
				{(this.props.rowEdited || this.props.rowSelected) && 
					<input type="text"
							className="form-control form-control-sm"
							uuid={this.props.uuid}
							col={this.props.col}
							defaultValue={this.props.value}
							ref={this.props.inputRef}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				}
				{!(this.props.rowEdited || this.props.rowSelected) && <span>{this.props.value}</span>}
			</td>
		)
	}
}