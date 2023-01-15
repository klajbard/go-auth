import { defineConfig, PluginOption } from 'vite'
import { extname, join, resolve } from 'path'
import fs from 'fs'

const HtmlFallback = ({rootDir}: {rootDir: string}): PluginOption => ({
  name: 'html-fallback',
  configureServer(server) {
    server.middlewares.use((req, _, next) => {
      if (req.originalUrl && req.originalUrl.length > 1 && !extname(req.originalUrl)) {
        if (fs.existsSync(join(rootDir, `${req.originalUrl}.html`))) {
          req.url += '.html'
        }
      }
      next()
    })
  }
})

export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        error: resolve(__dirname, '404.html'),
        logout: resolve(__dirname, 'logout.html'),
        serviecs: resolve(__dirname, 'services.html'),
      }
    }
  },
  plugins: [
    HtmlFallback({rootDir: __dirname})
  ]
})
