// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowAdded ? "table-warning" : ""
		return (
			<tr className={variant} onClick={() => this.props.onRowClick(this.props.row)}>
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
							className="btn btn-sm mx-0 p-0"
							onClick={() => this.deleteRows()}>
							<i className="bi bi-dash-circle"></i>
						</button>
					}
				</td>
				<GridCell uuid={this.props.row.uuid} col="uri" value={this.props.row.uri} rowSelected={this.props.rowSelected} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text01" value={this.props.row.text01} rowSelected={this.props.rowSelected} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text02" value={this.props.row.text02} rowSelected={this.props.rowSelected} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text03" value={this.props.row.text03} rowSelected={this.props.rowSelected} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text04" value={this.props.row.text04} rowSelected={this.props.rowSelected} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
			</tr>
		)
	}
}

class GridCell extends React.Component {
	constructor(props) {
		super(props)
		this.id = this.props.uuid + "-" + this.props.col
	}

	render() {
		return (
			<td>
				{(this.props.rowEdited || this.props.rowSelected) && 
					<input type="text"
							className="form-control form-control-sm"
							id={this.id}
							uuid={this.props.uuid}
							col={this.props.col}
							defaultValue={this.props.value}
							ref={this.props.inputRef}
							onInput={() => this.setChanged()} />
				}
				{!(this.props.rowEdited || this.props.rowSelected) && <span>{this.props.value}</span>}
			</td>
		)
	}

	setChanged() {
		console.log("changed.")
	}
}