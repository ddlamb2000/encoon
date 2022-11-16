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
						<tr><td>Text1</td><td>{this.props.row.text1}</td></tr>
						<tr><td>Text2</td><td>{this.props.row.text2}</td></tr>
						<tr><td>Text3</td><td>{this.props.row.text3}</td></tr>
						<tr><td>Text4</td><td>{this.props.row.text4}</td></tr>
						<tr><td>Int1</td><td>{this.props.row.int1}</td></tr>
						<tr><td>Int2</td><td>{this.props.row.int2}</td></tr>
						<tr><td>Int3</td><td>{this.props.row.int3}</td></tr>
						<tr><td>Int4</td><td>{this.props.row.int4}</td></tr>
						<tr><td>Version</td><td>{this.props.row.version}</td></tr>
						<tr><td>Created by</td><td>{this.props.row.createdBy}</td></tr>
						<tr><td>Updated by</td><td>{this.props.row.updatedBy}</td></tr>
					</tbody>
				</table>
			</div>
		)
	}
}
