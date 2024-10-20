import {trackEditor, tracks} from "../../wailsjs/go/models";
import {
  ComputePolylineProps as computePolylineProps,
  CreateNewTrack,
  DeleteTrack,
  GetTrackTree,
  MoveTrack,
  SaveTrack,
} from "../../wailsjs/go/backend/App";
import {GetTrack} from '../../wailsjs/go/trackEditor/TrackEditor';
import PolylineProps = tracks.PolylineProps;
import CreateTrack = tracks.CreateTrack;
import TrackDto = trackEditor.TrackDto;

export function useTracksApi() {
  async function getTrackTree(): Promise<tracks.TrackTreeNode> {
    return GetTrackTree();
  }
  async function getTrack(id: string): Promise<TrackDto>  {
    return await GetTrack(id)
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
  return { getTrackTree,getTrack, ComputePolylineProps, saveTrack, createTrack, deleteTrack, moveTrack };
}
