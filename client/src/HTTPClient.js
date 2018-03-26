import request from 'superagent'

function formatResponse(resp) {
	if (!resp.ok) {
		throw resp.error
	}
	if (resp.body.code !== 0) {
		throw new Error(resp.body)
	}
	return resp.body.data
}

/**
 * @class HTTPClient
 */
export default class HTTPClient {
	constructor(host, port) {
		this.host = host
		this.port = port
	}

	createWebSocket(url, onClose, onMessage) {
		const conn = new WebSocket(`ws://${this.host}:${this.port}${url}`)
		conn.onclose = onClose
		conn.onmessage = onMessage
		return conn
	}

	l(relURL) {
		return `http://${this.host}:${this.port}${relURL}`
	}

	/**
	 * @param {String} url
	 * @param {Boolean} json
	 * @return {request.Request}
	 */
	post(url, json = true) {
		return request.post(this.l(url)).use(req => this.makeRequest(req, json))
	}

	/**
	 * @param {String} url
	 * @return {request.Request}
	 */
	get(url) {
		return request.get(this.l(url)).use(this.makeRequest.bind(this))
	}

	makeRequest(req, json = true) {
		if (json) {
			req.set('Content-Type', 'application/json').accept('json')
		}

		const rawResponseThen = req.then.bind(req)

		req.then = (onFulfilled, onRejected) =>
			rawResponseThen(formatResponse.bind(this)).then(
				onFulfilled,
				(...args) => {
					if (onRejected) {
						onRejected(...args)
					}
				}
			)
	}
}
