// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowAdded || this.props.rowSelected ? "table-warning" : ""
		return (
			<tr className={variant} onClick={() => this.props.onRowClick(this.props.row)}>
				<td scope="row">
					{this.props.rowAdded && <span>*</span>}
					{!this.props.rowAdded && 
					<button
						type="button"
						className="btn btn-sm mx-0 p-0">
						<a href={this.props.row.path}><i className="bi bi-card-text"></i></a>
					</button>
					}
					<button
						type="button"
						className="btn btn-sm mx-0 p-0"
						onClick={() => this.addRow()}>
						<i className="bi bi-plus-circle"></i>
					</button>
					<button
						type="button"
						className="btn btn-sm mx-0 p-0"
						onClick={() => this.deleteRows()}>
						<i className="bi bi-dash-circle"></i>
					</button>
				</td>
				<GridCell uuid={this.props.row.uuid} col="uri" value={this.props.row.uri} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text01" value={this.props.row.text01} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text02" value={this.props.row.text02} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text03" value={this.props.row.text03} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text04" value={this.props.row.text04} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
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
			<td className="pt-0">
					<input type="text"
							className="form-control form-control-sm form-control-plaintext px-1"
							id={this.id}
							uuid={this.props.uuid}
							col={this.props.col}
							defaultValue={this.props.value}
							ref={this.props.inputRef} />
				{/* {!this.props.rowEdited && <span>{this.props.value}</span>} */}
			</td>
		)
	}
}