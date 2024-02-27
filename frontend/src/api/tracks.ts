import { backend } from "../../wailsjs/go/models";
import { GetTracks, GetGpxData } from "../../wailsjs/go/backend/App";

export function useTracksApi() {
  function getTracks(): Promise<backend.Track[]> {
    return GetTracks();
  }
  function getGpxData(baseName: string, variant: string): Promise<backend.GpxData> {
    return GetGpxData(baseName, variant);
  }
  return { getTracks, getGpxData };
}
