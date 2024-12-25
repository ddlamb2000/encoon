// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class Navigation extends React.Component {
	render() {
		if(trace) console.log("[Navigation.render()]")
		return (
			<nav className="position-sticky pt-4 sidebar-sticky">
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="My grids"
						gridSubTitle="Grids I own"
						filterColumnOwned={true}
						filterColumnName='relationship3'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid} />
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="Edit grids"
						gridSubTitle="Grids I can edit"
						filterColumnOwned={true}
						filterColumnName='relationship5'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid}
						noEdit={true} />
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="View grids"
						gridSubTitle="Grids I can view"
						filterColumnOwned={true}
						filterColumnName='relationship4'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid}
						noEdit={true} />
			</nav>
		)
	}
}
