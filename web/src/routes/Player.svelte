<script lang="ts">
  import { formatTime } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { Status } from "$lib/api/types";
  import {
    FastForward,
    Pause,
    Play,
    Rewind,
    SkipBack,
    SkipForward,
  } from "lucide-svelte";

  type Props = {
    apiClient: ApiClient;
    status?: Status;
  };

  const { apiClient, status }: Props = $props();
</script>

<p class="line-clamp-1 text-center text-lg">
  {status?.currentListName}
</p>

<div class="h-4"></div>

<div class="flex justify-center">
  <img class="aspect-square w-52" src={status?.currentTrack.coverUrl} alt="" />
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
  <!-- <button
    onclick={async () => {
      await apiClient.rewindTrack();
    }}
  >
    <RotateCcw size={32} />
  </button> -->

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

  <!-- <div class="min-h-8 min-w-8"></div> -->
</div>
