interface IState {
	state1: string,
	state2: number
}

interface IAction {
	payload: any
}

export default {

		namespace: 'example',

		state: {},

		reducers: {
			save(state: IState, action: IAction) {
				return { ...state, ...action.payload }
			},
		},

	}
