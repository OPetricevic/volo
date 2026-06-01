import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";
import { copyFileSync, mkdirSync, existsSync, readdirSync } from "fs";

// Plugin to copy manifest.json and icons to dist after build
function copyExtensionFiles() {
  return {
    name: "copy-extension-files",
    closeBundle() {
      copyFileSync(
        resolve(__dirname, "public/manifest.json"),
        resolve(__dirname, "dist/manifest.json")
      );
      const iconsDir = resolve(__dirname, "public/icons");
      const distIcons = resolve(__dirname, "dist/icons");
      if (!existsSync(distIcons)) mkdirSync(distIcons, { recursive: true });
      for (const file of readdirSync(iconsDir)) {
        if (file.endsWith(".png")) {
          copyFileSync(resolve(iconsDir, file), resolve(distIcons, file));
        }
      }
    },
  };
}

export default defineConfig({
  plugins: [react(), copyExtensionFiles()],
  resolve: {
    alias: {
      "@": resolve(__dirname, "src"),
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        popup: resolve(__dirname, "src/popup/index.html"),
        background: resolve(__dirname, "src/background/service-worker.ts"),
        content: resolve(__dirname, "src/content/content.ts"),
      },
      output: {
        entryFileNames: "[name].js",
        chunkFileNames: "chunks/[name].[hash].js",
        assetFileNames: "assets/[name].[ext]",
      },
    },
  },
});
