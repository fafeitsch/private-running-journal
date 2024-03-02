import { tracks } from "../../wailsjs/go/models";
import { GetGpxData, GetTracks } from "../../wailsjs/go/backend/App";

export function useTracksApi() {
  function getTracks(): Promise<tracks.Track[]> {
    return GetTracks();
  }
  function getGpxData(baseName: string): Promise<tracks.GpxData> {
    return GetGpxData(baseName);
  }
  return { getTracks, getGpxData };
}
