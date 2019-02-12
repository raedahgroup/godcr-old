const path = require('path')
const StyleLintPlugin = require('stylelint-webpack-plugin')
const webpack = require('webpack')

module.exports = {
  entry: {
    app: './src/index.js'
  },
  externals: {
  },
  optimization: {
    splitChunks: {
      chunks: 'all'
    }
  },
  module: {
    rules: [
      {
        test: /\.(scss)$/,
        use: [
          {
            // Adds CSS to the DOM by injecting a `<style>` tag
            loader: 'style-loader'
          },
          {
            // Interprets `@import` and `url()` like `import/require()` and will resolve them
            loader: 'css-loader'
          },
          {
            // Loader for webpack to process CSS with PostCSS
            loader: 'postcss-loader',
            options: {
              plugins: function () {
                return [
                  require('autoprefixer')
                ]
              }
            }
          },
          {
            // Loads a SASS/SCSS file and compiles it to CSS
            loader: 'sass-loader'
          }
        ]
      }
    ]
  },
  plugins: [
    new StyleLintPlugin()
  ],
  output: {
    filename: 'js/[name].bundle.js',
    path: path.resolve(__dirname, '../public/dist'),
    publicPath: '/dist/'
  }
}
