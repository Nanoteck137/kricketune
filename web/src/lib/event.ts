export interface ReconnectingEventSourceOptions extends EventSourceInit {
  reconnectInterval?: number;
  maxReconnectAttempts?: number;
}

export type EventCallback = (event: Event | MessageEvent) => void;

export class ReconnectingEventSource {
  private url: string;
  private options: ReconnectingEventSourceOptions;
  private reconnectInterval: number;
  private maxReconnectAttempts: number;
  private reconnectAttempts: number;
  private eventSource: EventSource | null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null;
  private listeners: Record<string, EventCallback[]>;

  constructor(url: string, options: ReconnectingEventSourceOptions = {}) {
    this.url = url;
    this.options = options;
    this.reconnectInterval = options.reconnectInterval || 3000;
    this.maxReconnectAttempts = options.maxReconnectAttempts || Infinity;
    this.reconnectAttempts = 0;
    this.eventSource = null;
    this.reconnectTimer = null;
    this.listeners = {};

    this.connect();
  }

  private connect(): void {
    try {
      this.eventSource = new EventSource(this.url, this.options);

      this.eventSource.onopen = (e: Event) => {
        console.log("EventSource connected");
        this.reconnectAttempts = 0;
        this.emit("open", e);
      };

      this.eventSource.onmessage = (e: MessageEvent) => {
        this.emit("message", e);
      };

      this.eventSource.onerror = (e: Event) => {
        console.error("EventSource error");
        this.emit("error", e);
        this.handleDisconnect();
      };
    } catch (error) {
      console.error("Failed to create EventSource:", error);
      this.handleDisconnect();
    }
  }

  private handleDisconnect(): void {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }

    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log("Max reconnect attempts reached");
      this.emit(
        "maxReconnectAttemptsReached",
        new Event("maxReconnectAttemptsReached"),
      );
      return;
    }

    this.reconnectAttempts++;
    console.log(`Reconnecting... Attempt ${this.reconnectAttempts}`);

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, this.reconnectInterval);
  }

  public addEventListener(event: string, callback: EventCallback): void {
    if (!this.listeners[event]) {
      this.listeners[event] = [];
    }
    this.listeners[event].push(callback);

    // Also add to the actual EventSource if it exists
    if (
      this.eventSource &&
      event !== "open" &&
      event !== "error" &&
      event !== "maxReconnectAttemptsReached"
    ) {
      this.eventSource.addEventListener(event, callback as EventListener);
    }
  }

  public removeEventListener(event: string, callback: EventCallback): void {
    if (this.listeners[event]) {
      this.listeners[event] = this.listeners[event].filter(
        (cb) => cb !== callback,
      );
    }
    if (this.eventSource) {
      this.eventSource.removeEventListener(event, callback as EventListener);
    }
  }

  private emit(event: string, data: Event | MessageEvent): void {
    if (this.listeners[event]) {
      this.listeners[event].forEach((callback) => callback(data));
    }
  }

  public close(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
    }
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
    this.reconnectAttempts = this.maxReconnectAttempts;
  }
}
