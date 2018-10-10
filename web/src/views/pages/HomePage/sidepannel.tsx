import { Button, Row, Tag } from 'antd'
import * as React from 'react'

import { copyToClipboard } from 'src/infra/copy'

const loadingText = "..."
function getShareLink(token: string | null): string {
	if (token === null) {
		return loadingText
	}
	return `${window.location.origin}/?token=${token}`
}

interface Props {
	onCallBtnClick: () => void
	token: string | null
	audioConnected: boolean
	codeConnected: boolean
}

export default class SidePannel extends React.Component<any, Props> {
	constructor(props: Props) {
		super(props)
		this.onClickCopyButton = this.onClickCopyButton.bind(this)
	}

	public render() {
		const {
			audioConnected,
			codeConnected,
			onCallBtnClick
		} = this.props

		return (
			<React.Fragment>
				<Row>
					<Tag>
						{getShareLink(this.props.token)}
					</Tag>
				</Row>
				<Row>
					<Button
						size="small"
						type="primary"
						onClick={this.onClickCopyButton}>
						Copy Link
					</Button>
					<Button
						size="small"
						type="primary"
						onClick={onCallBtnClick}>
						Audio call
					</Button>
				</Row>
				<Row>
					<p>
						audio: {this.connectionStatusText(audioConnected)}
						<br />
						code: {this.connectionStatusText(codeConnected)}
					</p>
				</Row>
			</React.Fragment>
		)
	}

	private onClickCopyButton() {
		copyToClipboard(getShareLink(this.props.token))
	}

	private connectionStatusText(connected: boolean): string {
		return connected ? 'connected' : 'disconnected'
	}
}
