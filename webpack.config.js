import webpack from 'webpack'

export default {
	entry: './js/index.js',
	output: {
		filename: 'tab.cjs',
		libraryTarget: 'commonjs',
	},
	target: 'node',
	mode: 'development',
}