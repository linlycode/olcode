import { connect } from 'react-redux'

import { setUser, setActors } from './actions'

const userMapper = {
	mapState: state => ({ user: state.user }),
	mapDispatch: dispatch => ({
		setUser: user => dispatch(setUser(user)),
	}),
}

const actorsMapper = {
	mapState: state => ({ actors: state.actors }),
	mapDispatch: dispatch => ({
		setActors: actors => dispatch(setActors(actors)),
	}),
}

const mappers = new Map([['user', userMapper], ['actors', actorsMapper]])

/**
 * @return {string[]}
 */
export function getAvailableObjects() {
	return Array.from(mappers.keys())
}

/**
 * @param {string[]} objects
 * @return {function((ReactComponent) => ReactComponent)}
 */
export function subscribe(objects) {
	function mapToProps(mapperName) {
		// target is state or dispatch
		return target =>
			objects.reduce((props, name) => {
				const m = mappers.get(name)
				if (m === undefined) {
					throw new Error(`object "${name}" not in store`)
				}
				const mapper = m[mapperName]
				return Object.assign(props, mapper(target))
			}, {})
	}

	return connect(mapToProps('mapState'), mapToProps('mapDispatch'))
}
