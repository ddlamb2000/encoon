// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const columns = getColumnValueForRow(this.props.columns, this.props.row, true)
		return (
			<div>
				<h6 className="card-subtitle mb-2 text-muted">{this.props.row.uuid}</h6>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						{columns && columns.map(
							col => <tr key={col.col}><td>{col.label}</td><td>{getCellValue(col.type, col.value)}</td></tr>
						)}
					</tbody>
				</table>
			</div>
		)
	}
}
