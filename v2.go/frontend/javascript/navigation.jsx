// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Navigation extends React.Component {
	render() {
		return (
			<nav className="navbar navbar-expand-lg bg-light">
				<div className="container-fluid">
					<a className="navbar-brand">{this.props.appName}</a>
					<div className="navbar-text">{this.props.dbName}</div>
					<div className="collapse navbar-collapse" id="navbarNavAltMarkup">
						<div className="navbar-nav">
							<a className="nav-link" href={`/${this.props.dbName}/grids`}>
								<button type="button" className="btn btn-outline-primary btn-sm">
									Grids <img src="/icons/grid-3x3-gap.svg" role="img" alt="Grid" />
								</button>
							</a> 
							<a className="nav-link" href={`/${this.props.dbName}/users`}>
								<button type="button" className="btn btn-outline-primary btn-sm">
									Users <img src="/icons/people.svg" role="img" alt="People" />
								</button>
							</a> 
						</div>
					</div>
					<div className="navbar-nav">
						<a className="nav-link">
							<button type="button"
									className="btn btn-outline-secondary btn-sm"
									onClick={
										() => {
											localStorage.removeItem(`access_token_${this.props.dbName}`)
											location.reload()
										}
									}>
								Log out <img src="/icons/box-arrow-right.svg" role="img" alt="Log out" />
							</button>
						</a>
					</div>
					<div className="navbar-text">{this.props.user}</div>
				</div>
			</nav>
		)
	}
}