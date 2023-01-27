// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridFooter extends React.Component {
	render() {
		const { grid, rows, uuid, canAddRows, rowsEdited, rowsAdded, rowsDeleted, isLoading, miniGrid } = this.props
		const countRows = rows ? rows.length : 0
		const countRowsAdded = rowsAdded ? rowsAdded.length : 0
		const countRowsEdited = rowsEdited ? rowsEdited.length : 0
		const countRowsDeleted = rowsDeleted ? rowsDeleted.length : 0
		return (
			<nav className='mb-3' onClick={() => this.props.onSelectRowClick()}>
				{isLoading && <Spinner />}
				{!isLoading && !miniGrid && countRows == 0 && <small className="text-muted px-1">No data</small>}
				{!isLoading && !miniGrid && countRows == 1 && uuid == '' && <small className="text-muted px-1">{countRows} row</small>}
				{!isLoading && !miniGrid && countRows > 1 && <small className="text-muted px-1">{countRows} rows</small>}
				{!isLoading && !miniGrid && countRowsAdded > 0 && <small className="text-muted px-1">({countRowsAdded} added)</small>}
				{!isLoading && !miniGrid && countRowsEdited > 0 && <small className="text-muted px-1">({countRowsEdited} edited)</small>}
				{!isLoading && !miniGrid && countRowsDeleted > 0 && <small className="text-muted px-1">({countRowsDeleted} deleted)</small>}
				{!isLoading && grid && uuid == "" && canAddRows &&
					<button type="button"
							className="btn btn-outline-success btn-sm mx-1"
							onClick={this.props.onAddRowClick}>
						Add <i className="bi bi-plus-circle"></i>
					</button>
				}
				{!isLoading && countRowsAdded + countRowsEdited + countRowsDeleted > 0 &&
					<button type="button"
							className="btn btn-outline-warning btn-sm mx-1"
							onClick={this.props.onLoadDataClick}>
						Cancel <i className="bi bi-arrow-counterclockwise"></i>
					</button>
				}
				{!isLoading && countRowsAdded + countRowsEdited + countRowsDeleted > 0 &&
					<button type="button"
							className="btn btn-outline-primary btn-sm mx-1"
							onClick={this.props.onSaveDataClick}>
						Save <i className="bi bi-save"></i>
					</button>
				}
			</nav>
		)
	}
}
