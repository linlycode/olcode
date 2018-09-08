const devConfig = {
	gatewayPort: 8081
}

const prodConfig = {
	gatewayPort: 80
}

export default process.env.NODE_ENV === 'production' ? prodConfig : devConfig
