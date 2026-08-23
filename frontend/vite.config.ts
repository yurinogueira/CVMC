import { defineConfig, type Plugin } from "vite";
import react from "@vitejs/plugin-react";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function spaFallbackPlugin(): Plugin {
  return {
    name: "spa-fallback-404",
    closeBundle() {
      const distDir = path.resolve(__dirname, "dist");
      const indexPath = path.join(distDir, "index.html");
      const notFoundPath = path.join(distDir, "404.html");
      if (fs.existsSync(indexPath)) {
        fs.copyFileSync(indexPath, notFoundPath);
      }
    },
  };
}

export default defineConfig({
  base: process.env.VITE_BASE_PATH || "/",
  plugins: [react(), spaFallbackPlugin()],
  server: {
    port: 5173,
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: process.env.VITE_API_PROXY_TARGET || "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    target: "esnext",
    cssCodeSplit: true,
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            if (id.includes("@mui/icons-material")) {
              return "vendor-icons";
            }
            if (
              id.includes("@mui/material") ||
              id.includes("@emotion/react") ||
              id.includes("@emotion/styled")
            ) {
              return "vendor-mui";
            }
            if (
              id.includes("react") ||
              id.includes("react-dom") ||
              id.includes("react-router-dom") ||
              id.includes("zustand") ||
              id.includes("axios")
            ) {
              return "vendor-core";
            }
            return "vendor";
          }
        },
      },
    },
  },
});
