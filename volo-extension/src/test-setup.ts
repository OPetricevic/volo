import "@testing-library/jest-dom/vitest";
import { vi } from "vitest";

// Mock chrome APIs globally for all tests
vi.stubGlobal("chrome", {
  runtime: {
    sendMessage: vi.fn((_msg, callback) => {
      if (callback) callback(null);
    }),
    onMessage: { addListener: vi.fn() },
    onInstalled: { addListener: vi.fn() },
  },
  storage: {
    local: {
      get: vi.fn(() => Promise.resolve({})),
      set: vi.fn(() => Promise.resolve()),
      remove: vi.fn(() => Promise.resolve()),
    },
  },
  tabs: {
    create: vi.fn(),
    query: vi.fn(() => Promise.resolve([])),
  },
});
