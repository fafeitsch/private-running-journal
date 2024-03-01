import {backend, tracks} from "../../wailsjs/go/models";
import { GetTracks, GetGpxData } from "../../wailsjs/go/backend/App";

export function useTracksApi() {
  function getTracks(): Promise<tracks.Track[]> {
    return GetTracks();
  }
  function getGpxData(baseName: string): Promise<tracks.GpxData> {
    return GetGpxData(baseName);
  }
  return { getTracks, getGpxData };
}
