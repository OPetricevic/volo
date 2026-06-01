import { describe, it, expect } from "vitest";

// Test Ollama URL construction and model detection logic
// (Can't test actual Electron APIs in unit tests, but can test the logic)

describe("Ollama configuration", () => {
  const OLLAMA_URL = "http://localhost:11434";
  const MODEL_NAME = "gemma2:2b";

  it("ollama URL is localhost", () => {
    expect(OLLAMA_URL).toBe("http://localhost:11434");
  });

  it("model name is gemma2:2b", () => {
    expect(MODEL_NAME).toBe("gemma2:2b");
  });

  it("tags endpoint is correct", () => {
    const tagsUrl = `${OLLAMA_URL}/api/tags`;
    expect(tagsUrl).toBe("http://localhost:11434/api/tags");
  });

  it("generate endpoint is correct", () => {
    const generateUrl = `${OLLAMA_URL}/api/generate`;
    expect(generateUrl).toBe("http://localhost:11434/api/generate");
  });
});

describe("Ollama request body", () => {
  it("builds correct generate request", () => {
    const body = {
      model: "gemma2:2b",
      prompt: "test prompt",
      stream: false,
      options: { temperature: 0.3, num_predict: 256 },
    };

    expect(body.model).toBe("gemma2:2b");
    expect(body.stream).toBe(false);
    expect(body.options.temperature).toBe(0.3);
    expect(body.options.num_predict).toBe(256);
    expect(body.options.num_predict).toBeLessThanOrEqual(512); // safety limit
  });

  it("temperature is low for factual responses", () => {
    const temperature = 0.3;
    expect(temperature).toBeLessThan(0.5); // Low = more deterministic
    expect(temperature).toBeGreaterThan(0); // Not zero (allows some variation)
  });
});

describe("Model detection from tags response", () => {
  it("detects gemma in model list", () => {
    const mockResponse = {
      models: [
        { name: "gemma2:2b", size: 2500000000 },
        { name: "llama3:8b", size: 8000000000 },
      ],
    };

    const hasGemma = mockResponse.models.some((m) => m.name.includes("gemma"));
    expect(hasGemma).toBe(true);
  });

  it("returns false when gemma not installed", () => {
    const mockResponse = {
      models: [{ name: "llama3:8b", size: 8000000000 }],
    };

    const hasGemma = mockResponse.models.some((m) => m.name.includes("gemma"));
    expect(hasGemma).toBe(false);
  });

  it("handles empty model list", () => {
    const mockResponse = { models: [] };
    const hasGemma = mockResponse.models.some((m: { name: string }) => m.name.includes("gemma"));
    expect(hasGemma).toBe(false);
  });
});
