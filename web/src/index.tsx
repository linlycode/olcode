import dva from 'dva'
import './index.less'
import Router from './Router'

const app = dva()

app.router(Router)
app.start('#root')
