const devConfig = {
	gatewayPort: 8081
}

const prodConfig = {
	gatewayPort: 443
}

export default process.env.NODE_ENV === 'production' ? prodConfig : devConfig
