<script lang="ts">
  import { cn } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { List, Status } from "$lib/api/types";

  type Props = {
    apiClient: ApiClient;
    status?: Status;
    lists: List[];
  };

  const { apiClient, status, lists }: Props = $props();
</script>

<div class="flex items-center justify-between border-b pb-4">
  <div class="flex items-center gap-4">
    <h2 class="text-xl font-semibold">Lists</h2>
  </div>
</div>

<div class="h-4"></div>

<div class="flex flex-col">
  <div class="flex flex-col gap-2">
    {#each lists as list}
      <button
        class={cn(
          "flex w-full gap-4 border-b p-3",
          status?.currentListId === list.id
            ? "rounded bg-primary text-primary-foreground"
            : "hover:bg-accent hover:text-accent-foreground",
        )}
        onclick={async () => {
          await apiClient.loadList(list.id);
        }}
      >
        <div class="flex flex-1 flex-col items-start">
          <div class="line-clamp-1 text-start font-medium">{list.name}</div>
          <!-- <div class="text-sm text-gray-400">42 tracks</div> -->
        </div>
      </button>
    {/each}
  </div>
</div>
