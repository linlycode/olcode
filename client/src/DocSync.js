import MsgType from './messages'

export default class DocSync {
	constructor(ws, closeCallback) {
		this.ws = ws
		this.ws.onopen = this.onOpen
		this.ws.onmessage = this.onMessage
		this.ws.onerror = this.onError
		this.ws.onclose = this.onClose
		this.callbacks = {
			onSyncDoc: () => {},
			onRoomDeleted: () => {},
			onRoomUserListChanged: () => {},
		}
		this.closeCallback = closeCallback
		this.connected = false
	}

	/**
	 * @param { { onSyncDoc: Function, onRoomDeleted: Function, onRoomUserListChanged: Function } } callbacks
	 */
	setCallbacks(callbacks) {
		this.callbacks = {
			onSyncDoc: callbacks.onSyncDoc || this.callbacks.onSyncDoc,
			onRoomDeleted: callbacks.onRoomDeleted || this.callbacks.onRoomDeleted,
			onRoomUserListChanged:
				callbacks.onRoomUserListChanged || this.callbacks.onRoomUserListChanged,
		}
	}

	sendJSON(json) {
		this.ws.send(JSON.stringify(json))
	}

	/**
	 * @param {String} text
	 */
	insert(text) {
		if (!this.connected) {
			console.error('push when websocket is not connected')
			return
		}
		this.sendJSON({ msg_type: MsgType.docInsert, text })
	}

	/**
	 * @param {} content
	 * @param {Boolean} before
	 */
	delete(content, before) {
		if (!this.connected) {
			console.error('push when websocket is not connected')
			return
		}
		this.sendJSON({ msg_type: MsgType.docDelete, before, content })
	}

	updateCursor(offset) {
		if (!this.connected) {
			console.error('push when websocket is not connected')
			return
		}
		this.sendJSON({ msg_type: MsgType.moveCursor, offset })
	}

	onOpen = () => {
		this.connected = true
	}

	onMessage = e => {
		const data = JSON.parse(e.data)

		const { onSyncDoc, onRoomDeleted, onRoomUserListChanged } = this.callbacks

		switch (data.msg_type) {
			case MsgType.syncDoc:
				onSyncDoc(data)
				break
			case MsgType.roomDeleted:
				onRoomDeleted(data)
				break
			case MsgType.roomUserListChanged:
				onRoomUserListChanged(data)
				break
			default:
				throw new Error('unknown msg_type', data.msg_type)
		}
		if (this.syncCallback) {
			this.syncCallback()
		}
	}

	onError = e => {
		console.error(e)
	}

	onClose = () => {
		if (this.closeCallback) {
			this.closeCallback()
		}
	}
}
