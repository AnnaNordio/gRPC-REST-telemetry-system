// Importa i file generati da protoc (CommonJS)
const messages = require('./telemetry_pb.js');
const services = require('./telemetry_grpc_web_pb.js');

// Esporta tutto in un unico oggetto che Vite potrà leggere
module.exports = {
  ...messages,
  ...services
};