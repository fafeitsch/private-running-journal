import {trackEditor, tracks} from "../../wailsjs/go/models";
import {
  DeleteTrack,
  GetTrackTree,
  MoveTrack,
} from "../../wailsjs/go/backend/App";
import {GetTrack, SaveTrack, GetPolylineMeta} from '../../wailsjs/go/trackEditor/TrackEditor';
import TrackDto = trackEditor.TrackDto;
import PolylineMeta = trackEditor.PolylineMeta;
import CoordinateDto = trackEditor.CoordinateDto;

export function useTracksApi() {
  async function getTrackTree(): Promise<tracks.TrackTreeNode> {
    return GetTrackTree();
  }
  async function getTrack(id: string): Promise<TrackDto>  {
    return await GetTrack(id)
  }
  async function getPolylineMeta(coordinates: CoordinateDto[]): Promise<PolylineMeta> {
    return GetPolylineMeta(coordinates);
  }
  function saveTrack(track: trackEditor.SaveTrackDto): Promise<void> {
    return SaveTrack(track);
  }
  function deleteTrack(trackId: string): Promise<void> {
    return DeleteTrack(trackId);
  }
  function moveTrack(trackId: string, newPath: string): Promise<tracks.Track> {
    return MoveTrack(trackId, newPath);
  }
  return { getTrackTree,getTrack,getPolylineMeta, saveTrack,  deleteTrack, moveTrack };
}
