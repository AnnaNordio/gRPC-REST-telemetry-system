import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  define: {
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
    host: true, 
    strictPort: true,
    hmr: {
      clientPort: 3000,
    },
    proxy: {
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