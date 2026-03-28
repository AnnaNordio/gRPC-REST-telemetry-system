const path = require('path');

module.exports = {
  mode: 'production',
  entry: './index.js',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'index.js',
    library: {
      type: 'umd', // Universal Module Definition: funziona ovunque
      name: 'TelemetryProto',
    },
    globalObject: 'this',
  },
  // gRPC-web ha bisogno di questi per non impazzire nel bundle
  resolve: {
    fallback: {
      "buffer": false,
      "fs": false,
      "path": false
    }
  }
};