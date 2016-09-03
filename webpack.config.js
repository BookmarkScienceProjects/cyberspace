module.exports = {
  entry: './src/entry.js',
  output: {
    path: __dirname + '/static/dist',
    filename: 'bundle.js'
  },
  module: {
    preLoaders: [
      {
        test: /\.js$/,
        loader: 'eslint-loader',
        exclude: /node_modules/
      }
    ],
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        query: {
          presets: ['es2015'],
          cacheDirectory: true
        }
      }
    ]
  },
  eslint: {
    failOnWarning: false,
    failOnError: false
  }
};
