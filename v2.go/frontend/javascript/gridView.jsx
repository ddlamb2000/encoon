// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridView extends React.Component {
	render() {
		return (
			<div>
				<h6 className="card-subtitle mb-2 text-muted">{this.props.row.uuid}</h6>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						<tr><td>Text01</td><td>{this.props.row.text01}</td></tr>
						<tr><td>Text02</td><td>{this.props.row.text02}</td></tr>
						<tr><td>Text03</td><td>{this.props.row.text03}</td></tr>
						<tr><td>Text04</td><td>{this.props.row.text04}</td></tr>
						<tr><td>Int01</td><td>{this.props.row.int01}</td></tr>
						<tr><td>Int02</td><td>{this.props.row.int02}</td></tr>
						<tr><td>Int03</td><td>{this.props.row.int03}</td></tr>
						<tr><td>Int04</td><td>{this.props.row.int04}</td></tr>
						<tr><td>Version</td><td>{this.props.row.version}</td></tr>
						<tr><td>Created by</td><td>{this.props.row.createdBy}</td></tr>
						<tr><td>Updated by</td><td>{this.props.row.updatedBy}</td></tr>
					</tbody>
				</table>
			</div>
		)
	}
}
