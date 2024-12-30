import type { Status } from "$lib/api/types";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  let status: Status | undefined = undefined;
  const res = await locals.apiClient.getStatus();
  if (res.success) {
    status = res.data;
  }

  return {
    status,
  };
};
