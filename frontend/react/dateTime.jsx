// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class DateTime extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			timeAgo: this.getTimeAgo()
		}
	}

	getTimeAgo() {
		const localNow = new Date
		const localDate = new Date(this.props.dateTime)
		const localNowUTC =  Date.UTC(localNow.getUTCFullYear(),
										localNow.getUTCMonth(),
										localNow.getUTCDate(),
										localNow.getUTCHours(),
										localNow.getUTCMinutes(),
										localNow.getUTCSeconds())
		const localDateUTC =  Date.UTC(localDate.getUTCFullYear(),
										localDate.getUTCMonth(),
										localDate.getUTCDate(),
										localDate.getUTCHours(),
										localDate.getUTCMinutes(),
										localDate.getUTCSeconds())
		const seconds = (localNowUTC - localDateUTC) / 1000
		const MINUTE = 60, HOUR = MINUTE * 60, DAY = HOUR * 24, WEEK = DAY * 7, MONTH = DAY * 30, YEAR = DAY * 365		
		if(seconds < MINUTE) return `${Math.round(seconds)} sec ago`
		if(seconds < HOUR) return `${Math.round(seconds / MINUTE)} min ago`
		if(seconds < DAY) return `${Math.round(seconds / HOUR)} hour ago`
		if(seconds < WEEK) return `${Math.round(seconds / DAY)} day ago`
		if(seconds < MONTH) return `${Math.round(seconds / WEEK)} week ago`
		if(seconds < YEAR) return `${Math.round(seconds / MONTH)} month ago`
		return `${Math.round(seconds / YEAR)} year ago`
	}

	componentDidMount() {
		this.timerID = setInterval(() => { this.setState({ timeAgo: this.getTimeAgo() }) }, 1000)
	}

	componentWillUnmount() {
		clearInterval(this.timerID)
	}

	render() {
		return (
			<span>{this.props.dateTime} <small><em>({this.state.timeAgo})</em></small></span>
		)
	}
}
