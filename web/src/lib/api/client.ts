import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";

export const GET_SETS_URL = "/api/v1/player/sets"
export const CHANGE_SET_URL = "/api/v1/player/sets/:index"
export const GET_STATUS_URL = "/api/v1/player/status"
export const PLAY_URL = "/api/v1/player/play"
export const PAUSE_URL = "/api/v1/player/pause"
export const NEXT_URL = "/api/v1/player/next"
export const PREV_URL = "/api/v1/player/prev"
export const CLEAR_QUEUE_URL = "/api/v1/player/clearQueue"

export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getSets(options?: ExtraOptions) {
    return this.request("/api/v1/player/sets", "GET", api.Status, z.undefined(), undefined, options)
  }
  
  changeSet(index: string, options?: ExtraOptions) {
    return this.request(`/api/v1/player/sets/${index}`, "POST", api.Status, z.undefined(), undefined, options)
  }
  
  getStatus(options?: ExtraOptions) {
    return this.request("/api/v1/player/status", "GET", api.Status, z.undefined(), undefined, options)
  }
  
  play(options?: ExtraOptions) {
    return this.request("/api/v1/player/play", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  pause(options?: ExtraOptions) {
    return this.request("/api/v1/player/pause", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  next(options?: ExtraOptions) {
    return this.request("/api/v1/player/next", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  prev(options?: ExtraOptions) {
    return this.request("/api/v1/player/prev", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  clearQueue(options?: ExtraOptions) {
    return this.request("/api/v1/player/clearQueue", "POST", z.undefined(), z.undefined(), undefined, options)
  }
}
