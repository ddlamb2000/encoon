// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridRowHeader extends React.Component {
	render() {
		const { column, filterColumnOwned, filterColumnName, filterColumnGridUuid } = this.props
		return (
			<th scope="col">
				{this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && <mark>{column.label}</mark>}
				{!this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && column.label}
				{!column.owned && <small><br />[{column.grid.displayString}]</small>}
				{trace && <small><br /><em>{column.gridUuid}</em></small>}
				{trace && <small><br /><em>{column.name}</em></small>}
			</th>
		)
	}

	matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) {
		const ownership = (column.owned && filterColumnOwned == 'true') || (!column.owned && filterColumnOwned != 'true')
		return ownership && column.name == filterColumnName && column.gridUuid == filterColumnGridUuid
	}
}
