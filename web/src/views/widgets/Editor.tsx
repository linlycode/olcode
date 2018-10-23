import CodeMirror from 'codemirror'
import 'codemirror/lib/codemirror.css'
import styled from 'styled-components'

import * as React from 'react'

const Wrapper = styled.div`
	height: 100%;

	.CodeMirror {
		height: 100%;
	}
`

interface Props {
	disabled?: boolean,
	code?: string,
	onCodeChange?: (code: string) => void
}

const EditorReferences: Editor[] = []

export default class Editor extends React.Component<any, Props> {
	public props: Props

	private cm: CodeMirror.Editor
	private textAreaElem: HTMLTextAreaElement
	private onCodeChange: (code: string) => void
	private disabled: boolean = true
	private code: string = ""
	private disableChangeCallback: boolean = false

	constructor(props: Props) {
		super(props)
		EditorReferences.push(this)
		this.disabled = props.disabled || this.disabled
		this.code = props.code as string
		this.onCodeChange = props.onCodeChange || (() => { return })
	}

	public componentWillReceiveProps(nextProps: Props) {
		const { disabled, code } = nextProps
		if (disabled !== undefined && disabled !== this.disabled) {
			this.textAreaElem.disabled = disabled
			this.disabled = disabled
		}
		if (code !== undefined && code !== this.code) {
			// setValue will fire the handleCodeChange but we need disable it
			// for it via setting disableChangeCallback
			this.disableChangeCallback = true
			this.cm.setValue(code)
			this.code = code
		}
	}

	public shouldComponentUpdate(nextProps: Props): boolean {
		return false
	}

	public componentDidMount() {
		const container = document.getElementById(this.id()) as HTMLElement
		this.cm = CodeMirror(container, {
			lineNumbers: true,
			value: this.code,
		})
		this.textAreaElem = this.cm.getInputField()
		this.textAreaElem.disabled = this.disabled
		this.cm.on('change', () => this.handleCodeChange())
	}

	public render() {
		return <Wrapper id={this.id()} />
	}

	private handleCodeChange() {
		if (this.disableChangeCallback) {
			this.disableChangeCallback = false
			return
		}
		// TODO: may use only change set
		this.code = this.cm.getValue()
		this.onCodeChange(this.code)
	}

	private id(): string {
		return `codemirror-editor-${EditorReferences.indexOf(this)}`
	}
}
