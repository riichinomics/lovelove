var HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

module.exports = env => {
	if (!env) {
		env = {}
	}
	return {
		entry: ["./src/index.ts"],
		mode: env.production ? 'production' : 'development',
		devtool: "source-map",
		devServer: {
			historyApiFallback: true
		},
		resolve: {
			extensions: [".ts", ".js", ".json"]
		},
		module: {
			rules: [
				{
					test: /\.ts$/,
					exclude: /node_modules/,
					use: [
						{
							loader: "awesome-typescript-loader"
						}
					]
				},
				{
					enforce: "pre",
					test: /\.js$/,
					loader: "source-map-loader"
				},
				{
					test: /\.s[ac]ss$/i,
					use: [
						"style-loader",
						"@teamsupercell/typings-for-css-modules-loader",
						{
							loader: "css-loader",
							options: { modules: true }
						},
						'sass-loader'
					],
				},
				{
					test: /\.(mp3)$/i,
					use: [
						{
							loader: 'file-loader',
						},
					],
				},
				{
					test: /\.(png|jp(e*)g|svg)$/,
					use: [
						{
							loader: 'url-loader',
						}
					]
				}
			]
		},
		optimization: {
			splitChunks: {
				chunks: 'all'
			}
		},
		plugins: [
			new HtmlWebpackPlugin({
				base: "/",
				title: "MUHJONG"
			})
		]
	};
};