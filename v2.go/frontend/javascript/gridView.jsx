// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const columns = getColumnValuesForRow(this.props.columns, this.props.row, true)
		return (
			<div>
				<h4 className="card-title">{this.props.row.displayString}</h4>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						{columns && columns.map(
							column => <tr key={column.name}>
											<td>{column.label}</td>
											<GridCell uuid={this.props.row.uuid}
														columnName={column.name}
														type={column.type}
														typeUuid={column.typeUuid}
														value={column.value}
														values={column.values}
														readonly={column.readonly}
														rowAdded={this.props.rowAdded}
														rowSelected={this.props.rowSelected}
														rowEdited={this.props.rowEdited}
														referencedValuesAdded={this.props.referencedValuesAdded}
														referencedValuesRemoved={this.props.referencedValuesRemoved}
														onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
														onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
														inputRef={this.props.inputRef} />
									</tr>
						)}
						<tr>
							<td>Created</td>
							<td>{this.props.row.created} by {this.props.row.createdBy}</td>
						</tr>
						<tr>
							<td>Updated</td>
							<td>{this.props.row.updated} by {this.props.row.updatedBy}</td>
						</tr>
					</tbody>
				</table>
			</div>
		)
	}
}
