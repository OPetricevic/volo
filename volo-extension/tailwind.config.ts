import type { Config } from "tailwindcss";

export default {
  content: ["./src/**/*.{ts,tsx,html}"],
  theme: {
    extend: {
      colors: {
        volo: {
          bg: "#1e1f22",
          surface: "#2b2d31",
          "surface-2": "#313338",
          input: "#383a40",
          border: "#3f4147",
          text: "#dbdee1",
          muted: "#949ba4",
          faint: "#6d6f78",
          accent: "#5865f2",
          "accent-hover": "#4752c4",
          success: "#23a559",
          warning: "#f0b232",
          danger: "#da373c",
        },
      },
    },
  },
  plugins: [],
} satisfies Config;
