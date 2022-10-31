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
								Dashbord <img src="/icons/house.svg" role="img" alt={this.props.dbName} />
							</a>
						</li>
						<li className="nav-item">
							<a className="nav-link" href={`/${this.props.dbName}/grids`}>
							Grids <img src="/icons/grid-3x3-gap.svg" role="img" alt="Grid" />
							</a>
						</li>
					</ul>

					<h6 className="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted">
					<span>Other grids</span>
					<a className="link-secondary" href="#" aria-label="Add a new report">
						<img src="/icons/plus-circle.svg" role="img" alt="Grid" />
					</a>
					</h6>
					<ul className="nav flex-column mb-2">
						<li className="nav-item">
							<a className="nav-link" href={`/${this.props.dbName}/users`}>
							<span data-feather="file-text" className="align-text-bottom"></span>
							Users <img src="/icons/people.svg" role="img" alt="People" />
							</a>
						</li>
					</ul>
				</div>
			</nav>
		)
	}
}