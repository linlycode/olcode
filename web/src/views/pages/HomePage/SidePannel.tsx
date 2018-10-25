import { Button, Icon, Layout, Row, Tag, Tooltip } from 'antd'
import * as React from 'react'
import styled from 'styled-components'

import { copyToClipboard } from 'src/infra/copy'

const Adress = styled.div`
	border: 1px solid #eee;
	border-radius: 2px;
	padding: 0 15px;
	font-size: 14px;
	line-height: 2;
	display: inline-block;
	width: 200px;
	overflow: auto;
	vertical-align: middle;
	white-space: nowrap;
`

const Sider = styled(Layout.Sider)`
	padding: 10px;
`

const Status = styled.ul`
	list-style: none;
	padding: 0;

	li {
		margin-top: 5px;
	}
`

const loadingText = '...'
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

interface AddonProps {
	onClick: () => void
	text: string
	icon: string
}

export default class SidePannel extends React.Component<any, Props> {
	constructor(props: Props) {
		super(props)
		this.onClickCopyButton = this.onClickCopyButton.bind(this)
	}

	public render() {
		const { audioConnected, codeConnected, onCallBtnClick } = this.props

		return (
			<Sider theme="light" width="300">
				<Row>
					<Adress>{getShareLink(this.props.token)}</Adress>
					<Addon
						text="copy link"
						icon="copy"
						onClick={this.onClickCopyButton}
					/>
					<Addon
						text="audio call"
						icon="video-camera"
						onClick={onCallBtnClick}
					/>
				</Row>
				<Row>
					<Status>
						<li>audio: {this.connectionStatusText(audioConnected)}</li>
						<li>code: {this.connectionStatusText(codeConnected)}</li>
					</Status>
				</Row>
			</Sider>
		)
	}

	private onClickCopyButton() {
		copyToClipboard(getShareLink(this.props.token))
	}

	private connectionStatusText(connected: boolean) {
		return connected ? (
			<Tag color="green">connected</Tag>
		) : (
			<Tag color="red">disconnected</Tag>
		)
	}
}

function Addon(props: AddonProps) {
	return (
		<Tooltip placement="top" title={props.text}>
			<Button
				style={{ marginLeft: '5px' }}
				size="small"
				onClick={props.onClick}
			>
				<Icon type={props.icon} theme="outlined" />
			</Button>
		</Tooltip>
	)
}
