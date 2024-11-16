// See https://svelte.dev/docs/kit/types#app.d.ts

import type { ApiClient } from "$lib/api/client";

// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      apiClient: ApiClient;
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {};
