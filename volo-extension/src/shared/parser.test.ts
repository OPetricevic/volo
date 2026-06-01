import { describe, it, expect } from "vitest";
import {
  parseCommand,
  buildExecutionUrl,
  containsWakeWord,
  extractAfterWakeWord,
} from "./parser";

// ─── Wake Word Detection ────────────────────────────────────────────────────

describe("containsWakeWord", () => {
  it("detects 'hey volo' at the start", () => {
    expect(containsWakeWord("hey volo search something")).toBe(true);
  });

  it("detects 'hey volo' in the middle of speech", () => {
    expect(containsWakeWord("um hey volo open youtube")).toBe(true);
  });

  it("is case-insensitive", () => {
    expect(containsWakeWord("Hey Volo search react")).toBe(true);
    expect(containsWakeWord("HEY VOLO open github")).toBe(true);
  });

  it("returns false when wake word is absent", () => {
    expect(containsWakeWord("search for react hooks")).toBe(false);
    expect(containsWakeWord("hey siri open youtube")).toBe(false);
    expect(containsWakeWord("hello")).toBe(false);
  });

  it("returns false for partial matches", () => {
    expect(containsWakeWord("hey vol")).toBe(false);
    expect(containsWakeWord("volo")).toBe(false);
  });
});

describe("extractAfterWakeWord", () => {
  it("extracts command after wake word", () => {
    expect(extractAfterWakeWord("hey volo search react hooks")).toBe("search react hooks");
  });

  it("handles extra whitespace", () => {
    expect(extractAfterWakeWord("  hey volo   open youtube  ")).toBe("open youtube");
  });

  it("returns empty string when nothing follows wake word", () => {
    expect(extractAfterWakeWord("hey volo")).toBe("");
  });

  it("returns null when wake word is not present", () => {
    expect(extractAfterWakeWord("open youtube")).toBeNull();
  });

  it("handles wake word in the middle of speech", () => {
    expect(extractAfterWakeWord("um hey volo search dogs")).toBe("search dogs");
  });
});

// ─── Command Parsing ────────────────────────────────────────────────────────

describe("parseCommand", () => {
  describe("search commands", () => {
    it("parses 'search [query]'", () => {
      const result = parseCommand("search react hooks tutorial");
      expect(result.action).toBe("search");
      expect(result.query).toBe("react hooks tutorial");
      expect(result.confidence).toBeGreaterThanOrEqual(0.9);
    });

    it("parses 'search for [query]'", () => {
      const result = parseCommand("search for best restaurants nearby");
      expect(result.action).toBe("search");
      expect(result.query).toBe("best restaurants nearby");
    });

    it("parses 'look up [query]'", () => {
      const result = parseCommand("look up weather tomorrow");
      expect(result.action).toBe("search");
      expect(result.query).toBe("weather tomorrow");
    });

    it("parses 'find [query]'", () => {
      const result = parseCommand("find golang concurrency patterns");
      expect(result.action).toBe("search");
      expect(result.query).toBe("golang concurrency patterns");
    });

    it("falls back to search for unrecognized commands", () => {
      const result = parseCommand("what is the meaning of life");
      expect(result.action).toBe("search");
      expect(result.query).toBe("what is the meaning of life");
      expect(result.confidence).toBeLessThan(0.8);
    });
  });

  describe("navigate commands", () => {
    it("parses 'open youtube'", () => {
      const result = parseCommand("open youtube");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("youtube");
      expect(result.confidence).toBeGreaterThanOrEqual(0.9);
    });

    it("parses 'go to github'", () => {
      const result = parseCommand("go to github");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("github");
    });

    it("parses 'goto reddit'", () => {
      const result = parseCommand("goto reddit");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("reddit");
    });

    it("handles fuzzy site names ('tube' → youtube)", () => {
      const result = parseCommand("open tube");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("youtube");
    });

    it("handles raw URLs", () => {
      const result = parseCommand("open example.com");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("https://example.com");
    });
  });

  describe("open-and-search commands", () => {
    it("parses 'open youtube [query]'", () => {
      const result = parseCommand("open youtube lofi beats");
      expect(result.action).toBe("open-and-search");
      expect(result.target).toBe("youtube");
      expect(result.query).toBe("lofi beats");
    });

    it("parses 'open spotify jazz'", () => {
      const result = parseCommand("open spotify jazz");
      expect(result.action).toBe("open-and-search");
      expect(result.target).toBe("spotify");
      expect(result.query).toBe("jazz");
    });

    it("parses 'open github react'", () => {
      const result = parseCommand("open github react");
      expect(result.action).toBe("open-and-search");
      expect(result.target).toBe("github");
      expect(result.query).toBe("react");
    });
  });

  describe("open-and-play commands", () => {
    it("parses 'open youtube play lofi hip hop'", () => {
      const result = parseCommand("open youtube play lofi hip hop");
      expect(result.action).toBe("open-and-play");
      expect(result.target).toBe("youtube");
      expect(result.query).toBe("lofi hip hop");
    });

    it("parses 'open spotify play jazz'", () => {
      const result = parseCommand("open spotify play jazz");
      expect(result.action).toBe("open-and-play");
      expect(result.target).toBe("spotify");
      expect(result.query).toBe("jazz");
    });
  });

  describe("browser control commands", () => {
    it("parses 'go back'", () => {
      const result = parseCommand("go back");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("back");
    });

    it("parses 'back'", () => {
      const result = parseCommand("back");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("back");
    });

    it("parses 'go forward'", () => {
      const result = parseCommand("go forward");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("forward");
    });

    it("parses 'close tab'", () => {
      const result = parseCommand("close tab");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("close-tab");
    });

    it("parses 'close this tab'", () => {
      const result = parseCommand("close this tab");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("close-tab");
    });

    it("parses 'new tab'", () => {
      const result = parseCommand("new tab");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("new-tab");
    });

    it("parses 'scroll down'", () => {
      const result = parseCommand("scroll down");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("scroll-down");
    });

    it("parses 'scroll up'", () => {
      const result = parseCommand("scroll up");
      expect(result.action).toBe("browser-control");
      expect(result.query).toBe("scroll-up");
    });
  });

  describe("case insensitivity", () => {
    it("handles uppercase input", () => {
      const result = parseCommand("OPEN YOUTUBE");
      expect(result.action).toBe("navigate");
      expect(result.target).toBe("youtube");
    });

    it("handles mixed case", () => {
      const result = parseCommand("Search For React Hooks");
      expect(result.action).toBe("search");
      expect(result.query).toBe("react hooks");
    });
  });
});

// ─── URL Building ───────────────────────────────────────────────────────────

describe("buildExecutionUrl", () => {
  it("builds Google search URL", () => {
    const url = buildExecutionUrl({ action: "search", query: "react hooks", confidence: 0.9 });
    expect(url).toBe("https://www.google.com/search?q=react%20hooks");
  });

  it("builds YouTube navigation URL", () => {
    const url = buildExecutionUrl({ action: "navigate", target: "youtube", confidence: 0.9 });
    expect(url).toBe("https://www.youtube.com");
  });

  it("builds YouTube search URL for open-and-search", () => {
    const url = buildExecutionUrl({
      action: "open-and-search",
      target: "youtube",
      query: "lofi beats",
      confidence: 0.9,
    });
    expect(url).toBe("https://www.youtube.com/results?search_query=lofi%20beats");
  });

  it("builds Spotify search URL", () => {
    const url = buildExecutionUrl({
      action: "open-and-search",
      target: "spotify",
      query: "jazz",
      confidence: 0.9,
    });
    expect(url).toBe("https://open.spotify.com/search/jazz");
  });

  it("returns null for browser-control actions", () => {
    const url = buildExecutionUrl({ action: "browser-control", query: "back", confidence: 0.95 });
    expect(url).toBeNull();
  });

  it("handles raw URL targets", () => {
    const url = buildExecutionUrl({
      action: "navigate",
      target: "https://example.com",
      confidence: 0.85,
    });
    expect(url).toBe("https://example.com");
  });
});
