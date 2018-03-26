import User from './User'

export default class UserAuth {
	constructor(httpClient) {
		this.httpClient = httpClient
	}

	login(name) {
		return this.httpClient
			.post('/api/login')
			.send({ name })
			.then(({ user_id }) => new User(user_id, name)) // eslint-disable-line
	}

	logout() {
		return this.httpClient.post('/api/logout').then(() => null)
	}

	createRoom() {
		return this.httpClient
			.post('/api/create_room')
			.then(({ room_id }) => room_id) // eslint-disable-line
	}

	attend(roomID) {
		return this.httpClient.createWebSocket(
			`/api/ws/attend?room_id=${roomID}`,
			evt => console.log('websocket will be closed', evt), // eslint-disable-line
			evt => console.log('receive message', evt.data) // eslint-disable-line
		)
	}
}
