import { RouterAPI } from 'dva'
import { Route, Router, Switch } from 'dva/router'
import * as React from 'react'
import HomePage from './pages/HomePage'


function RouterConfig({ history }: RouterAPI) {
	return (
		<Router history={history}>
			<Switch>
				<Route path="/" exact={true} component={HomePage} />
			</Switch>
		</Router>
	)
}

export default RouterConfig
