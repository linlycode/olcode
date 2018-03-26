const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const CleanWebpackPlugin = require('clean-webpack-plugin')

module.exports = {
	entry: {
		app: './src/index.js',
	},
	mode: 'development',
	output: {
		filename: '[name].js',
		path: path.resolve(__dirname, 'dist'),
		publicPath: 'static',
	},
	module: {
		rules: [
			{ test: /\.js$/, exclude: /node_modules/, loader: 'babel-loader' },
			{
				test: /\.css$/,
				use: [{ loader: 'style-loader' }],
			},
		],
	},
	optimization: {
		splitChunks: {
			chunks: 'initial',
		},
	},
	plugins: [
		new CleanWebpackPlugin(['dist']),
		new HtmlWebpackPlugin({ template: './index.html' }),
		new webpack.NamedModulesPlugin(),
		new webpack.HotModuleReplacementPlugin(),
	],
	devServer: {
		// contentBase: './dist',
		hot: true,
	},
}
