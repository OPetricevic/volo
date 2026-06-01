/**
 * App configuration.
 * API_URL defaults to production. Override via environment variable for development.
 */

// In production builds, this is the deployed API URL.
// In development, electron-vite injects VITE_API_URL from .env
const DEFAULT_API_URL = "https://api.volo.yourdomain.com";

export const config = {
  apiUrl: import.meta.env.VITE_API_URL || DEFAULT_API_URL,
  ollamaUrl: "http://localhost:11434",
  appVersion: "0.1.0",
} as const;

export type Config = typeof config;
