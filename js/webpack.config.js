const CopyPlugin = require("copy-webpack-plugin")
const { CleanWebpackPlugin } = require("clean-webpack-plugin")
const TerserPlugin = require("terser-webpack-plugin")
const MiniCssExtractPlugin = require("mini-css-extract-plugin")

const path = require("path")

module.exports = {
  entry:        path.resolve(__dirname, "src/index.js"),
  module:       {
    rules: [
      {
        test:    /\.scss$/,
        exclude: /node_modules/,
        use:     [
          MiniCssExtractPlugin.loader,
          "css-loader",
          "sass-loader",
        ],
      },
      {
        test:   /\.ttf$/,
        loader: "ignore-loader",
      },
    ],
  },
  optimization: {
    nodeEnv:     "production",
    minimizer:   [
      new TerserPlugin(),
    ],
    splitChunks: {
      cacheGroups: {
        commons: {
          test:               /node_modules/,
          name:               "vendors",
          reuseExistingChunk: true,
          chunks:             "all",
        },
        default: {
          reuseExistingChunk: true,
        },
      },
    },
  },
  output:       {
    filename: "[name].js",
    path:     path.resolve(__dirname, "dist"),
  },
  plugins:      [
    new MiniCssExtractPlugin(),
    new CleanWebpackPlugin(),
    new CopyPlugin({
      patterns: [
        { from: "assets", to: "" },
      ],
    }),
  ],
}
