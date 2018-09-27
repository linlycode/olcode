import 'antd/dist/antd.css'
import dva from 'dva'
import Router from './views/Router'

const app = dva()

app.router(Router)
app.start('#root')
