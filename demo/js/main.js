'use strict'

/** @type WebSocket*/
let ws;
const peerConnConfig = null
const peerConn = new RTCPeerConnection(peerConnConfig)
peerConn.onicecandidate = onIceCandidate
peerConn.ondatachannel = onDataChannelCreated
let dataChannel
let roomID, peerID

const sendButton = document.getElementById("send-button")
const codeTextarea = document.getElementById("code")
sendButton.onclick = () => {
	sendMessageViaDataCh(codeTextarea.value)
}

function sendMessageViaDataCh(msg) {
	dataChannel.send(msg)
}

function sendMessage(msgObj) {
	ws.send(JSON.stringify(msgObj))
}

// IceCandidate will be generated from the local
// It needs to be sent to the peer
function onIceCandidate(event) {
	console.log('icecandidate event:', event)
	if (event.candidate) {
		sendMessage({
			type: 'candidate',
			label: event.candidate.sdpMLineIndex,
			id: event.candidate.sdpMid,
			candidate: event.candidate.candidate
		})
	} else {
		console.log('End of candidates.')
	}
}

function signalingMessageCallback(message) {
	if (message.type === 'offer') {
		console.log('Got offer. Sending answer to peer.')
		peerConn.setRemoteDescription(
			new RTCSessionDescription(message),
			() => { },
			console.log)
		peerConn.createAnswer(onLocalSessionCreated, console.log)
	} else if (message.type === 'answer') {
		console.log('Got answer.')
		peerConn.setRemoteDescription(
			new RTCSessionDescription(message),
			() => { },
			console.log)
	} else if (message.type === 'candidate') {
		peerConn.addIceCandidate(new RTCIceCandidate({
			candidate: message.candidate
		}))
	}
}

function startWebrtcConnection() {
	console.log('Creating Peer connection + Data Channel')
	const channel = peerConn.createDataChannel('code');
	onDataChannelCreated({ channel });

	console.log('Creating an offer');
	peerConn.createOffer(onLocalSessionCreated, console.log);
}

function onLocalSessionCreated(desc) {
	console.log('local session created:', desc);
	peerConn.setLocalDescription(
		desc,
		() => {
			console.log('sending local desc:', peerConn.localDescription);
			sendMessage(peerConn.localDescription)
		},
		console.log);
}

function onDataChannelCreated({ channel }) {
	dataChannel = channel
	console.log('onDataChannelCreated:', channel);

	channel.onopen = () => {
		console.log('CHANNEL opened!!!');
		codeTextarea.disabled = false
	};

	channel.onclose = () => {
		console.log('Channel closed.');
	}
	channel.onmessage = ({ data }) => {
		console.log('data channel message:', data)
		codeTextarea.value = data
	}
}

function InitWSConn() {
	const onOpen = () => {
		console.log("ws connection opened")
		const token = new URLSearchParams(window.location.search).get('token')
		if (token) {
			ws.send(`HELLO ${token}`)
		} else {
			ws.send("HELLO")
		}
	}
	const onError = () => {
		console.log("ws connection error")
	}
	const onClose = (event) => {
		console.log("ws connection close", event)
	}

	function onMessage(event) {
		console.log("receive msg from server:", event.data)
		const handleError = msg => {
			console.log(msg)
			ws.close()
		}

		/** @type {string} */
		const msg = event.data
		if (msg.startsWith("ACKHELLO")) {
			const tokens = msg.split(" ")
			if (tokens.length !== 4) {
				handleError(`invalid message: ${msg}`)
				return
			}
			let _, success
			[_, success, roomID, peerID] = msg.split(" ")
			if (!success) {
				handleError(`failed ack hello: ${msg}`)
				return
			}
		} else if (msg.startsWith("PEER_JOINED")) {
			startWebrtcConnection()
		} else {
			// we assume all the other message is used for webrtc connection
			// including sdp/icecandidate
			signalingMessageCallback(JSON.parse(msg))
		}
	}

	const port = 8081
	let url = `ws://${window.location.hostname}:${port}/ws`
	ws = new WebSocket(url)
	ws.onopen = onOpen
	ws.onerror = onError
	ws.onclose = onClose
	ws.onmessage = onMessage
}

InitWSConn()
