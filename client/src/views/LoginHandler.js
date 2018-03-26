import React from 'react'
import PropTypes from 'prop-types'

import { subscribe } from '../store'
import UserAuth from '../UserAuth'

import LoginWidget from './LoginWidget'

class LoginHandler extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			errorMessage: '',
		}
	}

	onLogin = username => {
		const { userAuth } = this.props.actors

		userAuth
			.login(username)
			.then(
				user => this.props.onLoggedIn(user),
				e => this.setState({ errorMessage: e })
			)
	}

	render() {
		return (
			<LoginWidget
				onLogin={this.onLogin}
				errorMessage={this.state.errorMessage}
			/>
		)
	}
}

LoginHandler.propTypes = {
	actors: PropTypes.shape({
		userAuth: PropTypes.instanceOf(UserAuth),
	}).isRequired,
	onLoggedIn: PropTypes.func.isRequired,
}

export default subscribe(['actors'])(LoginHandler)
