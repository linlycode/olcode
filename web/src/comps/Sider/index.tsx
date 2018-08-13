import * as React from 'react'
import * as Layouts from 'src/styles/layouts'
import styled from 'styled-components'

const Wrapper = styled.aside`
	position: fixed;
	top: ${Layouts.topbarHeight}px;
	width: ${Layouts.sideBarWidth}px;
	left: 0;
	bottom: 0;
	background: white;
	border-right: 1px solid #eee;
	padding: 15px;
`

interface SiderProps extends React.HTMLAttributes<HTMLDivElement> {
	collapsed?: boolean
	defaultCollapsed?: boolean
	reverseArrow?: boolean
	trigger?: React.ReactNode
	width?: number | string
}


export default class Sider extends React.Component<SiderProps> {
	public state = {

	}

	public render() {
		return <Wrapper>{this.props.children}</Wrapper>
	}
}
