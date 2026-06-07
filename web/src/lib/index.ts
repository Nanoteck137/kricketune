// place files you want to import through the `$lib` alias in this folder.

import { ApiClient } from "$lib/api/client";
import { clsx, type ClassValue } from "clsx";
import { getContext, setContext } from "svelte";
import toast from "svelte-5-french-toast";
import { twMerge } from "tailwind-merge";

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function handleApiError(err: {
  code: number;
  type: string;
  message: string;
}) {
  toast.error(`API Error: ${err.type} (${err.code}): ${err.message}`);
  console.error("API Error", err);
}

const API_CLIENT_KEY = Symbol("API_CLIENT");

export function setApiClient(baseUrl: string) {
  const apiClient = new ApiClient(baseUrl);
  return setContext(API_CLIENT_KEY, apiClient);
}

export function setApiClientRaw(client: ApiClient) {
  return setContext(API_CLIENT_KEY, client);
}

export function getApiClient() {
  return getContext<ReturnType<typeof setApiClient>>(API_CLIENT_KEY);
}
