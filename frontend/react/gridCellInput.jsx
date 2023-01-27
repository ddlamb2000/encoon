// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridCellInput extends React.Component {
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
			const className = this.props.type == "checkbox" ? "form-check-input" : "form-control"
			return (
				<input type={this.props.type}
						typeuuid={this.props.typeUuid}
						className={className + " form-control-sm rounded-2 shadow gap-2 p-1 " + this.props.variantReadOnly}
						name={this.props.uuid}
						uuid={this.props.uuid}
						column={this.props.columnName}
						readOnly={this.props.readonly}
						defaultChecked={this.props.checkedBoolean}
						defaultValue={this.props.value}
						ref={this.props.inputRef}
						onInput={() => this.props.onEditRowClick(this.props.uuid)} />
			)
		}
	}

	componentDidMount() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.createRichTextField(this.props.id, this.props.value, false)
		}
	}

	componentWillUnmount() {
		if(trace) console.log("[GridCellInput.componentWillUnmount()]")
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.deleteRichTextField(this.props.id)
		}
	}
}
