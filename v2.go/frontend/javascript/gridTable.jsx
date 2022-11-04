// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm">
				<thead className="table-light">
					<tr>
						<th scope="col"><span>&nbsp;</span></th>
						<th scope="col">Uri</th>
						<th scope="col">Text01</th>
						<th scope="col">Text02</th>
						<th scope="col">Text03</th>
						<th scope="col">Text04</th>
						<th scope="col">Uuid</th>
					</tr>
				</thead>
				<tbody>
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.isRowSelected(row)}
										rowEdited={this.isRowEdited(row)}
										rowAdded={this.isRowAdded(row)}
										onRowClick={row => this.props.onRowClick(row)}
										inputRef={this.props.inputRef} />
					)}
				</tbody>
			</table>
		)
	}

	isRowSelected(row) {
		return this.props.rowsSelected.includes(row.uuid)
	}

	isRowEdited(row) {
		return this.props.rowsEdited.includes(row.uuid)
	}

	isRowAdded(row) {
		return this.props.rowsAdded.includes(row.uuid)
	}
}