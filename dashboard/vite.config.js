import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  define: {
    // Necessario per le librerie protobuf
    'global': 'window',
  },
  server: {
    port: 3000,
    proxy: {
      '/results': 'http://localhost:8080',
      '/set-mode': 'http://localhost:8080',
      '/get-mode': 'http://localhost:8080',
      '/set-size': 'http://localhost:8080', 
      '/get-size': 'http://localhost:8080',
      '/telemetry': 'http://localhost:8080'
    }
  },
  optimizeDeps: {
    include: ['my-grpc-protos']
  }
})