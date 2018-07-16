const initialState = {
	user: null,
	actors: { userAuth: null },
	room: null,
}

export default function rootReducer(state = initialState, action) {
	switch (action.type) {
		case 'SET_USER':
			return Object.assign({}, state, { user: action.user })
		case 'SET_ACTORS':
			return Object.assign({}, state, { actors: action.actors })
		case 'SET_ROOM':
			return Object.assign({}, state, { room: action.room })
		default:
			return state
	}
}
