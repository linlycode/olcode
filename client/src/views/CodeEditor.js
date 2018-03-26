import React from 'react'
import PropTypes from 'prop-types'
import AceEditor from 'react-ace'

// import 'brace/theme/monokai'
import 'brace/theme/github'

import 'brace/mode/plain_text'
import 'brace/mode/c_cpp'
import 'brace/mode/java'
import 'brace/mode/golang'
import 'brace/mode/python'
import 'brace/mode/javascript'
import 'brace/mode/jsx'

import { subscribe } from '../store'
import User from '../User'
import Room from '../Room'

const peerCursorClassName = 'peer-cursor'

/**
 * @param {string} content
 */
function contentTolines(content) {
	return content.split('\n')
}

/**
 * @param {String[]} lines
 * @param { {row: Number, column: Number} } pos
 */
function posToOffset(lines, pos) {
	let count = 0
	if (pos.row >= lines.length) {
		throw new Error('invalid row')
	}
	for (let i = 0; i < pos.row; i += 1) {
		count += lines[i].length
	}
	count += pos.column
	return count
}

/**
 * @param {String[]} lines
 * @param {Number} offset
 * @return { {row: Number, column: Number} }
 */
function offsetToPos(lines, offset) {
	let count = 0
	for (let i = 0; i < lines.length; i += 1) {
		count += lines[i].length
		if (count > offset) {
			return { row: i, column: lines[i].length - (count - offset) }
		}
	}
	throw new Error('invalid offset')
}

class CodeEditor extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			lang: 'plain_text',
			peerCursors: [],
		}
		this.editor = null
	}

	componentWillReceiveProps(nextProps) {
		if (!this.props.room && nextProps.room) {
			nextProps.room.docSync.setCallbacks({ onSyncDoc: this.syncDoc })
		}
	}

	onLanguageChange = e => {
		this.setState({ lang: e.target.value })
	}

	onContentChange = (_, e) => {
		const { room } = this.props
		if (!room) {
			return
		}
		const content = e.lines.join('\n')
		switch (e.action) {
			case 'insert':
				room.docSync.insert(content)
				break
			case 'remove':
				{
					const before =
						this.editor.selection.getCursor().column === e.start.column
					room.docSync.delete(content.length, before)
				}
				break
			default:
				throw new Error('unexpected action')
		}
	}

	onCursorChange = selection => {
		const { room } = this.props
		if (!room) {
			return
		}
		const { doc } = this.editor.session
		const lines = doc.getLines(0, doc.getLength())
		const offset = posToOffset(lines, selection.getCursor())
		room.docSync.updateCursor(offset)
	}

	onFocus = () => {}

	syncDoc({ content, cursor_map: cursors }) {
		if (!this.editor) {
			return
		}
		this.editor.setValue(content)
		const lines = contentTolines(content)

		const { user } = this.props
		if (user && user.id in cursors) {
			const pos = offsetToPos(cursors[user.id])
			this.editor.moveCursorTo(pos.row, pos.column)
			delete cursors[user.id]  // eslint-disable-line
		}

		const peerCursors = Object.keys(cursors).map(userID => {
			const pos = offsetToPos(lines, cursors[userID])
			return {
				userID,
				className: peerCursorClassName,
				startRow: pos.row,
				startCol: pos.column,
			}
		})

		this.setState({
			peerCursors,
		})
	}

	render() {
		return (
			<div>
				<select defaultValue="plain_text" onChange={this.onLanguageChange}>
					<option value="plain_text">Plain Text</option>
					<option value="c_cpp">C/C++</option>
					<option value="javascript">Java</option>
					<option value="golang">Go</option>
					<option value="python">Python</option>
					<option value="javascript">Javascript</option>
					<option value="jsx">React</option>
				</select>
				<AceEditor
					value="andres iniesta Lionel"
					ref={e => {
						if (e) {
							this.editor = e.editor
							console.log('editor set')
						} else {
							console.error('editor ref is', e)
						}
					}}
					mode={this.state.lang}
					theme="github"
					markers={this.state.peerCursors}
					onChange={this.onContentChange}
					onFocus={this.onFocus}
					onCursorChange={this.onCursorChange}
				/>
			</div>
		)
	}
}

CodeEditor.propTypes = {
	user: PropTypes.instanceOf(User),
	room: PropTypes.instanceOf(Room),
}

CodeEditor.defaultProps = {
	user: null,
	room: null,
}

export default subscribe(['room'])(CodeEditor)
