const path = require('path');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const { ESBuildMinifyPlugin } = require('esbuild-loader')

module.exports = {
  watch: false,
  watchOptions: {
    ignored: /node_modules/
  },
  mode: 'production',
  devtool: false,
  entry: {
    default: './src/default',
    app: './src/app/app',
    home: './src/home/home',
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: "esbuild-loader",
        }
      },
        /*
      {
        test: /\.m?js$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
          options: {
            presets: ['@babel/preset-env']
          }
        }
      },
        */
      {
        test: /\.(html|svelte)$/,
        use: {
            loader: 'svelte-loader',
            options: {
              hydratable: true,
            },
          },
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/i,
        type: 'asset/resource',
      },
    ]
  },
  resolve: {
    alias: {
      svelte: path.resolve('node_modules', 'svelte'),
      vue: 'vue/dist/vue.min.js',
    },
    extensions: ['.mjs', '.js', '.svelte', 'vue' ],
    mainFields: ['svelte', 'browser', 'module', 'main']
  },
  output: {
    filename: '[name].[contenthash].js',
    chunkFilename: '[name].[contenthash].js',
    path: path.resolve(__dirname, '../../static/js/'),
    publicPath: '/static/js/'
  },
  optimization: {
    minimize: true,
    usedExports: false,
    moduleIds: 'deterministic',
    removeAvailableModules: true,
    flagIncludedChunks: true,
    minimizer: [
      new ESBuildMinifyPlugin({
        target: 'es2015'  // Syntax to compile to (see options below for possible values)
       })
    ],
    splitChunks: {
      chunks: 'all',
    }
  },
  plugins: [
      new CleanWebpackPlugin(),
  ],
};

