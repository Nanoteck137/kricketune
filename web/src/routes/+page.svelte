<script lang="ts">
  import { ApiClient } from "$lib/api/client.js";
  import type { Status, List } from "$lib/api/types.js";
  import {
    FastForward,
    List as ListIcon,
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

  const { data } = $props();
  const apiClient = new ApiClient(data.apiAddress);

  let lists = $state<List[]>([]);
  let status = $state<Status | undefined>(data.status);

  async function getLists() {
    const res = await apiClient.getLists();
    if (!res.success) {
      console.error(res.error);
      return;
    }

    lists = res.data.lists;
  }

  async function updateStatus() {
    const res = await apiClient.getStatus();
    if (!res.success) {
      console.error(res.error);
      return;
    }

    status = res.data;
  }

  onMount(() => {
    getLists();
    updateStatus();

    const int = setInterval(async () => {
      updateStatus();
    }, 500);

    return () => {
      clearInterval(int);
    };
  });
</script>

<div class="container mx-auto">
  <div class="h-16"></div>

  <div class="">
    <p class="line-clamp-1 text-center text-2xl">
      {status?.trackName ?? "Unknown"}
    </p>
    <p class="line-clamp-1 text-center text-xl">
      {status?.trackArtist ?? "Unknown"}
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

        await updateStatus();
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
</div>
