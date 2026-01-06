import { resolve } from 'path'
import { existsSync, mkdirSync, cpSync } from 'fs'
import react from '@vitejs/plugin-react'
import { defineConfig, externalizeDepsPlugin } from 'electron-vite'

// Shared alias configuration
const aliases = {
  '@/app': resolve(__dirname, 'app'),
  '@/components': resolve(__dirname, 'components'),
  '@/lib': resolve(__dirname, 'lib'),
  '@/types': resolve(__dirname, 'types'),
  '@/resources': resolve(__dirname, 'resources'),
}

// Plugin to copy assets directory
const copyAssets = () => ({
  name: 'copy-assets',
  writeBundle() {
    const sourceDir = resolve(__dirname, 'app/assets')
    const targetDir = resolve(__dirname, 'out/renderer/assets')

    if (existsSync(sourceDir)) {
      mkdirSync(targetDir, { recursive: true })
      cpSync(sourceDir, targetDir, { recursive: true })
    }
  },
})

export default defineConfig({
  main: {
    build: {
      rollupOptions: {
        input: {
          main: resolve(__dirname, 'lib/main/main.ts'),
        },
      },
    },
    resolve: {
      alias: aliases,
    },
    plugins: [externalizeDepsPlugin()],
  },
  preload: {
    build: {
      rollupOptions: {
        input: {
          preload: resolve(__dirname, 'lib/preload/preload.ts'),
        },
      },
    },
    resolve: {
      alias: aliases,
    },
    plugins: [externalizeDepsPlugin()],
  },
  renderer: {
    root: './app',
    build: {
      rollupOptions: {
        input: {
          index: resolve(__dirname, 'app/index.html'),
        },
      },
    },
    resolve: {
      alias: aliases,
    },
    plugins: [react(), copyAssets()],
  },
})
