// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridRow extends React.Component {
	render() {
		const variant = (this.props.rowSelected ? "table-warning" : "")
		return (
			<tr className={variant} onClick={() => this.props.onRowClick(this.props.row)}>
				<td>
					{this.props.row.added && <span>*</span>}
					{!this.props.row.added && <span>&nbsp;</span>}
				</td>
				{this.props.row.editable && <td><input></input></td>}
				{!this.props.row.editable && <td>{this.props.row.uri}</td>}
				<td>{this.props.row.text01}</td>
				<td>{this.props.row.text02}</td>
				<td>{this.props.row.text03}</td>
				<td>{this.props.row.text04}</td>
				<td scope="row"><a href={this.props.row.path}>{this.props.row.uuid}</a></td>
			</tr>
		)
	}
}
