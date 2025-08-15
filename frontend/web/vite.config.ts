import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: [
        'vue',
        'vue-router',
        'pinia'
      ],
      dts: true
    }),
    Components({
      resolvers: [ElementPlusResolver()],
      dts: true
    })
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '~': resolve(__dirname, '.')
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        charset: false
      }
    }
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
    cors: true,
    proxy: {
      '/v1/user': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      '/v1/resume': {
        target: 'http://localhost:8000',
        changeOrigin: true
      },
      '/v1/file': {
        target: 'http://localhost:8001',
        changeOrigin: true
      },
      '/v1/ai': {
        target: 'http://localhost:8002',
        changeOrigin: true
      },
      '/v1/parser': {
        target: 'http://localhost:8003',
        changeOrigin: true
      }
    }
  },
  build: {
    target: 'es2015',
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
    rollupOptions: {
      output: {
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: '[ext]/[name]-[hash].[ext]'
      }
    }
  }
})
