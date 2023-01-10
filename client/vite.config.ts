import { defineConfig } from 'vite'
import { resolve } from 'path'

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        error: resolve(__dirname, '404.html'),
        logout: resolve(__dirname, 'logout/index.html'),
        serviecs: resolve(__dirname, 'services/index.html'),
      }
    }
  }
})
