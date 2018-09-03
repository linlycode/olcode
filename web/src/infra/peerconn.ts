import log from 'src/infra/log'
import 'webrtc-adapter'
import codec, { Codec } from './codec'

export interface Sender {
	send(msg: string): boolean
}

export interface DataChanCallbacks {
	onopen: EventHandler | null
	onmessage: (msg: string) => void | null
	onerror: (event: ErrorEvent) => void | null
	onclose: EventHandler | null
}

export interface IPeerConn {
	setSender(sender: Sender): void
	connect(): void
	sendData(msg: string): boolean
	closeDataChan(): void
	audioCall(onSuccess: () => void): void
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
	private audioElem: HTMLAudioElement
	private onRemoteAudioAdd: () => void | null
	private codec: Codec
	constructor(c: PeerConnConfig) {
		this.config = c
		log.info(c.iceServer)
		this.pc = new RTCPeerConnection({
			iceServers: [c.iceServer],
		})
		this.pc.onicecandidate = ev => this.onIceCandidate(ev)
		this.pc.ondatachannel = ({ channel }) => this.onDataChCreated(channel)
		this.dataCh = null

		// <audio id="audio2" autoplay></audio>
		this.audioElem = document.createElement('audio')
		this.audioElem.setAttribute('autoplay', 'true')
		this.audioElem.style.visibility = 'hidden'
		document.body.appendChild(this.audioElem)

		this.pc.ontrack = ev => this.onRemoteStream(ev)

		this.codec = codec
	}

	public setSender(sender: Sender): void {
		this.sender = sender
	}

	public connect(): void {
		// createDataChannel must be called before createOffer
		// TODO: add data channel should not be the default behavior of connect
		if (!this.dataCh) {
			this.onDataChCreated(this.pc.createDataChannel("code"))
		}
		this.pc.createOffer().then((offer: RTCSessionDescriptionInit) => {
			this.onLocalSessionCreated(offer)
		})
	}

	public sendData(msg: string): boolean {
		if (this.dataCh === null) {
			return false
		}

		log.info("will send message:", msg)
		this.dataCh.send(this.codec.encode(msg))
		return true
	}

	public closeDataChan(): void {
		if (this.dataCh !== null) {
			this.dataCh.close()
			this.dataCh = null
		}
	}

	public audioCall(onSuccess: () => void): void {
		this.onRemoteAudioAdd = onSuccess
		const constraint = { audio: true }
		navigator.mediaDevices
			.getUserMedia(constraint)
			.then(stream => {
				this.addAudioTracks(stream)
				log.info("audio call is made leading to renegotiation")
				this.connect()
			})
			.catch(e => log.error(e))
	}

	public handlePeerIceCandidate(candidate: RTCIceCandidateInit): boolean {
		this.pc.addIceCandidate(candidate)
		return true
	}

	public handlePeerSdp(message: RTCSessionDescriptionInit): boolean {
		switch (message.type) {
			case "offer":
				log.info('Got offer. Sending answer to peer.')
				this.pc.setRemoteDescription(message, () => null, log.info)
				const reply = () => this.pc.createAnswer().then((answer) => this.onLocalSessionCreated(answer))

				// found audio tracks in offer sdp
				// an answer with audio information should be generated
				if (message.sdp && message.sdp.indexOf('audio') >= 0) {
					this.attachAudio(reply)
				} else {
					reply()
				}
				break
			case "answer":
				log.info('Got answer.')
				this.pc.setRemoteDescription(message, () => null, log.info)
				break
			default:
				return false
		}
		return true
	}

	private onDataChCreated(ch: RTCDataChannel): void {
		log.info("data channel created")
		this.dataCh = ch
		const emptyFunc = () => null
		const cbs = this.config.dataChCallbacks

		this.dataCh.onopen = cbs.onopen || emptyFunc
		this.dataCh.onclose = cbs.onclose || emptyFunc
		this.dataCh.onerror = cbs.onerror || emptyFunc

		let onMsg: (msgEv?: MessageEvent) => void = emptyFunc
		if (cbs.onmessage) {
			onMsg = (ev: MessageEvent) => {
				const data = this.codec.decode(ev.data)
				cbs.onmessage(data)
			}
		}
		this.dataCh.onmessage = onMsg
	}

	private notifyPeer(o: object | null) {
		this.sender.send(JSON.stringify(o))
	}

	private onLocalSessionCreated(desc: RTCSessionDescriptionInit) {
		log.info('local session created:', desc)
		this.pc.setLocalDescription(
			desc,
			() => {
				log.info('sending local desc:', this.pc.localDescription)
				this.notifyPeer(this.pc.localDescription)
			},
			log.info)
	}

	// IceCandidate will be generated from the local
	// It needs to be sent to the peer
	private onIceCandidate(event: RTCPeerConnectionIceEventInit) {
		log.info('icecandidate event: ', event)
		if (!event.candidate) {
			log.info('End of candidates.')
			return
		}
		this.notifyPeer({
			candidate: event.candidate.candidate,
			sdpMLineIndex: event.candidate.sdpMLineIndex,
			sdpMid: event.candidate.sdpMid,
			type: 'candidate',
		})
	}

	private onRemoteStream(ev: RTCTrackEvent) {
		this.audioElem.srcObject = ev.streams[0]
		if (this.onRemoteAudioAdd) {
			this.onRemoteAudioAdd()
		}
	}

	private addAudioTracks(stream: MediaStream): boolean {
		const tracks = stream.getAudioTracks()
		if (tracks.length <= 0) {
			log.error('fail to get audio tracks from stream')
			return false
		}

		log.info('will use audio track:', tracks[0].label)
		stream.getTracks().forEach(track => this.pc.addTrack(track))
		return true
	}

	private attachAudio(onAttached: () => void) {
		const constraint = { audio: true }
		navigator.mediaDevices
			.getUserMedia(constraint)
			.then(stream => {
				// TODO: so may be failed here without external awareness
				if (this.addAudioTracks(stream)) {
					onAttached()
				}
			})
			.catch(e => log.error("fail to append audio track"))
	}
}

export default function MakePeerConnection(c: PeerConnConfig): IPeerConn {
	return new PeerConn(c)
}
