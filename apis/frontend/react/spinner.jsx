// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class Spinner extends React.Component {
	render() {
		return (
			<span className="spinner-grow spinner-grow-sm ms-1" role="status">
				<span className="visually-hidden">Loading...</span>
			</span>
		)
	}
}
