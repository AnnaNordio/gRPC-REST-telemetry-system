import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Ogni chiamata che inizia con /results o /set-mode 
      // viene mandata a Go (porta 8080)
      '/results': 'http://localhost:8080',
      '/set-mode': 'http://localhost:8080',
      '/telemetry': 'http://localhost:8080'
    }
  }
})