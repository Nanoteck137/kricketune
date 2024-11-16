import { env } from "$env/dynamic/public";
import { ApiClient } from "$lib/api/client";
import type { Handle } from "@sveltejs/kit";

const apiAddress = env.PUBLIC_API_ADDRESS ? env.PUBLIC_API_ADDRESS : "";

export const handle: Handle = async ({ event, resolve }) => {
  const url = new URL(event.request.url);

  let addr = apiAddress;
  if (addr === "") {
    addr = url.origin;
  }

  const client = new ApiClient(addr);
  event.locals.apiClient = client;

  const response = await resolve(event);
  return response;
};
