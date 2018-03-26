import User from './User'

export default class UserAuth {
	constructor(httpClient) {
		this.httpClient = httpClient
	}

	login(name) {
		return this.httpClient.post('/api/login')
			.send({ name })
			.then(({ user_id }) => new User(user_id, name))
	}

	logout() {
		return new Promise(resolve => resolve(true))
	}

	createRoom() {
		return this.httpClient
			.post('/api/create_room')
			.then(({ room_id }) => {
				return room_id
			})
	}

	attend(roomID) {
		const conn = new WebSocket(`ws://localhost:5432/api/ws/attend?room_id=${roomID}`);
		conn.onclose = function (evt) {
			console.log("websocket will be closed", evt)
		};
		conn.onmessage = function (evt) {
			console.log("receive message, data=%v", evt.data)
		};
		return conn
	}
}
