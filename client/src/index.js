import React from 'react'
import ReactDOM from 'react-dom'
import { createStore } from 'redux'
import { Provider } from 'react-redux'

import App from './App'
import config from './dev.config'
import reducer from './reducer'

const store = createStore(reducer)

ReactDOM.render(
	<Provider store={store}>
		<App config={config} />
	</Provider>,
	document.getElementById('root')
)
