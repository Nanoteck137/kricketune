import { PUBLIC_API_ADDRESS } from "$env/static/public";
import { ApiClient } from "$lib/api/client";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ url }) => {
  let addr = PUBLIC_API_ADDRESS;
  if (addr === "") {
    addr = url.origin;
  }

  const apiClient = new ApiClient(addr);

  return {
    apiClient,
  };
};
