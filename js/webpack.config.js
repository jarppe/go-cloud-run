const CopyPlugin = require("copy-webpack-plugin")
const { CleanWebpackPlugin } = require("clean-webpack-plugin")
const TerserPlugin = require("terser-webpack-plugin")

const path = require("path")

module.exports = {
  entry:        [
    path.resolve(__dirname, "src/index.js"),
    path.resolve(__dirname, "src/style.scss"),
    path.resolve(__dirname, "src/vendors.scss"),
  ],
  module:       {
    rules: [
      {
        test:    /\.scss$/,
        exclude: /node_modules/,
        use:     [
          {
            loader:  "file-loader",
            options: {
              name: "[name].css",
            },
          },
          "css-loader",
          "sass-loader",
        ],
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
    new CleanWebpackPlugin(),
    new CopyPlugin({
      patterns: [
        { from: "assets", to: "" },
      ],
    }),
  ],
}
