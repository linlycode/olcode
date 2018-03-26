import React from 'react'
import PropTypes from 'prop-types'

import { subscribe } from '../store'
import Room from '../Room'
import RoomManager from '../RoomManager'

class RoomWidget extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			roomID: '',
		}
	}

	onRoomIDChange = e => {
		this.setState({ roomID: e.target.value })
	}

	onEnterRoomClick = () => {
		const { roomManager } = this.props.actors
		this.props.setRoom(roomManager.attend(this.state.roomID))
	}

	onCreateRoomClick = () => {
		const { roomManager } = this.props.actors
		roomManager.create().then(roomID => {
			this.props.setRoom(roomManager.attend(roomID))
		})
	}

	onLeaveRoomClick = () => {
		this.props.actors.roomManager.leave(this.props.room)
	}

	render() {
		const { room } = this.props
		return room ? (
			<span>
				<span>Room ID: </span>
				<input readOnly="readonly" value={room.id} />
				<button onClick={this.onLeaveRoomClick}>Leave Room</button>
			</span>
		) : (
			<span>
				<div>
					<button onClick={this.onCreateRoomClick}>Create Room</button>
				</div>
				<input
					placeholder="Room ID"
					value={this.state.roomID}
					onChange={this.onRoomIDChange}
				/>
				<button onClick={this.onEnterRoomClick}>Enter Room</button>
			</span>
		)
	}
}

RoomWidget.propTypes = {
	actors: PropTypes.shape({
		roomManager: PropTypes.instanceOf(RoomManager),
	}).isRequired,
	room: PropTypes.instanceOf(Room),
	setRoom: PropTypes.func.isRequired,
}

RoomWidget.defaultProps = {
	room: null,
}

export default subscribe(['room', 'actors'])(RoomWidget)
