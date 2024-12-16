// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridCellDisplay extends React.Component {
	render() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			return (
				<div id={"richtext-" + this.props.id}>
					<div id={this.props.id}
							typeuuid={this.props.typeUuid}
							uuid={this.props.uuid}
							ref={this.props.inputRef}
							column={this.props.columnName}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				</div>
			)
		}
		else {
			const variantMonospace = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : ""
			return (
				<span className={variantMonospace}>
					{this.getCellValue(this.props.typeUuid, this.props.value)} {this.props.icon && <i className={`bi bi-${this.props.icon}`}></i>}
				</span>
			)
		}
	}

	getCellValue(typeUuid, value) {
		switch(typeUuid) {
			case UuidPasswordColumnType: return '*****'
			case UuidBooleanColumnType: return value == 'true' ? '✔︎' : ''
		}
		return value
	}

	componentDidMount() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.createRichTextField(this.props.id, this.props.value, true)
		}
	}

	componentWillUnmount() {
		if(trace) console.log("[GridCellDisplay.componentWillUnmount()]")
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.deleteRichTextField(this.props.id)
		}
	}
}
