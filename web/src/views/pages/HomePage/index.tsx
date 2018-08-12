import * as React from 'react'
import Conn, { ConnConfig } from 'src/domain/conn'
import log from 'src/infra/log'
import { DataChanCallbacks } from 'src/infra/peerconn'
import SideBar from 'src/views/widgets/SideBar'
import TopBar from 'src/views/widgets/TopBar'
import style from './App.less'

interface State {
	codeTextareaDisabled: boolean
	codeText: string
	token: string | null
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
			onRecvToken: (token) => this.onRecvToken(token),
			port: 8081,
			// TODO: this should be passed by props
			token: new URLSearchParams(window.location.search).get('token'),
		}
		if (!c.token) {
			this.setState({ token: c.token })
		}

		log.info("token:", c.token)

		this.conn = new Conn(c)
		this.state = {
			codeText: "",
			codeTextareaDisabled: true,
			token: null
		}
		this.updateCodeText = this.updateCodeText.bind(this)
	}

	public componentDidMount() {
		this.conn.connect()
	}

	public render() {
		return (
			<div className={style.app}>
				<TopBar title="Online Code" />
				<textarea className={style.code}
					disabled={this.state.codeTextareaDisabled}
					value={this.state.codeText}
					onChange={this.updateCodeText}
					placeholder="Press Start(or be started), enter some text, then press Send." />
				<SideBar token={this.state.token} />
			</div>
		)
	}

	private onRecvToken(token: string) {
		this.setState({ token })
	}

	private updateCodeText(ev: React.ChangeEvent<HTMLTextAreaElement>) {
		this.setState({ codeText: ev.target.value })
		this.conn.sync(ev.target.value)
	}

	private onDataChanOpen(ev: Event) {
		log.info('Channel opened!!!')
		this.setState({ codeTextareaDisabled: false })
	}

	private onDataChanClose(ev: Event) {
		log.info('Channel closed!!!')
	}

	private onDataChanMessage(ev: MessageEvent) {
		log.info('data channel message: ', ev.data)
		this.setState({ codeText: ev.data })
	}

	private onDataChanError(ev: ErrorEvent) {
		log.info('data channel error: ', ev)
	}
}

export default App
