import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { Popup } from "./Popup";

// Mock the storage/api imports
vi.mock("@/shared/api", () => ({
  getHistory: vi.fn(() => Promise.resolve(null)),
  clearHistory: vi.fn(() => Promise.resolve(true)),
  healthCheck: vi.fn(() => Promise.resolve(false)),
}));

vi.mock("@/shared/storage", () => ({
  getOnboarded: vi.fn(() => Promise.resolve(true)),
  setOnboarded: vi.fn(),
  getMicMode: vi.fn(() => Promise.resolve("always")),
}));

function mockState(state: string, micMode: string, apiOnline: boolean) {
  (chrome.runtime.sendMessage as ReturnType<typeof vi.fn>).mockImplementation(
    (_msg: unknown, callback?: (response: unknown) => void) => {
      if (callback) {
        callback({ state, micMode, apiOnline });
      }
    }
  );
}

describe("Popup", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockState("idle", "always", true);
  });

  it("renders the Volo header", async () => {
    render(<Popup />);
    await waitFor(() => {
      expect(screen.getByText("Volo")).toBeInTheDocument();
    });
  });

  it("renders Status and History tabs", async () => {
    render(<Popup />);
    await waitFor(() => {
      expect(screen.getByText("Status")).toBeInTheDocument();
      expect(screen.getByText("History")).toBeInTheDocument();
    });
  });

  describe("API offline state", () => {
    beforeEach(() => {
      mockState("listening", "always", false);
    });

    it("shows offline warning when API is down", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText(/Offline/)).toBeInTheDocument();
      });
    });

    it("shows explanation that history is not recording", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText(/history not recording/)).toBeInTheDocument();
      });
    });
  });

  describe("Voice error state", () => {
    beforeEach(() => {
      mockState("error", "always", true);
    });

    it("shows voice unavailable message", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText("Voice unavailable")).toBeInTheDocument();
      });
    });

    it("shows mic permission guidance", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText(/microphone permissions/i)).toBeInTheDocument();
      });
    });
  });

  describe("Mic off state", () => {
    beforeEach(() => {
      mockState("idle", "off", true);
    });

    it("shows Off option as active", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText("Off")).toBeInTheDocument();
      });
    });
  });

  describe("Healthy state", () => {
    it("does not show offline banner when API is online", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText("Volo")).toBeInTheDocument();
      });
      expect(screen.queryByText("API Offline")).not.toBeInTheDocument();
    });

    it("does not show voice error when voice is idle", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText("Volo")).toBeInTheDocument();
      });
      expect(screen.queryByText("Voice unavailable")).not.toBeInTheDocument();
    });

    it("shows listening status when active", async () => {
      mockState("listening", "always", true);
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText(/Listening for/)).toBeInTheDocument();
      });
    });
  });

  describe("Microphone mode display", () => {
    it("renders all three mic options", async () => {
      render(<Popup />);
      await waitFor(() => {
        expect(screen.getByText("Always Listen")).toBeInTheDocument();
        expect(screen.getByText("Click to Talk")).toBeInTheDocument();
        expect(screen.getByText("Off")).toBeInTheDocument();
      });
    });
  });
});
