import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent }) => {
  const data = await parent();

  const status = await data.apiClient.getStatus();
  if (!status.success) {
    error(status.error.code, { message: status.error.message });
  }

  return {
    ...data,
    status: status.data,
  };
};
