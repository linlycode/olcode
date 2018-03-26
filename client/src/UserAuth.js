import User from './User'

export default class UserAuth {
	constructor(httpClient) {
		this.httpClient = httpClient
	}

	login(name) {
		// const resp = await this.httpClient.post('/login', { name })
		return new Promise(resolve => resolve(new User(12, name)))
	}

	logout() {
		return new Promise(resolve => resolve(true))
	}
}
