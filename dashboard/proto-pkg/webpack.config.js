const path = require('path');

module.exports = {
  mode: 'production',
  entry: './telemetry_grpc_web_pb.js', 
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
    library: {
      type: 'commonjs' // Permette a Vite di importarlo facilmente
    }
  },
  target: 'web',
  externals: {
    // Lasciamo che Vite gestisca queste librerie core
    'google-protobuf': 'google-protobuf',
    'grpc-web': 'grpc-web'
  }
};