import type { Config } from "tailwindcss";

export default {
  content: ["./src/**/*.{ts,tsx,html}"],
  theme: {
    extend: {
      colors: {
        volo: {
          bg: "#0a0a0c",
          surface: "rgba(255, 255, 255, 0.03)",
          "surface-2": "#121314",
          input: "rgba(255, 255, 255, 0.03)",
          border: "rgba(255, 255, 255, 0.06)",
          text: "#d0d6e0",
          muted: "#9c9da1",
          faint: "#62666d",
          accent: "#0ea5e9",
          "accent-hover": "#0284c7",
          success: "#22c55e",
          warning: "#f59e0b",
          danger: "#ef4444",
        },
      },
    },
  },
  plugins: [],
} satisfies Config;
