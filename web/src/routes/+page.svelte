<script lang="ts">
  import { enhance } from "$app/forms";
  import { invalidate } from "$app/navigation";
  import { GET_STATUS_URL } from "$lib/api/client.js";
  import { Pause, Play, SkipBack, SkipForward } from "lucide-svelte";
  import { onMount } from "svelte";

  const { data } = $props();

  onMount(() => {
    const int = setInterval(() => {
      invalidate(GET_STATUS_URL);
    }, 500);

    return () => {
      clearInterval(int);
    };
  });
</script>

<div class="flex w-full flex-col items-center justify-center gap-6 py-20">
  <div class="flex flex-col items-center justify-center gap-2">
    <p class="line-clamp-1 text-2xl">{data.status.trackName}</p>
    <!-- <p>{data.status.trackAlbum}</p> -->
    <p class="line-clamp-1 text-xl">{data.status.trackArtist}</p>
  </div>

  <div class="flex items-center">
    <form action="?/prev" method="post" use:enhance>
      <button>
        <SkipBack size={50} />
      </button>
    </form>

    {#if data.status.isPlaying}
      <form action="?/pause" method="post" use:enhance>
        <button>
          <Pause size={60} />
        </button>
      </form>
    {:else}
      <form action="?/play" method="post" use:enhance>
        <button>
          <Play size={60} />
        </button>
      </form>
    {/if}

    <form action="?/next" method="post" use:enhance>
      <button>
        <SkipForward size={50} />
      </button>
    </form>
  </div>
</div>
