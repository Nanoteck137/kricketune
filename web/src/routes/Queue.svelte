<script lang="ts">
  import { cn } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { Status, Track } from "$lib/api/types";
  import { Button } from "@nanoteck137/nano-ui";
  import { ListX } from "lucide-svelte";

  type Props = {
    apiClient: ApiClient;
    status?: Status;
    queue: Track[];
  };

  const { apiClient, status, queue }: Props = $props();
</script>

<div class="flex items-center justify-between border-b pb-4">
  <div class="flex items-center gap-4">
    <h2 class="text-xl font-semibold">Queue</h2>
    <span class="text-sm text-muted-foreground">
      {(status?.queueIndex ?? 0) + 1} / {status?.numTracks}
    </span>
  </div>

  <Button
    variant="outline"
    onclick={async () => {
      await apiClient.clearQueue();
    }}
  >
    <ListX />
    Clear Queue
  </Button>
</div>

<div class="h-8"></div>

<div class="flex flex-col gap-2">
  {#each queue as t, i}
    <button
      class={cn(
        "flex w-full gap-4 border-b p-3",
        i === status?.queueIndex
          ? "rounded bg-primary text-primary-foreground"
          : "hover:bg-accent hover:text-accent-foreground",
      )}
      onclick={() => {
        apiClient.setQueueIndex({ index: i });
      }}
    >
      <img
        class="aspect-square w-12 rounded object-cover"
        src={t.coverUrl}
        alt=""
        loading="lazy"
      />

      <div class="flex flex-1 flex-col items-start">
        <div class="line-clamp-1 text-start font-medium">
          {#if i === status?.queueIndex}
            &gt;
          {/if}
          {t.name}
        </div>

        <div class="line-clamp-1 text-start text-sm font-medium opacity-60">
          {t.artists.join(", ")}
        </div>
      </div>
    </button>
  {/each}
</div>
