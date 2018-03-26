import React from 'react'
import { hot } from 'react-hot-loader'
import PropTypes from 'prop-types'

import { subscribe } from '../store'
import User from '../User'

import UserNameWidget from './UserNameWidget'
import CodeEditor from './CodeEditor'

class MainPage extends React.Component {
	constructor(props) {
		super(props)
	}

	enterRoomHandler() {
		const { userAuth } = this.props.actors;
		userAuth.createRoom()
			.then(roomID => userAuth.attend(roomID))
	}

	render() {
		const { user } = this.props
		return (
			<div>
				{user && <button onClick={() => this.enterRoomHandler()}>Enter Room</button>}
				<UserNameWidget />
				<CodeEditor />
			</div>
		)
	}
}

MainPage.propTypes = {
	user: PropTypes.instanceOf(User),
}

MainPage.defaultProps = {
	user: null,
}

export default hot(module)(subscribe(['user', 'actors'])(MainPage))
