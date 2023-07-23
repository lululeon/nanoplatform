import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  return {
    plugins: [react()],
    define: {
      // only expose what's necessary
      'process.env': {
        GQL_ENDPOINT: env.VITE_GQL_ENDPOINT,
      },
    },
  }
})
