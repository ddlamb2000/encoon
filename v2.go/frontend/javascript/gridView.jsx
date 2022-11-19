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
							col => <tr key={col.col}>
										<td>{col.label}</td>
										<td><GridCell uuid={this.props.row.uuid}
														key={col.col}
														col={col.col}
														type={col.type}
														typeUuid={col.typeUuid}
														value={col.value}
														values={col.values}
														readonly={col.readonly}
														rowAdded={this.props.rowAdded}
														rowSelected={this.props.rowSelected}
														rowEdited={this.props.rowEdited}
														onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
														onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
														inputRef={this.props.inputRef} /></td>
									</tr>
						)}
					</tbody>
				</table>
			</div>
		)
	}
}
