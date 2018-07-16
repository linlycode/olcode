import User from './User'

export default class UserAuth {
	constructor(httpClient) {
		this.httpClient = httpClient
	}

	/**
	 * @return {Promise<User>}
	 */
	loginFromCookie() {
		return this.httpClient
			.get('/api/login')
			.then(({ user_id: id, name }) => new User(id, name))
	}

	login(name) {
		return this.httpClient
			.post('/api/login')
			.send({ name })
			.then(({ user_id: id }) => new User(id, name))
	}

	logout() {
		return this.httpClient.post('/api/logout').then(() => null)
	}
}
