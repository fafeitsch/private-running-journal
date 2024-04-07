import { tracks } from "../../wailsjs/go/models";
import {
  ComputePolylineProps as computePolylineProps,
  CreateNewTrack, DeleteTrack,
  GetGpxData,
  GetTracks, MoveTrack,
  SaveTrack,
} from "../../wailsjs/go/backend/App";
import PolylineProps = tracks.PolylineProps;
import CreateTrack = tracks.CreateTrack;

export function useTracksApi() {
  async function getTracks(): Promise<tracks.Track[]> {
    return GetTracks();
  }
  function getGpxData(baseName: string): Promise<tracks.GpxData> {
    return GetGpxData(baseName);
  }
  function ComputePolylineProps(coordinates: tracks.Coordinates[]): Promise<PolylineProps> {
    return computePolylineProps(coordinates);
  }
  function saveTrack(track: tracks.SaveTrack): Promise<tracks.Track> {
    return SaveTrack(track);
  }
  function createTrack(track: CreateTrack): Promise<tracks.Track> {
    return CreateNewTrack(track);
  }
  function deleteTrack(trackId: string): Promise<void> {
    return DeleteTrack(trackId);
  }
  function moveTrack(trackId: string, newPath: string): Promise<tracks.Track> {
    return MoveTrack(trackId, newPath);
  }
  return { getTracks, getGpxData, ComputePolylineProps, saveTrack, createTrack, deleteTrack, moveTrack };
}
