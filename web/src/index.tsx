import dva from 'dva'
import './index.less'
import Router from './Router'

const app = dva()

//  Model
// app.model(require('./models/example').default);

app.router(Router)
app.start('#root')
