const initialState = {
	user: null,
	actors: { userAuth: null },
}

export default function rootReducer(state = initialState, action) {
	switch (action.type) {
		case 'SET_USER':
			return Object.assign({}, state, { user: action.user })
		case 'SET_ACTORS':
			return Object.assign({}, state, { actors: action.actors })
		default:
			return state
	}
}
