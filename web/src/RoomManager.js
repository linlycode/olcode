import Room from './Room'
import DocSync from './DocSync'

export default class RoomManager {
	constructor(httpClient) {
		this.httpClient = httpClient
	}

	/**
	 * @return {Promise<String>}
	 */
	create() {
		return this.httpClient
			.post('/api/create_room')
			.then(({ room_id: roomID }) => roomID)
	}

	/**
	 * @param {String} roomID
	 * @return {Room}
	 */
	attend(roomID) {
		const ws = this.httpClient.createWebSocket(
			`/api/ws/attend?room_id=${roomID}`
		)
		return new Room(roomID, new DocSync(ws))
	}
}
