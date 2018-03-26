import MsgType from './messages'

export default class DocSync {
	constructor(ws, closeCallback) {
		this.ws = ws
		this.ws.onopen = this.onOpen
		this.ws.onmessage = this.onMessage
		this.ws.onerror = this.onError
		this.ws.onclose = this.onClose
		this.syncCallback = null
		this.closeCallback = closeCallback
		this.connected = false
	}

	setSyncCallback(callback) {
		this.syncCallback = callback
	}

	sendJSON(json) {
		this.ws.send(JSON.stringify(json))
	}

	modifyDoc(action, diff) {
		if (!this.connected) {
			console.error('push when websocket is not connected')
			return
		}
		let type = null
		switch (action) {
			case 'insert':
				type = MsgType.docInsert
				break
			case 'remove':
				type = MsgType.docDelete
				break
			default:
				throw new Error('unexpected action', action)
		}
		this.sendJSON({ msg_type: type, content: diff })
	}

	updateCursor(offset) {
		if (!this.connected) {
			console.error('push when websocket is not connected')
			return
		}
		this.sendJSON({ msg_type: 'update_cursor', offset })
	}

	onOpen = () => {
		this.connected = true
	}

	onMessage = e => {
		if (this.syncCallback) {
			this.syncCallback(JSON.parse(e.data))
		}
	}

	onError = e => {
		console.log(e)
	}

	onClose = () => {
		if (this.closeCallback) {
			this.closeCallback()
		}
	}
}
