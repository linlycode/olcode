import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Typography from '@material-ui/core/Typography'
import * as React from 'react'

interface Props {
	title: string
}

function TopBar ({title}: Props) {
	return (
		<AppBar>
			<Toolbar>
				<Typography  variant="title" color="inherit" >{title}</Typography>
			</Toolbar>
		</AppBar>
	)
}

export default TopBar
