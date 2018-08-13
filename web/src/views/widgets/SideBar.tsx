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
	// TODO: it seems no a best type for classes
	classes: any
	onCallBtnClick: () => void
	token: string | null
}

// TODO: fix the style
const styles = {
	button: {
		margin: '0 20px 0 0'
	},
	linkContainer: {
		marginBottom: '20px',
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
				<Token className={this.props.classes.linkContainer}>
					{getShareLink(this.props.token)}
				</Token>
				<div >
					<Button size="small"
						variant="contained"
						color="primary"
						className={this.props.classes.button}
						onClick={this.onClickCopyButton}>
						Copy Link
					</Button>
					<Button size="small"
						variant="contained"
						color="primary"
						className={this.props.classes.button}
						onClick={this.props.onCallBtnClick}>
						Audio call
					</Button>
				</div>
			</Sider >
		)
	}
}
export default withStyles(styles)(SideBar)
