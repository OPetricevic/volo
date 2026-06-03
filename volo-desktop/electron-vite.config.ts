import { defineConfig } from "electron-vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";

export default defineConfig({
  main: {
    build: {
      rollupOptions: {
        input: resolve(__dirname, "src/main/index.ts"),
      },
    },
  },
  preload: {
    build: {
      rollupOptions: {
        input: resolve(__dirname, "src/preload/index.ts"),
        output: {
          format: "cjs",
          entryFileNames: "[name].js",
        },
      },
    },
  },
  renderer: {
    root: "src/renderer",
    build: {
      rollupOptions: {
        input: resolve(__dirname, "src/renderer/index.html"),
      },
    },
    plugins: [react()],
    resolve: {
      alias: {
        "@": resolve(__dirname, "src/renderer"),
      },
    },
  },
});
