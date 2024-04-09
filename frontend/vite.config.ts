import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import Components from "unplugin-vue-components/dist/vite";
import { PrimeVueResolver } from "unplugin-vue-components/dist/resolvers";
import license from "rollup-plugin-license"
import path from "path"
import * as child from "child_process";

const commitHash = child.execSync("git rev-parse --short HEAD").toString();

// https://vitejs.dev/config/
export default defineConfig({
  define: {
    __APP_VERSION__: JSON.stringify(process.env.npm_package_version),
    __COMMIT_HASH__: JSON.stringify(commitHash)
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
