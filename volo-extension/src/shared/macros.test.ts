import { describe, it, expect, vi, beforeEach } from "vitest";

// Mock chrome.storage.local
const mockStorage: Record<string, unknown> = {};

vi.stubGlobal("chrome", {
  storage: {
    local: {
      get: vi.fn((keys: string | string[]) => {
        if (typeof keys === "string") {
          return Promise.resolve({ [keys]: mockStorage[keys] });
        }
        const result: Record<string, unknown> = {};
        for (const k of keys) result[k] = mockStorage[k];
        return Promise.resolve(result);
      }),
      set: vi.fn((items: Record<string, unknown>) => {
        Object.assign(mockStorage, items);
        return Promise.resolve();
      }),
      remove: vi.fn((key: string) => {
        delete mockStorage[key];
        return Promise.resolve();
      }),
    },
  },
});

// Mock fetch for API calls
vi.stubGlobal("fetch", vi.fn());

import { fetchMacros, getCachedMacros, matchMacro, isMacroCacheStale, type Macro } from "./api";

const SAMPLE_MACROS: Macro[] = [
  {
    id: "1",
    trigger_phrase: "morning routine",
    name: "Morning Routine",
    actions: [
      { type: "navigate", url: "https://gmail.com" },
      { type: "navigate", url: "https://youtube.com" },
    ],
    enabled: true,
  },
  {
    id: "2",
    trigger_phrase: "work mode",
    name: "Work Mode",
    actions: [
      { type: "navigate", url: "https://github.com" },
      { type: "navigate", url: "https://linear.app" },
      { type: "navigate", url: "https://slack.com" },
    ],
    enabled: true,
  },
];

describe("Macros", () => {
  beforeEach(() => {
    // Clear mock storage
    for (const key of Object.keys(mockStorage)) delete mockStorage[key];
    vi.clearAllMocks();
  });

  describe("fetchMacros", () => {
    it("fetches from API and caches result", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockResolvedValueOnce({
        json: () => Promise.resolve({ data: { macros: SAMPLE_MACROS } }),
      });

      const result = await fetchMacros();

      expect(result).toEqual(SAMPLE_MACROS);
      expect(mockStorage["volo_macros_cache"]).toEqual(SAMPLE_MACROS);
      expect(mockStorage["volo_macros_last_fetch"]).toBeDefined();
    });

    it("returns cached macros when API is offline", async () => {
      // Pre-populate cache
      mockStorage["volo_macros_cache"] = SAMPLE_MACROS;

      (fetch as ReturnType<typeof vi.fn>).mockRejectedValueOnce(new Error("Network error"));

      const result = await fetchMacros();

      expect(result).toEqual(SAMPLE_MACROS);
    });

    it("returns empty array when API offline and no cache", async () => {
      (fetch as ReturnType<typeof vi.fn>).mockRejectedValueOnce(new Error("Network error"));

      const result = await fetchMacros();

      expect(result).toEqual([]);
    });
  });

  describe("getCachedMacros", () => {
    it("returns cached macros", async () => {
      mockStorage["volo_macros_cache"] = SAMPLE_MACROS;

      const result = await getCachedMacros();

      expect(result).toEqual(SAMPLE_MACROS);
    });

    it("returns empty array if no cache", async () => {
      const result = await getCachedMacros();

      expect(result).toEqual([]);
    });
  });

  describe("isMacroCacheStale", () => {
    it("returns true if never fetched", async () => {
      const stale = await isMacroCacheStale();

      expect(stale).toBe(true);
    });

    it("returns false if fetched recently", async () => {
      mockStorage["volo_macros_last_fetch"] = Date.now();

      const stale = await isMacroCacheStale();

      expect(stale).toBe(false);
    });

    it("returns true if fetched more than 5 minutes ago", async () => {
      mockStorage["volo_macros_last_fetch"] = Date.now() - 6 * 60 * 1000;

      const stale = await isMacroCacheStale();

      expect(stale).toBe(true);
    });
  });

  describe("matchMacro", () => {
    beforeEach(() => {
      mockStorage["volo_macros_cache"] = SAMPLE_MACROS;
    });

    it("matches exact trigger phrase", async () => {
      const result = await matchMacro("morning routine");

      expect(result).not.toBeNull();
      expect(result!.name).toBe("Morning Routine");
    });

    it("matches case-insensitive", async () => {
      const result = await matchMacro("Morning Routine");

      expect(result).not.toBeNull();
      expect(result!.name).toBe("Morning Routine");
    });

    it("matches when trigger phrase is contained in transcript", async () => {
      const result = await matchMacro("start my morning routine please");

      expect(result).not.toBeNull();
      expect(result!.name).toBe("Morning Routine");
    });

    it("returns null for no match", async () => {
      const result = await matchMacro("open youtube");

      expect(result).toBeNull();
    });

    it("returns null with empty cache", async () => {
      mockStorage["volo_macros_cache"] = [];

      const result = await matchMacro("morning routine");

      expect(result).toBeNull();
    });

    it("matches work mode correctly", async () => {
      const result = await matchMacro("work mode");

      expect(result).not.toBeNull();
      expect(result!.actions).toHaveLength(3);
      expect(result!.actions[0].url).toBe("https://github.com");
    });
  });
});
