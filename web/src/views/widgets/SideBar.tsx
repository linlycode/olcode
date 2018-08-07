import Button from '@material-ui/core/Button'
import Card from '@material-ui/core/Card'
import CardActions from '@material-ui/core/CardActions'
import CardContent from '@material-ui/core/CardContent'
import { withStyles } from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'
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
			<Card>
				<CardContent>
					<Typography className={this.props.classes.title} color="textSecondary">
						{getShareLink(this.props.token)}
					</Typography>
				</CardContent>
				<CardActions>
					<Button size="small"
						onClick={this.onClickCopyButton}>
						Copy
					</Button>
				</CardActions>
			</Card >
		)
	}
}
export default withStyles(styles)(SideBar)
