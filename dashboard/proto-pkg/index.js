const messages = require('./telemetry_pb.js');
const services = require('./telemetry_grpc_web_pb.js');

module.exports = {
  ...messages,
  ...services
};