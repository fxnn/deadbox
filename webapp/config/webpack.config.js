'use strict';

const paths = require('./paths');
const polyfills = require.resolve('./polyfills');
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
  resolve: {
    modules: [
      paths.appSrc, // HINT: allow for absolute path ES6 imports from 'src' directory
      paths.appNodeModules
    ]
  },
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          // HINT: style-loader generates <style> tags, css-loader just resolves the CSS name to a URL
          // TODO: don't use style tags, but https://github.com/webpack-contrib/extract-text-webpack-plugin
          { loader: "style-loader" },
          { loader: "css-loader" }
        ]
      },
      {
        test: /\.js$/,
        exclude: /(node_modules|bower_components)/,
        use: {
          loader: 'babel-loader',
          options: {
            // NOTE that we're not transpiling to ES5, but just resolving modules.
            // Therefore, browsers need to know ES6, but don't need to know module loading.
            plugins: ["transform-es2015-modules-commonjs"]
          }
        }
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
