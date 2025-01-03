// DO NOT EDIT THIS: This file was generated by the Pyrin Typescript Generator
import { z } from "zod";

export const Status = z.object({
  trackName: z.string(),
  trackArtist: z.string(),
  trackAlbum: z.string(),
  isPlaying: z.boolean(),
  volume: z.number(),
  mute: z.boolean(),
  queueIndex: z.number(),
  numTracks: z.number(),
  position: z.number(),
  duration: z.number(),
});
export type Status = z.infer<typeof Status>;

export const List = z.object({
  id: z.string(),
  name: z.string(),
});
export type List = z.infer<typeof List>;

export const GetLists = z.object({
  lists: z.array(List),
});
export type GetLists = z.infer<typeof GetLists>;

export const SeekBody = z.object({
  skip: z.number(),
});
export type SeekBody = z.infer<typeof SeekBody>;

