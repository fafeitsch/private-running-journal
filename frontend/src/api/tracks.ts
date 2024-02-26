import { backend } from "../../wailsjs/go/models";
import { GetTracks } from "../../wailsjs/go/backend/App";

export function useTracksApi() {
  function getTracks(): Promise<backend.Track[]> {
    return GetTracks();
  }
  return { getTracks };
}
