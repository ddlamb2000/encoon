// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridCellDropDown extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: "",
			isLoaded: false,
			isLoading: false,
			rows: [],
		}
	}

	render() {
		const { isLoading, error, rows } = this.state
		const { gridPromptUuid, values, referencedValuesAdded, referencedValuesRemoved } = this.props
		const referencedValuesIncluded = error ? [] : values.
			concat(referencedValuesAdded).
			filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		const referencedValuesNotIncluded = error ? [] : rows.
			concat(referencedValuesRemoved).
			filter(ref => !referencedValuesAdded.map(ref => ref.uuid).includes(ref.uuid)).
			filter(ref => !referencedValuesIncluded.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		const countRows = referencedValuesNotIncluded ? referencedValuesNotIncluded.length : 0
		if(trace) {
			console.log("[GridCellDropDown.render()] this.props.columnName=", this.props.columnName)
			console.log("[GridCellDropDown.render()] values=", values)
			console.log("[GridCellDropDown.render()] referencedValuesAdded=", referencedValuesAdded)
			console.log("[GridCellDropDown.render()] referencedValuesRemoved=", referencedValuesRemoved)
			console.log("[GridCellDropDown.render()] referencedValuesIncluded=", referencedValuesIncluded)
			console.log("[GridCellDropDown.render()] referencedValuesNotIncluded=", referencedValuesNotIncluded)
		}
		return (
			<ul className="list-unstyled mb-0">
				{referencedValuesIncluded.map(ref => 
					<li key={ref.uuid}>
						<span>
							<button type="button" className="btn text-danger btn-sm mx-0 p-0"
									onClick={() => 
										this.props.onRemoveReferencedValueClick({fromUuid: this.props.uuid, 
																				 columnUuid: this.props.columnUuid,
																				 owned: this.props.owned,
																				 columnName: this.props.columnName,
																				 toGridUuid: this.props.gridPromptUuid,
																				 uuid: ref.uuid,
																				 displayString: ref.displayString})}>
								<i className="bi bi-box-arrow-down pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				)}
				{gridPromptUuid &&
					<li>
						<input type="search"
								className="form-control form-control-sm rounded-2 shadow gap-2 p-1"
								autoComplete="false"
								placeholder="Search..."
								onInput={(e) => this.loadDropDownData(gridPromptUuid, e.target.value)} />
					</li>
				}
				{isLoading && <li><Spinner /></li>}
				{error && !isLoading && <li className="alert alert-danger" role="alert">{error}</li>}
				{referencedValuesNotIncluded && countRows > 0 && referencedValuesNotIncluded.map(ref => (
					<li key={ref.uuid}>
						<span>
							<button type="button"
									className="btn text-success btn-sm mx-0 p-0"
									onClick={() => this.props.onAddReferencedValueClick({fromUuid: this.props.uuid, 
																						 columnUuid: this.props.columnUuid,
																						 owned: this.props.owned,
																						 columnName: this.props.columnName,
																						 toGridUuid: this.props.gridPromptUuid, 
																						 uuid: ref.uuid, 
																						 displayString: ref.displayString})}>
								<i className="bi bi-box-arrow-up pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				))}
			</ul>
		)
	}

	loadDropDownData(gridPromptUuid, value) {
		this.setState({isLoading: true})
		if(value.length > 0) {
			const uri = `/${this.props.dbName}/api/v1/${gridPromptUuid}`
			fetch(uri, {
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json',
					'Authorization': 'Bearer ' + this.props.token
				}
			})
			.then(response => {
				const contentType = response.headers.get("content-type")
				if(contentType && contentType.indexOf("application/json") != -1) {
					return response.json().then(	
						(result) => {
							if(result.response != undefined) {
								this.setState({
									isLoading: false,
									isLoaded: true,
									rows: result.response.rows,
									error: result.response.error
								})
							} else {
								this.setState({
									isLoading: false,
									isLoaded: true,
									error: result.error
								})
							}
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
		} else {
			this.setState({
				isLoading: false,
				isLoaded: false,
				rows: [],
				error: ""
			})
		}
	}
}
