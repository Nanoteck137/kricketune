<script lang="ts">
  import { setApiClientRaw } from "$lib";
  import { Toaster } from "svelte-5-french-toast";

  let { children, data } = $props();

  setApiClientRaw(data.apiClient);
</script>

<Toaster position="bottom-right" />

{@render children()}
