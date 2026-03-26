import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000, // Imposta la porta del frontend a 3000
    strictPort: true,
    proxy: {
      '/results': 'http://localhost:8080',
      '/set-mode': 'http://localhost:8080',
      '/get-mode': 'http://localhost:8080',
      '/set-size': 'http://localhost:8080', 
      '/get-size': 'http://localhost:8080',
      '/telemetry': 'http://localhost:8080'
    }
  }
})