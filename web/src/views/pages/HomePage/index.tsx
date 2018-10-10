import { Layout } from 'antd'
import * as React from 'react'

import config from 'src/config'
import Conn, { ConnConfig } from 'src/domain/conn'
import log from 'src/infra/log'
import { AVCallbacks, DataChanCallbacks } from 'src/infra/peerconn'
import CodeEditor from 'src/views/widgets/Editor'
import SidePannel from './sidepannel'

interface State {
	codeTextareaDisabled: boolean
	codeText: string
	token: string | null,
	audioConnected: boolean,
	codeConnected: boolean,
}

class App extends React.Component<any, State>{
	private conn: Conn
	constructor(props: any) {
		super(props)
		const dataChCallbacks: DataChanCallbacks = {
			onclose: this.onDataChanClose.bind(this),
			onerror: this.onDataChanError.bind(this),
			onmessage: this.onDataChanMessage.bind(this),
			onopen: this.onDataChanOpen.bind(this),
		}
		const avCallbacks: AVCallbacks = {
			onRemoteAudioAdd: this.onRemoteAudioAdd.bind(this)
		}
		const c: ConnConfig = {
			avCallbacks,
			dataChCallbacks,
			hostname: window.location.hostname,
			onRecvToken: (token) => this.onRecvToken(token),
			port: config.gatewayPort,
			ssl: window.location.protocol.startsWith("https"),
			// TODO: this should be passed by props
			token: new URLSearchParams(window.location.search).get('token'),
		}

		log.info("token:", c.token)

		this.conn = new Conn(c)
		this.state = {
			audioConnected: false,
			codeConnected: false,
			codeText: "// Press Start(or be started), enter some text, then press Send.",
			codeTextareaDisabled: true,
			token: c.token || null,
		}
		this.updateCodeText = this.updateCodeText.bind(this)
		this.onCallBtnClick = this.onCallBtnClick.bind(this)
	}

	public componentDidMount() {
		this.conn.connect()
	}

	public render() {
		const { audioConnected, codeConnected } = this.state
		return (
			<Layout>
				<Layout.Header>
					Online Code
				</Layout.Header>
				<Layout>
					<Layout.Sider theme='light' width="16.7%">
						<SidePannel token={this.state.token}
							onCallBtnClick={this.onCallBtnClick}
							audioConnected={audioConnected}
							codeConnected={codeConnected}
						/>
					</Layout.Sider>
					<Layout.Content>
						<CodeEditor
							disabled={this.state.codeTextareaDisabled}
							code={this.state.codeText}
							onCodeChange={this.updateCodeText}
						/>
					</Layout.Content>
				</Layout>
			</Layout>
		)
	}

	private onRecvToken(token: string) {
		this.setState({ token })
	}

	private updateCodeText(code: string) {
		this.setState({ codeText: code })
		this.conn.sync(code)
	}

	private onDataChanOpen(ev: Event) {
		log.info('Channel opened!!!')
		this.setState({ codeTextareaDisabled: false, codeConnected: true })
	}

	private onRemoteAudioAdd() {
		log.info("audio call succeed")
		this.setState({ audioConnected: true })
	}

	private onCallBtnClick() {
		this.conn.audioCall()
	}

	private onDataChanClose(ev: Event) {
		log.info('Channel closed!!!')
	}

	private onDataChanMessage(msg: string) {
		log.info('data channel message: ', msg)
		this.setState({ codeText: msg })
	}

	private onDataChanError(ev: ErrorEvent) {
		log.info('data channel error: ', ev)
	}
}

export default App
