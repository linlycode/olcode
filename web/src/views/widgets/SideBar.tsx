import Button from '@material-ui/core/Button'
import { withStyles } from '@material-ui/core/styles'
import * as React from 'react'
import styled from 'styled-components'

import Sider from '../../comps/Sider'

import { copyToClipboard } from '../../infra/copy'

const Token = styled.div`
	padding: 7px 15px;
	border: 1px solid #eee;
	border-radius: 4px;
	box-shadow: 0 0 10px 1px #eee inset;
`


const loadingText = "..."
function getShareLink(token: string | null): string {
	if (token === null) {
		return loadingText
	}
	return `${window.location.origin}/?token=${token}`
}

interface Props {
	token: string | null
	// TODO: it seems no a best type for classes
	classes: any
}

// TODO: fix the style
const styles = {
	bullet: {
		display: 'inline-block',
		margin: '0 2px',
		transform: 'scale(0.8)',
	},
	card: {
		minWidth: 275,
	},
	pos: {
		marginBottom: 12,
	},
	title: {
		fontSize: 14,
		marginBottom: 16,
	},
}


class SideBar extends React.Component<any, Props> {
	constructor(props: Props) {
		super(props)
		this.onClickCopyButton = this.onClickCopyButton.bind(this)
	}

	public onClickCopyButton() {
		copyToClipboard(getShareLink(this.props.token))
	}

	public render() {
		return (
			<Sider>
				<Token>
					{getShareLink(this.props.token)}
				</Token>
				<div>
					<Button size="small"
						onClick={this.onClickCopyButton}>
						Copy
					</Button>
				</div>
			</Sider >
		)
	}
}
export default withStyles(styles)(SideBar)
