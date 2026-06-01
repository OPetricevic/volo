/**
 * Retry Queue
 *
 * Stores failed API calls in chrome.storage.local and retries them
 * when the API becomes reachable again.
 *
 * - Max 50 queued commands (oldest dropped if exceeded)
 * - Retries on health check success
 * - Each item retried max 3 times before being dropped
 */

const QUEUE_KEY = "volo_retry_queue";
const MAX_QUEUE_SIZE = 50;
const MAX_RETRIES = 3;

interface QueuedCommand {
  id: string;
  transcript: string;
  currentUrl?: string;
  timestamp: string;
  retries: number;
}

/**
 * Add a failed command to the retry queue.
 */
export async function enqueue(transcript: string, currentUrl?: string): Promise<void> {
  const queue = await getQueue();

  const item: QueuedCommand = {
    id: crypto.randomUUID(),
    transcript,
    currentUrl,
    timestamp: new Date().toISOString(),
    retries: 0,
  };

  queue.push(item);

  // Cap at max size (drop oldest)
  while (queue.length > MAX_QUEUE_SIZE) {
    queue.shift();
  }

  await saveQueue(queue);
  console.log(`[Volo Queue] Enqueued: "${transcript}" (${queue.length} in queue)`);
}

/**
 * Process all queued commands. Called when API comes back online.
 * Returns the number of successfully sent commands.
 */
export async function flush(
  sendFn: (transcript: string, currentUrl?: string) => Promise<boolean>
): Promise<number> {
  const queue = await getQueue();
  if (queue.length === 0) return 0;

  console.log(`[Volo Queue] Flushing ${queue.length} queued commands...`);

  const remaining: QueuedCommand[] = [];
  let sent = 0;

  for (const item of queue) {
    const success = await sendFn(item.transcript, item.currentUrl);
    if (success) {
      sent++;
    } else {
      item.retries++;
      if (item.retries < MAX_RETRIES) {
        remaining.push(item);
      } else {
        console.log(`[Volo Queue] Dropped after ${MAX_RETRIES} retries: "${item.transcript}"`);
      }
    }
  }

  await saveQueue(remaining);
  console.log(`[Volo Queue] Flushed: ${sent} sent, ${remaining.length} remaining`);
  return sent;
}

/**
 * Get the current queue size.
 */
export async function getQueueSize(): Promise<number> {
  const queue = await getQueue();
  return queue.length;
}

/**
 * Clear the entire queue.
 */
export async function clearQueue(): Promise<void> {
  await chrome.storage.local.set({ [QUEUE_KEY]: [] });
}

// ─── Internal ────────────────────────────────────────────

async function getQueue(): Promise<QueuedCommand[]> {
  const result = await chrome.storage.local.get(QUEUE_KEY);
  return result[QUEUE_KEY] || [];
}

async function saveQueue(queue: QueuedCommand[]): Promise<void> {
  await chrome.storage.local.set({ [QUEUE_KEY]: queue });
}
