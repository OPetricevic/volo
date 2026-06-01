/**
 * Volo API Client
 *
 * Handles all communication with the Go backend.
 * Falls back to local-only mode if the API is unreachable.
 */

const API_BASE = "http://localhost:8080/api/v1";

// Storage keys
const TOKEN_KEY = "volo_auth_token";
const DEVICE_ID_KEY = "volo_device_id";

interface ApiResponse<T> {
  data?: T;
  error?: {
    code: string;
    message: string;
    internal_error?: string;
  };
}

interface AuthResponse {
  token: string;
  user: { id: string; role: string };
  device: { id: string; device_id: string };
}

interface CommandResponse {
  action: string;
  target?: string;
  query?: string;
  confidence: number;
  suggestions?: string[];
  execute?: { url: string; auto_play?: boolean };
}

interface HistoryResponse {
  commands: Array<{
    id: string;
    transcript: string;
    parsed_action: string;
    parsed_target?: string;
    parsed_query?: string;
    confidence: number;
    executed_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
}

// ─── Device ID ───────────────────────────────────────────

function generateDeviceId(): string {
  return "ext-" + crypto.randomUUID();
}

async function getDeviceId(): Promise<string> {
  const result = await chrome.storage.local.get(DEVICE_ID_KEY);
  if (result[DEVICE_ID_KEY]) return result[DEVICE_ID_KEY];

  const id = generateDeviceId();
  await chrome.storage.local.set({ [DEVICE_ID_KEY]: id });
  return id;
}

// ─── Token Management ────────────────────────────────────

async function getToken(): Promise<string | null> {
  const result = await chrome.storage.local.get(TOKEN_KEY);
  return result[TOKEN_KEY] || null;
}

async function setToken(token: string): Promise<void> {
  await chrome.storage.local.set({ [TOKEN_KEY]: token });
}

async function clearToken(): Promise<void> {
  await chrome.storage.local.remove(TOKEN_KEY);
}

// ─── HTTP Helpers ────────────────────────────────────────

async function request<T>(
  method: string,
  path: string,
  body?: unknown,
  auth = true
): Promise<ApiResponse<T>> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (auth) {
    const token = await getToken();
    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }
  }

  try {
    const response = await fetch(`${API_BASE}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    const data: ApiResponse<T> = await response.json();
    return data;
  } catch {
    // Network error — API unreachable
    return {
      error: {
        code: "NETWORK_ERROR",
        message: "Cannot reach Volo API. Running in offline mode.",
      },
    };
  }
}

// ─── Public API ──────────────────────────────────────────

/**
 * Register this device with the backend.
 * Called once on first install. Returns a JWT token.
 */
export async function registerDevice(): Promise<boolean> {
  const deviceId = await getDeviceId();

  const platform = "extension";
  const deviceName = `Chrome on ${navigator.platform}`;

  const result = await request<AuthResponse>("POST", "/auth/device", {
    device_id: deviceId,
    device_name: deviceName,
    platform,
  }, false);

  if (result.data?.token) {
    await setToken(result.data.token);
    return true;
  }

  console.warn("[Volo API] Registration failed:", result.error);
  return false;
}

/**
 * Check if we have a valid token stored.
 */
export async function isAuthenticated(): Promise<boolean> {
  const token = await getToken();
  return token !== null;
}

/**
 * Send a voice command to the backend for processing + history storage.
 * Returns the parsed command, or null if API is unreachable.
 */
export async function sendCommand(
  transcript: string,
  currentUrl?: string
): Promise<CommandResponse | null> {
  const result = await request<CommandResponse>("POST", "/command", {
    transcript,
    context: {
      current_url: currentUrl || "",
      timestamp: new Date().toISOString(),
    },
  });

  if (result.data) return result.data;

  // API error or offline — return null (caller falls back to local parsing)
  if (result.error?.code !== "NETWORK_ERROR") {
    console.warn("[Volo API] Command failed:", result.error);
  }
  return null;
}

/**
 * Get command history.
 */
export async function getHistory(
  page = 1,
  pageSize = 20
): Promise<HistoryResponse | null> {
  const result = await request<HistoryResponse>(
    "GET",
    `/history?page=${page}&page_size=${pageSize}`
  );
  return result.data || null;
}

/**
 * Clear all command history.
 */
export async function clearHistory(): Promise<boolean> {
  const result = await request<{ status: string }>("DELETE", "/history");
  return result.data?.status === "cleared";
}

/**
 * Logout (revoke current session).
 */
export async function logout(): Promise<void> {
  await request("POST", "/auth/logout");
  await clearToken();
}

/**
 * Check if the API is reachable.
 */
export async function healthCheck(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE.replace("/api/v1", "")}/health`, {
      method: "GET",
      signal: AbortSignal.timeout(3000),
    });
    return response.ok;
  } catch {
    return false;
  }
}
