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
}
