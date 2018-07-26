import { Route, Router, Switch } from 'dva/router'
import * as React from 'react'
import HomePage from './pages/HomePage/index'


function RouterConfig({ history }: any) {
	return (
		<Router history={history}>
			<Switch>
				<Route path="/" exact={true} component={HomePage} />
			</Switch>
		</Router>
	)
}

export default RouterConfig
