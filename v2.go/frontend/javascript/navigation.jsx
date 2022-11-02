// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Navigation extends React.Component {
	render() {
		return (
			<nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div className="position-sticky pt-3 sidebar-sticky">
					<ul className="nav flex-column">
						<li className="nav-item">
							<a className="nav-link active" aria-current="page" href={`/${this.props.dbName}`}>
								Dashboard <i className="bi bi-box"></i>
							</a>
						</li>
					</ul>
					<ul className="nav flex-column mb-2">
						<li className="nav-item">
							<a className="nav-link" href={`/${this.props.dbName}/grids`}>
								Grids <i className="bi bi-grid-3x3"></i>
							</a>
						</li>
						<li className="nav-item">
							<a className="nav-link" href={`/${this.props.dbName}/users`}>
								Users <i className="bi bi-person"></i>
							</a>
						</li>
					</ul>
				</div>
			</nav>
		)
	}
}