import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/dist/vite";
import { PrimeVueResolver } from "unplugin-vue-components/dist/resolvers";
import license from "rollup-plugin-license"
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(process.env.npm_package_version),
  },
  plugins: [
    vue(),
    Components({
      resolvers: [PrimeVueResolver()],
    }),
  ],
  build: {
    rollupOptions: {
      plugins: [
        license({
          thirdParty: {
            output: path.resolve(__dirname, "./src/about/vendor.LICENSE.txt"),
          },
        }),
      ],
    },
  },
});
