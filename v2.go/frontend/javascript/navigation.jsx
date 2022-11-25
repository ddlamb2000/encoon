// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Navigation extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: "",
			isLoaded: false,
			isLoading: false,
			rows: []
		}
	}

	componentDidMount() {
		this.loadData()
	}

	render() {
		const { isLoading, isLoaded, error, rows } = this.state
		return (
			<nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div className="position-sticky pt-3 sidebar-sticky">
					<ul className="nav flex-column">
						<li className="nav-item">
							<a className="nav-link active border-bottom" aria-current="page" href={`/${this.props.dbName}`}>
								Dashboard <i className="bi bi-box"></i>
							</a>
						</li>
					</ul>
					{isLoading && <Spinner />}
					{error && !isLoading && <div className="alert alert-primary" role="alert">{error}</div>}
					<ul className="nav flex-column mb-2">
						{isLoaded && rows && rows.map(row => 
							<li className="nav-item" key={row.uuid}>
								<a className="nav-link" href={`/${this.props.dbName}/${row.uuid}`}>
									{row.text1} {row.text3 && <i className={`bi bi-${row.text3}`}></i>}
								</a>
							</li>
						)}
					</ul>
				</div>
			</nav>
		)
	}

	loadData() {
		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${UuidGrids}`
		fetch(uri, {
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + this.props.token
			}
		})
		.then(response => {
			const contentType = response.headers.get("content-type")
			if(contentType && contentType.indexOf("application/json") !== -1) {
				return response.json().then(	
					(result) => {
						this.setState({
							isLoading: false,
							isLoaded: true,
							rows: result.rows,
							error: result.error
						})
					},
					(error) => {
						this.setState({
							isLoading: false,
							isLoaded: false,
							rows: [],
							error: error.message
						})
					}
				)
			} else {
				this.setState({
					isLoading: false,
					isLoaded: false,
					rows: [],
					error: `[${response.status}] Internal server issue.`
				})
			}
		})
	}
}