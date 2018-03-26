import request from 'superagent'

function formatBodyResponse(resp) {
	if (!resp.ok) {
		throw resp.error
	}
	return resp.body
}

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

	l(relURL) {
		if (relURL[0] != "/") {
			relURL = "/".concat(relURL)
		}
		return `http://${this.host}:${this.port}${relURL}`
	}

  /**
   * @param {String} url
   * @param {Boolean} json
   * @return {request.Request}
   */
	post(url, json = true) {
		return request.post(this.l(url)).use(req => this._makeRequest(req, json))
	}

  /**
   * @param {String} url
   * @return {request.Request}
   */
	get(url) {
		return request.get(this.l(url)).use(this._makeRequest.bind(this))
	}

	_makeRequest(req, json = true) {
		if (json) {
			req.set('Content-Type', 'application/json').accept('json')
		}

		const rawResponseThen = req.then.bind(req)

		req.then = (onFulfilled, onRejected) =>
			rawResponseThen(formatResponse.bind(this))
				.then(onFulfilled, (...args) => {
					onRejected && onRejected(...args)
				})
	}
}