<script lang="ts">
  import { ApiClient } from "$lib/api/client.js";
  import { Status, Track, type List } from "$lib/api/types.js";
  import { onMount } from "svelte";
  import { cn } from "$lib";
  import { z } from "zod";
  import Player from "./Player.svelte";
  import Queue from "./Queue.svelte";
  import ListTab from "./List.svelte";

  const { data } = $props();
  const apiClient = new ApiClient(data.apiAddress);

  let tab = $state<"player" | "list" | "queue" | "snapcast">("player");

  let lists = $state<List[]>([]);
  let queue = $state<Track[]>([]);
  let status = $state<Status | undefined>(data.status);

  async function getLists() {
    const res = await apiClient.getLists();
    if (!res.success) {
      console.error(res.error);
      return;
    }

    lists = res.data.lists;
  }

  async function getQueue() {
    const res = await apiClient.getQueue();
    if (!res.success) {
      console.error(res.error);
      return;
    }

    queue = res.data.tracks;
  }

  onMount(() => {
    getLists();
    getQueue();
  });

  const ConnectedEvent = z.object({});

  const Event = z.discriminatedUnion("type", [
    z.object({
      type: z.literal("connected"),
      data: ConnectedEvent,
    }),
  ]);

  onMount(() => {
    console.log("OnMount");
    const eventSource = new EventSource(apiClient.url.sseHandler());

    eventSource.onmessage = (e) => {
      const event = Event.parse(JSON.parse(e.data));
      console.log(event);

      switch (event.type) {
        case "connected":
          console.log("Connected to SSE");
          break;
      }
    };

    eventSource.addEventListener("connected", (e) => {
      console.log(e.data);
    });

    const StatusEvent = Status.extend({});

    eventSource.addEventListener("status", (e) => {
      const data = StatusEvent.parse(JSON.parse(e.data));
      status = data;
    });

    eventSource.addEventListener("queueChanged", (e) => {
      console.log("queueChanged", e.data);
      getQueue();
    });

    return () => {
      eventSource.close();
    };
  });

  /*
  onMount(() => {
    // Usage example:
    const sse = new ReconnectingEventSource(
      data.apiAddress + "/api/v1/player/sse",
      {
        reconnectInterval: 3000, // Wait 3 seconds before reconnecting
        maxReconnectAttempts: 10, // Try up to 10 times (use Infinity for unlimited)
      },
    );

    sse.addEventListener("open", (e: Event) => {
      console.log("Connection opened");
    });

    sse.addEventListener("message", (e: Event | MessageEvent) => {
      if ("data" in e) {
        console.log("Received:", e.data);
      }
    });

    sse.addEventListener("error", (e: Event) => {
      console.log("Connection error", e);
    });

    sse.addEventListener("maxReconnectAttemptsReached", () => {
      console.log("Could not reconnect after maximum attempts");
    });

    // Listen for custom event types
    sse.addEventListener("customEvent", (e: Event | MessageEvent) => {
      if ("data" in e) {
        console.log("Custom event:", e.data);
      }
    });

    // To close the connection:
    return () => {
      sse.close();
    };
  });
  */
</script>

<div class="container mx-auto">
  <div class="h-4"></div>

  <div
    class="flex flex-col items-center justify-center gap-4 border-b sm:flex-row"
  >
    <button
      class={cn(
        "relative w-24 py-2 text-xl",
        tab === "player" ? "border-b border-blue-200 text-blue-200" : "",
      )}
      onclick={() => {
        tab = "player";
      }}
    >
      Player
    </button>

    <button
      class={cn(
        "w-24 py-2 text-xl",
        tab === "list" ? "border-b border-blue-200 text-blue-200" : "",
      )}
      onclick={() => {
        tab = "list";
      }}
    >
      List
    </button>

    <button
      class={cn(
        "w-24 py-2 text-xl",
        tab === "queue" ? "border-b border-blue-200 text-blue-200" : "",
      )}
      onclick={() => {
        tab = "queue";
      }}
    >
      Queue
    </button>

    <button
      class={cn(
        "w-24 py-2 text-xl",
        tab === "snapcast" ? "border-b border-blue-200 text-blue-200" : "",
      )}
      onclick={() => {
        tab = "snapcast";
      }}
    >
      Snapcast
    </button>
  </div>

  <div class="h-16"></div>

  {#if tab === "player"}
    <Player {apiClient} {status} />
  {/if}

  {#if tab === "list"}
    <ListTab {apiClient} {lists} />
  {/if}

  {#if tab === "queue"}
    <Queue {apiClient} {status} {queue} />
  {/if}
</div>
