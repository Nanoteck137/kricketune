import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getStatus(options?: ExtraOptions) {
    return this.request("/api/v1/player/status", "GET", api.Status, z.any(), undefined, options)
  }
  
  getLists(options?: ExtraOptions) {
    return this.request("/api/v1/player/lists", "GET", api.GetLists, z.any(), undefined, options)
  }
  
  loadList(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/player/lists/${id}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  play(options?: ExtraOptions) {
    return this.request("/api/v1/player/play", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  pause(options?: ExtraOptions) {
    return this.request("/api/v1/player/pause", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  next(options?: ExtraOptions) {
    return this.request("/api/v1/player/next", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  prev(options?: ExtraOptions) {
    return this.request("/api/v1/player/prev", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  rewindTrack(options?: ExtraOptions) {
    return this.request("/api/v1/player/rewindTrack", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  seek(body: api.SeekBody, options?: ExtraOptions) {
    return this.request("/api/v1/player/seek", "POST", z.undefined(), z.any(), body, options)
  }
  
  clearQueue(options?: ExtraOptions) {
    return this.request("/api/v1/player/clearQueue", "POST", z.undefined(), z.any(), undefined, options)
  }
}
