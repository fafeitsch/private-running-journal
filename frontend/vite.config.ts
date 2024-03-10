import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/dist/vite';
import {PrimeVueResolver} from 'unplugin-vue-components/dist/resolvers';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), Components({
    resolvers: [
      PrimeVueResolver()
    ]
  })],
  build: {
    rollupOptions: {external: ['journal/JournalPage.vue']}
  }
})
