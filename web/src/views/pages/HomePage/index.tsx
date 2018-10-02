import * as React from 'react'
import config from 'src/config'
import Conn, { ConnConfig } from 'src/domain/conn'
import log from 'src/infra/log'
import { AVCallbacks, DataChanCallbacks } from 'src/infra/peerconn'
import * as Layouts from 'src/styles/layouts'
import SideBar from 'src/views/widgets/SideBar'
import TopBar from 'src/views/widgets/TopBar'
import styled from 'styled-components'

const Wrapper = styled.div`
	text-align: center;
`

const Content = styled.div`
	position: fixed;
	top: ${Layouts.topbarHeight}px;
	right: 0;
	bottom: 0;
	left: ${Layouts.sideBarWidth}px;
	overflow: auto;
`

const Textarea = styled.textarea`
	margin: 100px;
	width: 80%;
	height: 300px;
	border-color: #e6e6e6;
	padding: 15px;
	font-size: 14px;
`

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
			codeText: "",
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
		return (
			<Wrapper>
				<TopBar title="Online Code" />
				<SideBar token={this.state.token}
					onCallBtnClick={this.onCallBtnClick}
				/>
				<Content>
					<p>audio connection status: {this.connectionStatusText(this.state.audioConnected)}</p>
					<p>code connection status: {this.connectionStatusText(this.state.codeConnected)}</p>
					<Textarea
						disabled={this.state.codeTextareaDisabled}
						value={this.state.codeText}
						onChange={this.updateCodeText}
						placeholder="Press Start(or be started), enter some text, then press Send." />
				</Content>
			</Wrapper>
		)
	}


	private connectionStatusText(connected: boolean): string {
		return connected ? 'connected' : 'disconnected'
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
