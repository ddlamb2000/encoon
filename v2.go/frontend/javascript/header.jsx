// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Header extends React.Component {
	render() {
		return (
            <header className="navbar sticky-top bg-light flex-md-nowrap p-0 shadow">
                <a className="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#">{this.props.appName}</a>
                <button className="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"></span>
                </button>
                <input className="form-control form-control-dark w-100 rounded-0 border-0" type="text" placeholder="Search" aria-label="Search" />
                <div className="navbar-text">
                    <div className="nav-item text-nowrap px-4">
                        {this.props.user}
                    </div>
                </div>
                <div className="navbar-nav">
                    <div className="nav-item text-nowrap">
                        <a className="nav-link px-3">
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
                </div>
            </header>
		)
	}
}