import React from 'react'
import { hot } from 'react-hot-loader'
import PropTypes from 'prop-types'

import { subscribe } from '../store'
import User from '../User'

import UserNameWidget from './UserNameWidget'
import CodeEditor from './CodeEditor'

function MainPage(props) {
	const { user } = props
	return (
		<div>
			{user && <button>Enter Room</button>}
			<UserNameWidget />
			<CodeEditor />
		</div>
	)
}

MainPage.propTypes = {
	user: PropTypes.instanceOf(User),
}

MainPage.defaultProps = {
	user: null,
}

export default hot(module)(subscribe(['user'])(MainPage))
