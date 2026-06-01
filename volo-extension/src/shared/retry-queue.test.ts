import { describe, it, expect, beforeEach, vi } from "vitest";

// Mock chrome.storage.local
const storage: Record<string, unknown> = {};
const chromeMock = {
  storage: {
    local: {
      get: vi.fn(async (key: string) => ({ [key]: storage[key] })),
      set: vi.fn(async (data: Record<string, unknown>) => {
        Object.assign(storage, data);
      }),
    },
  },
};
vi.stubGlobal("chrome", chromeMock);
vi.stubGlobal("crypto", { randomUUID: () => Math.random().toString(36).slice(2) });

// Import after mocking
import { enqueue, flush, getQueueSize, clearQueue } from "./retry-queue";

beforeEach(() => {
  // Clear storage between tests
  for (const key of Object.keys(storage)) delete storage[key];
});

describe("retry-queue", () => {
  describe("enqueue", () => {
    it("adds a command to the queue", async () => {
      await enqueue("search react hooks", "https://google.com");
      const size = await getQueueSize();
      expect(size).toBe(1);
    });

    it("stores transcript and url", async () => {
      await enqueue("open youtube", "https://example.com");
      const result = await chromeMock.storage.local.get("volo_retry_queue");
      const queue = result["volo_retry_queue"] as Array<{ transcript: string; currentUrl: string }>;
      expect(queue[0].transcript).toBe("open youtube");
      expect(queue[0].currentUrl).toBe("https://example.com");
    });

    it("caps at 50 items (drops oldest)", async () => {
      for (let i = 0; i < 55; i++) {
        await enqueue(`command ${i}`);
      }
      const size = await getQueueSize();
      expect(size).toBe(50);

      // First item should be "command 5" (0-4 were dropped)
      const result = await chromeMock.storage.local.get("volo_retry_queue");
      const queue = result["volo_retry_queue"] as Array<{ transcript: string }>;
      expect(queue[0].transcript).toBe("command 5");
    });
  });

  describe("flush", () => {
    it("sends all queued commands", async () => {
      await enqueue("cmd 1");
      await enqueue("cmd 2");
      await enqueue("cmd 3");

      const sendFn = vi.fn(async () => true);
      const sent = await flush(sendFn);

      expect(sent).toBe(3);
      expect(sendFn).toHaveBeenCalledTimes(3);
      expect(await getQueueSize()).toBe(0);
    });

    it("keeps failed items in queue with incremented retries", async () => {
      await enqueue("will fail");

      const sendFn = vi.fn(async () => false);
      const sent = await flush(sendFn);

      expect(sent).toBe(0);
      expect(await getQueueSize()).toBe(1);

      // Check retries incremented
      const result = await chromeMock.storage.local.get("volo_retry_queue");
      const queue = result["volo_retry_queue"] as Array<{ retries: number }>;
      expect(queue[0].retries).toBe(1);
    });

    it("drops items after 3 failed retries", async () => {
      await enqueue("persistent failure");

      const sendFn = vi.fn(async () => false);

      // Flush 3 times — each time retries increments
      await flush(sendFn);
      await flush(sendFn);
      await flush(sendFn);

      // After 3 failures, item should be dropped
      expect(await getQueueSize()).toBe(0);
    });

    it("handles empty queue gracefully", async () => {
      const sendFn = vi.fn(async () => true);
      const sent = await flush(sendFn);
      expect(sent).toBe(0);
      expect(sendFn).not.toHaveBeenCalled();
    });

    it("partially succeeds — keeps only failures", async () => {
      await enqueue("success 1");
      await enqueue("fail 1");
      await enqueue("success 2");

      let callCount = 0;
      const sendFn = vi.fn(async () => {
        callCount++;
        return callCount !== 2; // second call fails
      });

      const sent = await flush(sendFn);
      expect(sent).toBe(2);
      expect(await getQueueSize()).toBe(1);
    });
  });

  describe("clearQueue", () => {
    it("empties the queue", async () => {
      await enqueue("cmd 1");
      await enqueue("cmd 2");
      expect(await getQueueSize()).toBe(2);

      await clearQueue();
      expect(await getQueueSize()).toBe(0);
    });
  });
});
