const path = require('path')
const CleanWebpackPlugin = require('clean-webpack-plugin')
const StyleLintPlugin = require('stylelint-webpack-plugin')

module.exports = {
  entry: {
    app: './src/index.js'
  },
  externals: {
    jquery: 'jQuery',
    turbolinks: 'Turbolinks'
  },
  optimization: {
    splitChunks: {
      chunks: 'all'
    }
  },
  module: {
    
  },
  plugins: [
    new CleanWebpackPlugin(['dist']),
    new StyleLintPlugin()
  ],
  output: {
    filename: 'js/[name].bundle.js',
    path: path.resolve(__dirname, '../public/dist'),
    publicPath: '/dist/'
  }
}
