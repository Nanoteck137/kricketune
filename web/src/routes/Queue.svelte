<script lang="ts">
  import type { ApiClient } from "$lib/api/client";
  import type { Status, Track } from "$lib/api/types";
  import { ListX } from "lucide-svelte";

  type Props = {
    apiClient: ApiClient;
    status?: Status;
    queue: Track[];
  };

  const { apiClient, status, queue }: Props = $props();
</script>

<button
  onclick={async () => {
    await apiClient.clearQueue();
  }}
>
  <ListX size={42} />
</button>

<div class="flex flex-col">
  <p>{(status?.queueIndex ?? 0) + 1} / {status?.numTracks}</p>
  {#each queue as t, i}
    <p class={i === status?.queueIndex ? "text-red-200" : ""}>
      {i + 1} - {t.name}
    </p>
  {/each}
</div>
