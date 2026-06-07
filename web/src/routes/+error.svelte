<script>
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { Button } from "@nanoteck137/nano-ui";
</script>

<div class="py-4 text-center">
  <p class="text-3xl font-semibold">{$page.status}</p>
  <div class="h-2"></div>
  <p class="text-lg">{$page.error?.message}</p>

  {#if $page.status === 400}
    <Button
      onclick={() => {
        localStorage.removeItem("token");
        invalidateAll();
        goto("/");
      }}
    >
      Logout
    </Button>
  {/if}
</div>
