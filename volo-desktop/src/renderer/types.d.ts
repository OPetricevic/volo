import type { VoloAPI } from "../preload/index";

declare global {
  interface Window {
    volo: VoloAPI;
  }
}
