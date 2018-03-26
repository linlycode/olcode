export function setUser(user) {
	return {
		type: 'SET_USER',
		user,
	}
}

export function setActors(actors) {
	return {
		type: 'SET_ACTORS',
		actors,
	}
}
