<script lang="ts">
  import { ApiClient } from "$lib/api/client.js";
  import { Status, Track, type List } from "$lib/api/types.js";
  import {
    FastForward,
    List as ListIcon,
    ListX,
    Pause,
    Play,
    Rewind,
    RotateCcw,
    SkipBack,
    SkipForward,
  } from "lucide-svelte";
  import { onMount } from "svelte";
  import { Button, Sheet } from "@nanoteck137/nano-ui";
  import { formatTime } from "$lib";
  import { z } from "zod";
  import { ReconnectingEventSource } from "$lib/event.js";

  const { data } = $props();
  const apiClient = new ApiClient(data.apiAddress);

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
    const eventSource = new EventSource(
      data.apiAddress + "/api/v1/player/sse",
    );

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
  <div class="h-16"></div>

  <div class="flex justify-center">
    <img
      class="aspect-square w-52"
      src={status?.currentTrack.coverUrl}
      alt=""
    />
  </div>

  <div class="h-4"></div>

  <div class="">
    <p class="line-clamp-1 text-center text-2xl">
      {status?.currentTrack.name ?? "Unknown"}
    </p>
    <p class="line-clamp-1 text-center text-xl">
      {status?.currentTrack.artists.join(", ") ?? "Unknown"}
    </p>
    {#if status}
      <p class="line-clamp-1 text-center text-lg">
        {formatTime(status.position / 1000)} / {formatTime(
          status.duration / 1000,
        )}
      </p>
    {/if}
  </div>

  <div class="h-4"></div>

  <div class="flex items-center justify-center gap-4">
    <button
      onclick={async () => {
        await apiClient.clearQueue();
      }}
    >
      <ListX size={42} />
    </button>

    <button
      onclick={async () => {
        await apiClient.rewindTrack();
      }}
    >
      <RotateCcw size={32} />
    </button>

    <button
      onclick={async () => {
        if (!status) return;

        await apiClient.seek({ skip: -15 });
        status.position -= 15 * 1000;
      }}
    >
      <Rewind size={32} />
    </button>

    <button
      onclick={async () => {
        await apiClient.prev();
      }}
    >
      <SkipBack size={42} />
    </button>

    <button
      onclick={async () => {
        if (!status) return;

        if (status.isPlaying) {
          await apiClient.pause();
        } else {
          await apiClient.play();
        }
      }}
    >
      {#if status?.isPlaying}
        <Pause size={48} />
      {:else}
        <Play size={48} />
      {/if}
    </button>

    <button
      onclick={async () => {
        await apiClient.next();
      }}
    >
      <SkipForward size={42} />
    </button>

    <button
      onclick={async () => {
        if (!status) return;

        await apiClient.seek({ skip: 15 });
        status.position += 15 * 1000;
      }}
    >
      <FastForward size={32} />
    </button>

    <Sheet.Root>
      <Sheet.Trigger>
        <ListIcon size={32} />
      </Sheet.Trigger>
      <Sheet.Content side="right" class="overflow-y-scroll">
        <Sheet.Header>
          <Sheet.Title>Lists Available</Sheet.Title>
        </Sheet.Header>

        <div class="h-4"></div>

        <div class="flex flex-col gap-2">
          {#each lists as list}
            <Button
              class="justify-start"
              variant="ghost"
              onclick={async () => {
                await apiClient.loadList(list.id);
              }}
            >
              {list.name}
            </Button>
          {/each}
        </div>
      </Sheet.Content>
    </Sheet.Root>
  </div>

  <div class="flex flex-col">
    <p>{(status?.queueIndex ?? 0) + 1} / {status?.numTracks}</p>
    {#each queue as t, i}
      <p class={i === status?.queueIndex ? "text-red-200" : ""}>
        {i + 1} - {t.name}
      </p>
    {/each}
  </div>
</div>
