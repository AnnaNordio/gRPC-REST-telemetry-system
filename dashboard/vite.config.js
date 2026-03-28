import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  define: {
    // Indispensabile per evitare "global is not defined" nei pacchetti google-protobuf
    'global': 'window',
  },
  resolve: {
    alias: {
      // @ punta alla tua cartella src per comodità
      '@': path.resolve(__dirname, './src'),
      // Mappa il nome che usi nell'import al bundle generato da Webpack
      // Assicurati che il percorso sia corretto rispetto a dove si trova questo file
      'telemetry-proto-bundle': path.resolve(__dirname, './proto-pkg/dist/index.js'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/results': 'http://localhost:8080',
      '/set-mode': 'http://localhost:8080',
      '/get-mode': 'http://localhost:8080',
      '/set-size': 'http://localhost:8080', 
      '/get-size': 'http://localhost:8080',
      '/telemetry': 'http://localhost:8080',
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      }
    }
  },
  optimizeDeps: {
    // Forza Vite a pre-ottimizzare il bundle per una maggiore velocità di caricamento
    include: ['telemetry-proto-bundle']
  }
})