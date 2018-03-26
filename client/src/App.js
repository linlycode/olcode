import React from 'react'
import PropTypes from 'prop-types'
import HTTPClient from './HTTPClient'

import { subscribe } from './store'
import UserAuth from './UserAuth'
import RoomManager from './RoomManager'

import MainPage from './views/MainPage'

class App extends React.Component {
	constructor(props) {
		super(props)
		const { config } = this.props

		this.httpClient = new HTTPClient(config.host, config.port)
		const actors = {
			userAuth: new UserAuth(this.httpClient),
			roomManager: new RoomManager(this.httpClient),
		}
		this.props.setActors(actors)
	}

	render() {
		return <MainPage />
	}
}

App.propTypes = {
	config: PropTypes.shape({
		protocol: PropTypes.string.isRequired,
		host: PropTypes.string.isRequired,
		port: PropTypes.number.isRequired,
	}).isRequired,
	setActors: PropTypes.func.isRequired,
}

export default subscribe(['actors'])(App)
