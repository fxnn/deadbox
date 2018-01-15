'use strict';

const paths = require('./config/paths');
const polyfills = require.resolve('./config/polyfills');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = {
  // Don't attempt to continue if there are any errors.
  bail: true,
  // In production, we only want to load the polyfills and the app code.
  entry: [polyfills, paths.appIndexJs],
  output: {
    // The build folder.
    path: paths.appBuild,
    filename: 'bundle.js',
    // We inferred the "public path" (such as / or /my-project) from homepage.
    publicPath: paths.servedPath,
  },
  devServer: {
    // server webpack dev server from build directory
    contentBase: paths.appBuild,
  },
  module: {
    loaders: [
      {
        // process ES6 files using babel
        loader: 'babel-loader',
        test: /\.(js|mjs)$/,
      }
    ]
  },
  plugins: [
    // Copy static files over. This makes webpack perform a task otherwise handled by build tools
    new CopyWebpackPlugin([
        { from: paths.appPublic } // to: output.path
    ]),
  ]
};
