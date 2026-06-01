import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        volo: {
          bg: "#0a0a0f",
          surface: "#141419",
          "surface-2": "#1c1c24",
          border: "#2a2a35",
          text: "#f0f0f5",
          muted: "#8888a0",
          accent: "#0ea5e9",
          "accent-glow": "rgba(14, 165, 233, 0.15)",
          teal: "#14b8a6",
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [],
} satisfies Config;
