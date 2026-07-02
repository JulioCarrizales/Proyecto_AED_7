import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// En desarrollo, Vite corre en el puerto 5173 y redirige las llamadas /api
// al backend de Go (puerto 8080). En producción se compila con `npm run build`
// y el propio servidor de Go sirve los archivos generados.
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
