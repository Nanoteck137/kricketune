import { GET_STATUS_URL } from "$lib/api/client";
import { error } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, depends }) => {
  depends(GET_STATUS_URL);

  const res = await locals.apiClient.getStatus();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    status: res.data,
  };
};

export const actions: Actions = {
  play: async ({ locals }) => {
    const res = await locals.apiClient.play();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },

  pause: async ({ locals }) => {
    const res = await locals.apiClient.pause();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },

  next: async ({ locals }) => {
    const res = await locals.apiClient.next();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },

  prev: async ({ locals }) => {
    const res = await locals.apiClient.prev();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },
};
