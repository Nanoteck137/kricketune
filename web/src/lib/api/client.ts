import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createUrl, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  url: ClientUrls;

  constructor(baseUrl: string) {
    super(baseUrl);
    this.url = new ClientUrls(baseUrl);
  }
  
  clearQueue(options?: ExtraOptions) {
    return this.request("/api/v1/player/clearQueue", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  getLists(options?: ExtraOptions) {
    return this.request("/api/v1/player/lists", "GET", api.GetLists, z.any(), undefined, options)
  }
  
  getQueue(options?: ExtraOptions) {
    return this.request("/api/v1/player/queue", "GET", api.GetQueue, z.any(), undefined, options)
  }
  
  getStatus(options?: ExtraOptions) {
    return this.request("/api/v1/player/status", "GET", api.Status, z.any(), undefined, options)
  }
  
  loadList(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/player/lists/${id}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  next(options?: ExtraOptions) {
    return this.request("/api/v1/player/next", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  pause(options?: ExtraOptions) {
    return this.request("/api/v1/player/pause", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  play(options?: ExtraOptions) {
    return this.request("/api/v1/player/play", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  prev(options?: ExtraOptions) {
    return this.request("/api/v1/player/prev", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  refreshList(options?: ExtraOptions) {
    return this.request("/api/v1/player/lists/refresh", "POST", api.GetLists, z.any(), undefined, options)
  }
  
  rewindTrack(options?: ExtraOptions) {
    return this.request("/api/v1/player/rewindTrack", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  seek(body: api.SeekBody, options?: ExtraOptions) {
    return this.request("/api/v1/player/seek", "POST", z.undefined(), z.any(), body, options)
  }
  
  setQueueIndex(body: api.SetQueueIndexBody, options?: ExtraOptions) {
    return this.request("/api/v1/player/queueIndex", "POST", z.undefined(), z.any(), body, options)
  }
  
}

export class ClientUrls {
  baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }
  
  clearQueue() {
    return createUrl(this.baseUrl, "/api/v1/player/clearQueue")
  }
  
  getLists() {
    return createUrl(this.baseUrl, "/api/v1/player/lists")
  }
  
  getQueue() {
    return createUrl(this.baseUrl, "/api/v1/player/queue")
  }
  
  getStatus() {
    return createUrl(this.baseUrl, "/api/v1/player/status")
  }
  
  loadList(id: string) {
    return createUrl(this.baseUrl, `/api/v1/player/lists/${id}`)
  }
  
  next() {
    return createUrl(this.baseUrl, "/api/v1/player/next")
  }
  
  pause() {
    return createUrl(this.baseUrl, "/api/v1/player/pause")
  }
  
  play() {
    return createUrl(this.baseUrl, "/api/v1/player/play")
  }
  
  prev() {
    return createUrl(this.baseUrl, "/api/v1/player/prev")
  }
  
  refreshList() {
    return createUrl(this.baseUrl, "/api/v1/player/lists/refresh")
  }
  
  rewindTrack() {
    return createUrl(this.baseUrl, "/api/v1/player/rewindTrack")
  }
  
  seek() {
    return createUrl(this.baseUrl, "/api/v1/player/seek")
  }
  
  setQueueIndex() {
    return createUrl(this.baseUrl, "/api/v1/player/queueIndex")
  }
  
  sseHandler() {
    return createUrl(this.baseUrl, "/api/v1/player/sse")
  }
}
