import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  define: {
    // Risolve il problema dei pacchetti gRPC-web/protobuf che cercano 'global'
    'global': 'window',
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      'telemetry-proto-bundle': path.resolve(__dirname, './proto-pkg/dist/index.js'),
    },
  },
  server: {
    port: 3000,
    // host: true permette a Vite di ascoltare su 0.0.0.0 (indispensabile per Docker)
    host: true, 
    strictPort: true,
    hmr: {
      // Assicura che l'Hot Module Replacement punti alla porta corretta sul tuo browser
      clientPort: 3000,
    },
    proxy: {
      // Sostituiamo 'localhost' con 'gateway' (il nome del servizio in docker-compose)
      '/results': 'http://gateway:8080',
      '/set-mode': 'http://gateway:8080',
      '/get-mode': 'http://gateway:8080',
      '/set-size': 'http://gateway:8080', 
      '/get-size': 'http://gateway:8080',
      '/telemetry': 'http://gateway:8080',
    }
  },
  optimizeDeps: {
    include: ['telemetry-proto-bundle']
  }
})