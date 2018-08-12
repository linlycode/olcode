import log from 'src/infra/log'
// TODO: add global symbol for root path
import MakePeerConnection, { DataChanCallbacks, IceServerConfig, IPeerConn, PeerConnConfig, Sender } from "src/infra/peerconn"

// TODO: save this to particular config file
const iConfig: IceServerConfig = {
	credential: 'demotoken2018',
	urls: 'turn:39.105.142.163',
	username: 'demo',
}

export interface ConnConfig {
	port: number
	token: string | null
	hostname: string
	// no more details of webrtc exposed to the external
	dataChCallbacks: DataChanCallbacks
	onRecvToken: (token: string) => void
}

function createPeerConnectionConfig(c: ConnConfig): PeerConnConfig {
	const pConfig: PeerConnConfig = {
		dataChCallbacks: c.dataChCallbacks,
		iceServer: iConfig
	}
	return pConfig
}

export default class Conn implements Sender {
	private c: ConnConfig
	private roomID: string
	private peerID: string
	private ws: WebSocket
	private pc: IPeerConn
	constructor(c: ConnConfig) {
		this.c = c
		this.pc = MakePeerConnection(createPeerConnectionConfig(c))
	}

	public sync(msg: string): boolean {
		return this.pc.sendData(msg)
	}

	// implements the Sender
	public send(msg: string): boolean {
		this.ws.send(msg)
		return true
	}

	public connect() {
		this.ws = new WebSocket(`ws://${this.c.hostname}:${this.c.port}/ws`)
		this.ws.onopen = this.onOpen.bind(this)
		this.ws.onclose = this.onClose.bind(this)
		this.ws.onmessage = this.onMessage.bind(this)
		this.ws.onerror = this.onError.bind(this)
		this.pc.setSender(this)
	}

	private onOpen(ev: Event) {
		log.info('ws connection opened')
		if (this.c.token) {
			this.ws.send(`HELLO ${this.c.token}`)
		} else {
			this.ws.send('HELLO')
		}
	}

	private onError(ev: Event) {
		log.info('ws connection error', ev)
	}

	private onClose(ev: CloseEvent) {
		log.info('ws connection close', ev)
	}

	private onMessage(ev: MessageEvent) {
		log.info("receive msg from server:", ev.data)
		const handleError = (errMsg: string) => {
			log.info(errMsg)
			this.ws.close()
		}

		const msg: string = ev.data
		if (msg.startsWith("ACKHELLO")) {
			const tokens = msg.split(" ")
			if (tokens.length !== 4) {
				handleError(`invalid message: ${msg}`)
				return
			}
			let success
			[, success, this.roomID, this.peerID] = msg.split(" ")
			if (!success) {
				handleError(`failed ack hello: ${msg}`)
				return
			}
			this.c.onRecvToken(`${this.roomID}`)
		} else if (msg.startsWith("PEER_JOINED")) {
			this.pc.connect()
		} else {
			// we assume all the other message is used for webrtc connection
			// including sdp/icecandidate
			this.handleSignalledMessage(JSON.parse(msg))
		}
	}

	// TODO: type of param message should be defined in ts
	private handleSignalledMessage(message: any) {
		if (message.type === 'offer' || message.type === "answer") {
			const sdp: RTCSessionDescriptionInit = {
				sdp: message.sdp,
				type: message.type,
			}
			this.pc.handlePeerSdp(sdp)
		} else if (message.type === 'candidate') {
			const candidate: RTCIceCandidateInit = {
				candidate: message.candidate,
				sdpMLineIndex: message.sdpMLineIndex,
				sdpMid: message.sdpMid,
			}
			this.pc.handlePeerIceCandidate(candidate)
		}
	}
}
