import React from 'react'
import PropTypes from 'prop-types'

import User from '../User'
import UserAuth from '../UserAuth'

import { subscribe } from '../store'
import LoginHandler from './LoginHandler'

class UserNameWidget extends React.Component {
	onLoggedIn = user => {
		this.props.setUser(user)
	}

	onLogoutClick = () => {
		const { userAuth } = this.props.actors
		userAuth.logout()
		this.props.setUser(null)
	}

	render() {
		const { user } = this.props
		return user ? (
			<span>
				<span>{user.name}</span>
				<button onClick={this.onLogoutClick}>Logout</button>
			</span>
		) : (
			<LoginHandler onLoggedIn={this.onLoggedIn} />
		)
	}
}

UserNameWidget.propTypes = {
	user: PropTypes.instanceOf(User),
	actors: PropTypes.shape({
		userAuth: PropTypes.instanceOf(UserAuth),
	}).isRequired,
	setUser: PropTypes.func.isRequired,
}

UserNameWidget.defaultProps = {
	user: null,
}

export default subscribe(['user', 'actors'])(UserNameWidget)
