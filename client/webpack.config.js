var HtmlWebpackPlugin = require("html-webpack-plugin");
const path = require("path");

module.exports = env => {
	if (!env) {
		env = {};
	}
	return {
		entry: ["./src/global.sass", "./src/index.tsx"],
		mode: env.production ? "production" : "development",
		devtool: "source-map",
		devServer: {
			historyApiFallback: true,
			host: "0.0.0.0",
			port: 8080,
			// disableHostCheck: true,
		},
		resolve: {
			extensions: [".ts", ".tsx", ".js", ".json", ".sass", ".scss"]
		},
		module: {
			rules: [
				{
					test: /\.tsx?$/,
					exclude: /node_modules/,
					use: [
						{
							loader: "ts-loader"
						},
						{
							loader: "astroturf/loader",
							options: { extension: ".module.scss" },
						},
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
						"sass-loader"
					],
				},
				{
					test: /\.(mp3)$/i,
					use: [
						{
							loader: "file-loader",
						},
					],
				},
				{
					test: /\.(png|jp(e*)g|svg)$/,
					use: [
						{
							loader: "url-loader",
						}
					]
				},
			]
		},
		optimization: {
			splitChunks: {
				chunks: "all"
			}
		},
		plugins: [
			new HtmlWebpackPlugin({
				base: "/",
				title: "Hanafuda"
			})
		]
	};
};
