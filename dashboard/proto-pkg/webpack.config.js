const path = require('path');

module.exports = {
  mode: 'production',
  entry: './index.js',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'index.js',
    library: {
      type: 'umd', 
      name: 'TelemetryProto',
    },
    globalObject: 'this',
  },
  resolve: {
    fallback: {
      "buffer": false,
      "fs": false,
      "path": false
    }
  }
};