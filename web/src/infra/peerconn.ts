import 'webrtc-adapter'

export interface Sender {
	send(msg: string): boolean
}

export interface DataChanCallbacks {
	onopen: EventHandler | null
	onmessage: (event: MessageEvent) => void | null
	onerror: (event: ErrorEvent) => void | null
	onclose: EventHandler | null
}

export interface IPeerConn {
	setSender(sender: Sender): void
	connect(): void
	sendData(msg: string): boolean
	closeDataChan(): void
	handlePeerIceCandidate(candidate: RTCIceCandidateInit): boolean
	handlePeerSdp(message: RTCSessionDescriptionInit): boolean
}

export interface IceServerConfig {
	urls: string | string[]
	username: string
	credential: string
}

export interface PeerConnConfig {
	iceServer: IceServerConfig
	dataChCallbacks: DataChanCallbacks
}

class PeerConn implements IPeerConn {
	private pc: RTCPeerConnection
	// sender accessible if connect is called
	private sender: Sender
	private config: PeerConnConfig
	private dataCh: RTCDataChannel | null
	constructor(c: PeerConnConfig) {
		this.config = c
		console.log(c.iceServer)
		this.pc = new RTCPeerConnection({
			iceServers: [c.iceServer],
		})
		this.pc.onicecandidate = ev => this.onIceCandidate(ev)
		this.pc.ondatachannel = ({ channel }) => this.onDataChCreated(channel)
		this.dataCh = null
	}

	public handlePeerIceCandidate(candidate: RTCIceCandidateInit): boolean {
		this.pc.addIceCandidate(candidate)
		return true
	}

	public handlePeerSdp(message: RTCSessionDescriptionInit): boolean {
		switch (message.type) {
			case "offer":
				console.log('Got offer. Sending answer to peer.')
				this.pc.setRemoteDescription(message, () => null, console.log)
				this.pc.createAnswer().then((answer) => this.onLocalSessionCreated(answer))
				break
			case "answer":
				console.log('Got answer.')
				this.pc.setRemoteDescription(message, () => null, console.log)
				break
			default:
				return false
		}
		return true
	}

	public setSender(sender: Sender): void {
		this.sender = sender
	}

	public connect(): void {
		// createDataChannel must be called before createOffer
		this.onDataChCreated(this.pc.createDataChannel("code"))
		this.pc.createOffer().then((offer: RTCSessionDescriptionInit) => {
			this.onLocalSessionCreated(offer)
		})
	}


	public sendData(msg: string): boolean {
		if (this.dataCh === null) {
			return false
		}

		this.dataCh.send(msg)
		return true
	}

	public closeDataChan(): void {
		if (this.dataCh !== null) {
			this.dataCh.close()
			this.dataCh = null
		}
	}

	private onDataChCreated(ch: RTCDataChannel): void {
		console.log("data channel created")
		this.dataCh = ch
		const emptyFunc = () => null
		const cbs = this.config.dataChCallbacks

		this.dataCh.onopen = cbs.onopen || emptyFunc
		this.dataCh.onclose = cbs.onclose || emptyFunc
		this.dataCh.onmessage = cbs.onmessage || emptyFunc
		this.dataCh.onerror = cbs.onerror || emptyFunc
	}

	private notifyPeer(o: object | null) {
		this.sender.send(JSON.stringify(o))
	}

	private onLocalSessionCreated(desc: RTCSessionDescriptionInit) {
		console.log('local session created:', desc)
		this.pc.setLocalDescription(
			desc,
			() => {
				console.log('sending local desc:', this.pc.localDescription)
				this.notifyPeer(this.pc.localDescription)
			},
			console.log)
	}

	// IceCandidate will be generated from the local
	// It needs to be sent to the peer
	private onIceCandidate(event: RTCPeerConnectionIceEventInit) {
		console.log('icecandidate event: ', event)
		if (!event.candidate) {
			console.log('End of candidates.')
			return
		}
		this.notifyPeer({
			candidate: event.candidate.candidate,
			sdpMLineIndex: event.candidate.sdpMLineIndex,
			sdpMid: event.candidate.sdpMid,
			type: 'candidate',
		})
	}
}

export default function MakePeerConnection(c: PeerConnConfig): IPeerConn {
	return new PeerConn(c)
}
