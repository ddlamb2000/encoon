// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowAdded || this.props.rowSelected ? "table-warning" : ""
		return (
			<tr className={variant} onClick={() => this.props.onRowClick(this.props.row)}>
				<td>
					{this.props.rowAdded && <span>*</span>}
					{!this.props.rowAdded && <span>&nbsp;</span>}
				</td>
				<GridCell uuid={this.props.row.uuid} col="uri" value={this.props.row.uri} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text01" value={this.props.row.text01} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text02" value={this.props.row.text02} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text03" value={this.props.row.text03} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<GridCell uuid={this.props.row.uuid} col="text04" value={this.props.row.text04} rowEdited={this.props.rowEdited} inputRef={this.props.inputRef} />
				<td scope="row"><a href={this.props.row.path}>{this.props.row.uuid}</a></td>
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
				{this.props.rowEdited &&
					<input id={this.id}
							uuid={this.props.uuid}
							type="text"
							col={this.props.col}
							defaultValue={this.props.value}
							ref={this.props.inputRef} />
				}
				{!this.props.rowEdited && <span>{this.props.value}</span>}
			</td>
		)
	}
}