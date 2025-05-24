import {projection, trackEditor} from "../../wailsjs/go/models";
import { GetTrackTree } from "../../wailsjs/go/backend/App";
import { DeleteTrack, GetPolylineMeta, GetTrack, SaveTrack } from "../../wailsjs/go/trackEditor/TrackEditor";
import TrackDto = trackEditor.TrackDto;
import PolylineMeta = trackEditor.PolylineMeta;
import CoordinateDto = trackEditor.CoordinateDto;
import TrackTreeNode = projection.TrackTreeNode;

export function useTracksApi() {
  async function getTrackTree(name: string = ""): Promise<TrackTreeNode> {
    return GetTrackTree(name);
  }
  async function getTrack(id: string): Promise<TrackDto> {
    return await GetTrack(id);
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
  return { getTrackTree, getTrack, getPolylineMeta, saveTrack, deleteTrack };
}
