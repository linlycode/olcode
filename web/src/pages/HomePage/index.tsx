import * as React from 'react'
import Conn, { ConnConfig } from '../../domain/conn'
import { DataChanCallbacks } from '../../infra/peerconn'
import './App.less'

interface State {
	codeTextareaDisabled: boolean
	codeText: string
}
class App extends React.Component<any, State>{
	private conn: Conn
	constructor(props: any) {
		super(props)
		const cbs: DataChanCallbacks = {
			onclose: this.onDataChanClose.bind(this),
			onerror: this.onDataChanError.bind(this),
			onmessage: this.onDataChanMessage.bind(this),
			onopen: this.onDataChanOpen.bind(this),
		}
		const c: ConnConfig = {
			dataChCallbacks: cbs,
			hostname: window.location.hostname,
			port: 8081,
			// TODO: this should be passed by props
			token: new URLSearchParams(window.location.search).get('token'),
		}

		console.log("token:", c.token)

		this.conn = new Conn(c)
		this.state = {
			codeText: "",
			codeTextareaDisabled: true,
		}
		this.updateCodeText = this.updateCodeText.bind(this)
	}

	public componentDidMount() {
		this.conn.connect()
	}

	public render() {
		return (
			<div className="App">
				<header className="App-header">
					<h1 className="App-title">Welcome olcode</h1>
				</header>
				<textarea id="code" disabled={this.state.codeTextareaDisabled}
					value={this.state.codeText}
					onChange={this.updateCodeText}
					placeholder="Press Start(or be started), enter some text, then press Send." />
			</div>
		)
	}

	private updateCodeText(ev: React.ChangeEvent<HTMLTextAreaElement>) {
		this.setState({ codeText: ev.target.value })
		this.conn.sync(ev.target.value)
	}

	private onDataChanOpen(ev: Event) {
		console.log('Channel opened!!!')
		this.setState({ codeTextareaDisabled: false })
	}

	private onDataChanClose(ev: Event) {
		console.log('Channel closed!!!')
	}

	private onDataChanMessage(ev: MessageEvent) {
		console.log('data channel message: ', ev.data)
		this.setState({ codeText: ev.data })
	}

	private onDataChanError(ev: ErrorEvent) {
		console.log('data channel error: ', ev)
	}
}

export default App
