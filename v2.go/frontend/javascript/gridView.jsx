// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const columns = getColumnValuesForRow(this.props.columns, this.props.row, true)
		const audits = this.props.row ? this.props.row.audits : []
		return (
			<div>
				<h4 className="card-title">{this.props.row.displayString}</h4>
				<div className="card-subtitle mb-2 text-muted">{this.props.grid.displayString}</div>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						{columns && columns.map(
							column => <tr key={column.uuid}>
											<td>
												{column.label}
												{!column.owned && <small>&nbsp;[{column.grid.displayString}]</small>}
												<small>&nbsp;<em>{column.name}</em></small>
											</td>
											<GridCell uuid={this.props.row.uuid}
														columnUuid={column.uuid}
														owned={column.owned}
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
														inputRef={this.props.inputRef}
														navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
									</tr>
						)}
						{audits && <tr>
								<td>Audit</td>
								<td>
									<ul className="list-unstyled mb-0">
										{audits.map(
											audit => {
												return (
													<li key={audit.uuid}>
														{audit.actionName} on <DateTime dateTime={audit.created} />,
														by <a href="#" onClick={() => this.props.navigateToGrid(UuidUsers, audit.createdBy)}>
															{audit.createdByName}
														</a>
													</li>
												)
											}
										)}
									</ul>
								</td>
							</tr>
						}
					</tbody>
				</table>
			</div>
		)
	}
}
