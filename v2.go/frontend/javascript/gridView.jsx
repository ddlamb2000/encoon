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
							col => <tr key={col.name}>
										<td>{col.label}</td>
										<GridCell uuid={this.props.row.uuid}
													col={col.name}
													type={col.type}
													typeUuid={col.typeUuid}
													value={col.value}
													values={col.values}
													readonly={col.readonly}
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
					</tbody>
				</table>
			</div>
		)
	}
}
