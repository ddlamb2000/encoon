// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		const columns = getColumnValuesForRow(this.props.columns, this.props.row, true)
		const createdByUri = `/${this.props.dbName}/${UuidUsers}/${this.props.row.createdBy}`
		const updatedByUri = `/${this.props.dbName}/${UuidUsers}/${this.props.row.updatedBy}`
		const audits = this.props.row ? this.props.row.audits : []
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
							<td>Audit</td>
							<td>
								<ul className="list-unstyled mb-0">
									{audits && audits.map(
										audit => {
											const uri = `/${this.props.dbName}/${UuidUsers}/${audit.createdBy}`
											return (
												<li key={audit.uuid}>
													{audit.actionName}
													&nbsp;on <DateTime dateTime={audit.created} />
													&nbsp;by <a href={uri}>{audit.createdByName}</a>
												</li>
											)
										}
									)}
								</ul>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		)
	}
}
