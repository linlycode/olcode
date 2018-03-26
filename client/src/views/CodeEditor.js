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
 * @param { {row: Number, column: Number} } cursor
 */
function cursorToOffset(lines, cursor) {
	let count = 0
	if (cursor.row >= lines.length) {
		throw new Error('invalid row')
	}
	for (let i = 0; i < cursor.row; i += 1) {
		count += lines[i].length
	}
	count += cursor.column
	return count
}

/**
 * @param {String[]} lines
 * @param {Number} offset
 * @return { {startRow: Number, startCol: Number} }
 */
function offsetToPos(lines, offset) {
	let count = 0
	for (let i = 0; i < lines.length; i += 1) {
		count += lines[i].length
		if (count > offset) {
			return { startRow: i, startCol: lines[i].length - (count - offset) }
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
					room.docSync.delete(content, before)
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
		const offset = cursorToOffset(lines, selection.getCursor())
		room.docSync.updateCursor(offset)
	}

	onFocus = () => {}

	syncDoc(editting) {
		if (!this.editor) {
			return
		}
		this.editor.setValue(editting.doc.content)
		const lines = contentTolines(editting.doc.content)
		const userCursors = editting.user_edittings.map(e =>
			Object.assign(
				{
					userID: e.user.id,
					className: peerCursorClassName,
				},
				offsetToPos(lines, e.cursor_pos)
			)
		)

		const { user } = this.props
		if (user) {
			const cursor = userCursors.find(c => c.userID === user.id)
			if (cursor) {
				this.editor.moveCursorTo(cursor.startRow, cursor.startCol)
			}
		}

		this.setState({
			peerCursors: userCursors.filter(c => !user || c.userID !== user.id),
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
