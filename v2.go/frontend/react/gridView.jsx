// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const columns = getColumnValuesForRow(this.props.columns, this.props.row, true)
		const columnsUsage = getColumnValuesForRow(this.props.columnsUsage, this.props.row, false)
		const audits = this.props.row ? this.props.row.audits : []
		return (
			<div>
				<h4 className="card-title">{this.props.row.displayString}</h4>
				<div className="card-subtitle mb-2 text-muted">
					{this.props.grid.displayString}
					<a href="#" onClick={() => this.props.navigateToGrid(UuidGrids, this.props.grid.uuid)}>
						<i className="bi bi-box-arrow-up-right mx-1"></i>
					</a>
				</div>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						{columns && columns.map(column => 
							<tr key={column.uuid}>
								<td>
									{column.label}
									{!column.owned && <small>&nbsp;[{column.grid.displayString}]</small>}
									{trace && <small>&nbsp;<em>{column.name}</em></small>}
								</td>
								<GridCell uuid={this.props.row.uuid}
											columnUuid={column.uuid}
											owned={column.owned}
											columnName={column.name}
											columnLabel={column.label}
											type={column.type}
											typeUuid={column.typeUuid}
											value={column.value}
											values={column.values}
											gridPromptUuid={column.gridPromptUuid}
											readonly={column.readonly}
											bidirectional={column.bidirectional}
											grid={this.props.grid}
											displayString={this.props.row.displayString}
											rowAdded={this.props.rowAdded}
											rowSelected={this.props.rowSelected}
											rowEdited={this.props.rowEdited}
											referencedValuesAdded={this.props.referencedValuesAdded}
											referencedValuesRemoved={this.props.referencedValuesRemoved}
											onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
											onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
											inputRef={this.props.inputRef}
											navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
											dbName={this.props.dbName}
											token={this.props.token} />
							</tr>
						)}
						{audits &&
							<tr>
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
				{this.props.grid.uuid == UuidGrids &&
					<div>
						<Grid token={this.props.token}
									dbName={this.props.dbName} 
									gridUuid={this.props.row.uuid} 
									navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
					</div>
				}
				{columnsUsage && 
					<div>
						{columnsUsage.map(column =>
							<Grid token={this.props.token}
									dbName={this.props.dbName} 
									key={column.uuid}
									gridUuid={column.grid.uuid} 
									filterColumnOwned='true'
									filterColumnName={column.name}
									filterColumnLabel={column.label}
									filterColumnGridUuid={column.gridUuid}
									filterColumnValue={this.props.row.uuid}
									filterColumnDisplayString={this.props.row.displayString}
									navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
						)}
					</div>
				}
			</div>
		)
	}
}
