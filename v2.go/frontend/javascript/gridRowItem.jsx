// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowAdded ?
							"table-warning" :
							(this.props.rowSelected ? "table-secondary" : "")
		return (
			<tr className={variant} onClick={() => this.props.onRowClick(this.props.row)}>
				<td>
					{this.props.rowAdded && <span>*</span>}
					{!this.props.rowAdded && <span>&nbsp;</span>}
				</td>
				{this.props.rowEdited && <td><input></input></td>}
				{!this.props.rowEdited && <td>{this.props.row.uri}</td>}
				<td>{this.props.row.text01}</td>
				<td>{this.props.row.text02}</td>
				<td>{this.props.row.text03}</td>
				<td>{this.props.row.text04}</td>
				<td scope="row"><a href={this.props.row.path}>{this.props.row.uuid}</a></td>
			</tr>
		)
	}
}
