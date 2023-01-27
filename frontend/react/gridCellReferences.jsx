// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridCellReferences extends React.Component {
	render() {
		const { values, referencedValuesAdded, referencedValuesRemoved } = this.props
		const referencedValuesIncluded = values.
			concat(referencedValuesAdded).
			filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		if(referencedValuesIncluded.length > 0) {
			return (
				<ul className="list-unstyled mb-0">
					{values.map(value => 
						<li key={value.uuid}>
							<a className="gap-2 p-0" href="#" onClick={() => this.props.navigateToGrid(value.gridUuid, value.uuid)}>
								{value.displayString}
							</a>
						</li>
					)}
				</ul>
			)
		}
	}
}
