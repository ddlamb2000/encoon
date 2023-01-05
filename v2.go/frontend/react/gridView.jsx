// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const { row, rowsAdded, rowsSelected, rowsEdited, referencedValuesAdded, referencedValuesRemoved } = this.props
		const columns = getColumnValuesForRow(this.props.columns, this.props.row, true)
		const columnsUsage = getColumnValuesForRow(this.props.columnsUsage, this.props.row, false)
		const audits = this.props.row ? this.props.row.audits : []
		const activeGrid = this.props.grid.uuid == UuidGrids ? "show active" : ""
		const activeDefinition = this.props.grid.uuid != UuidGrids ? "show active" : ""
		return (
			<div>
				<h5 className="card-title">{this.props.row.displayString}</h5>
				{this.props.grid.uuid == UuidGrids &&
					<div className="card-subtitle mb-2 text-muted">
						{this.props.row.text2} {this.props.row.text3 && <small><i className={`bi bi-${this.props.row.text3} mx-1`}></i></small>}
					</div>
				}
				{this.props.grid.uuid != UuidGrids &&
					<div className="card-subtitle mb-2 text-muted">
						{this.props.grid.displayString}
						<a href="#" onClick={() => this.props.navigateToGrid(UuidGrids, this.props.grid.uuid)}>
							<i className="bi bi-box-arrow-up-right mx-1"></i>
						</a>
					</div>
				}
				<nav className="nav nav-tabs" id="tabs" role="tablist">
					{this.props.grid.uuid == UuidGrids &&
						<a className={"nav-link " + activeGrid}
								id="rows-tab"
								data-bs-toggle="tab"
								data-bs-target="#rows-tab-pane"
								type="button"
								role="tab"
								aria-controls="rows-tab-pane"
								aria-selected="false">
							Data rows
						</a>
					}
					<a className={"nav-link " + activeDefinition}
							id="definition-tab"
							data-bs-toggle="tab"
							data-bs-target="#definition-tab-pane"
							type="button"
							role="tab"
							aria-controls="definition-tab-pane"
							aria-selected="true">
						Definition
					</a>
					{columnsUsage &&
						<a className="nav-link" 
								id="usages-tab"
								data-bs-toggle="tab"
								data-bs-target="#usages-tab-pane"
								type="button"
								role="tab"
								aria-controls="usages-tab-pane"
								aria-selected="false">
							Referenced by
						</a>
					}
				</nav>
				<div className="tab-content">
					<div className={"tab-pane fade " + activeDefinition}
							id="definition-tab-pane"
							role="tabpanel"
							aria-labelledby="definition-tab"
							tabIndex="0">
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
										<GridCell uuid={row.uuid}
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
													canEditRow={row.canEditRow}
													grid={this.props.grid}
													displayString={row.displayString}
													rowSelected={rowsSelected.includes(row.uuid)}
													rowEdited={rowsEdited.includes(row.uuid)}
													rowAdded={rowsAdded.includes(row.uuid)}
													referencedValuesAdded={referencedValuesAdded.filter(ref => ref.columnUuid == column.uuid)}
													referencedValuesRemoved={referencedValuesRemoved.filter(ref => ref.columnUuid == column.uuid)}
													onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
													onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
													onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
													onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)}
													inputRef={this.props.inputRef}
													navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
													dbName={this.props.dbName}
													token={this.props.token}
													loadParentData={() => this.props.loadParentData()} />
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
					</div>
					{this.props.grid.uuid == UuidGrids &&
						<div className={"tab-pane fade " + activeGrid}
								id="rows-tab-pane"
								role="tabpanel"
								aria-labelledby="rows-tab"
								tabIndex="1">
							<Grid token={this.props.token}
										dbName={this.props.dbName} 
										gridUuid={this.props.row.uuid} 
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
										innerGrid={true} />
						</div>
					}
					{columnsUsage &&
						<div className="tab-pane fade"
								id="usages-tab-pane"
								role="tabpanel"
								aria-labelledby="usages-tab"
								tabIndex="2">
							{columnsUsage && columnsUsage.map(column =>
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
			</div>
		)
	}
}
