import React from 'react'
import AceEditor from 'react-ace'

import 'brace/theme/monokai'

import 'brace/mode/plain_text'
import 'brace/mode/c_cpp'
import 'brace/mode/java'
import 'brace/mode/golang'
import 'brace/mode/python'
import 'brace/mode/javascript'
import 'brace/mode/jsx'

export default class CodeEditor extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			lang: 'plain_text',
		}
	}
	onLanguageChange = e => {
		this.setState({ lang: e.target.value })
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
				<AceEditor mode={this.state.lang} theme="monokai" />
			</div>
		)
	}
}
