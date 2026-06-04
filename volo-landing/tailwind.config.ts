import type { Config } from "tailwindcss";

export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        volo: {
          bg: "#060609",
          surface: "rgba(255, 255, 255, 0.03)",
          "surface-2": "rgba(255, 255, 255, 0.05)",
          border: "rgba(255, 255, 255, 0.08)",
          "border-strong": "rgba(255, 255, 255, 0.12)",
          text: "#f0f0f5",
          muted: "#7a7a95",
          accent: "#0ea5e9",
          "accent-glow": "rgba(14, 165, 233, 0.15)",
          teal: "#14b8a6",
        },
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      boxShadow: {
        'glass': '0 20px 40px -12px rgba(0, 0, 0, 0.5), inset 0 0 0 1px rgba(255, 255, 255, 0.02)',
        'glass-lg': '0 25px 50px -12px rgba(0, 0, 0, 0.6), inset 0 0 0 1px rgba(255, 255, 255, 0.04)',
        'glow': '0 0 40px -12px rgba(14, 165, 233, 0.3)',
        'glow-sm': '0 0 20px -6px rgba(14, 165, 233, 0.2)',
      },
    },
  },
  plugins: [],
} satisfies Config;
