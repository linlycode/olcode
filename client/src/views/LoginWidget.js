import React from 'react'
import PropTypes from 'prop-types'

export default class LoginWidget extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			username: '',
			loading: false,
		}
	}

	componentWillReceiveProps(nextProps) {
		if (this.props.errorMessage === '' && nextProps.errorMessage !== '') {
			this.setState({ loading: false })
		}
	}

	onUsernameChange = e => {
		this.setState({ username: e.target.value })
	}

	onLogin = () => {
		this.setState({ loading: true })
		this.props.onLogin(this.state.username)
	}

	render() {
		const { errorMessage } = this.props
		return (
			<span>
				<input
					name="name"
					value={this.state.username}
					onChange={this.onUsernameChange}
				/>
				<button disabled={this.state.loading} onClick={this.onLogin}>
					Login
				</button>
				{errorMessage && <span>{errorMessage}</span>}
			</span>
		)
	}
}

LoginWidget.propTypes = {
	onLogin: PropTypes.func.isRequired,
	errorMessage: PropTypes.string,
}

LoginWidget.defaultProps = {
	errorMessage: '',
}
