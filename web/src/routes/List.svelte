<script lang="ts">
  import type { ApiClient } from "$lib/api/client";
  import type { List, Status, Track } from "$lib/api/types";
  import { Button } from "@nanoteck137/nano-ui";

  type Props = {
    apiClient: ApiClient;
    lists: List[];
  };

  const { apiClient, lists }: Props = $props();
</script>

<div class="flex flex-col">
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
</div>
