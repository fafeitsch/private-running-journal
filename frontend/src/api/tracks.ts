import {backend} from '../../wailsjs/go/models';

const mockTracks = [
  {
    "variant": "Seeuferweg",
    "length": 16007,
    "id": "seeuferweg",
    "baseId": "waldrunde",
    "baseName": "Waldrunde",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 7864,
    "id": "kurz",
    "baseId": "waldrunde",
    "baseName": "Waldrunde",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 17153,
    "id": "hügelvariante",
    "baseId": "waldrunde",
    "baseName": "Waldrunde",
    "waypointData": ""
  },
  {
    "variant": "Rundkurs",
    "length": 18890,
    "id": "rundkurs",
    "baseId": "waldrunde",
    "baseName": "Waldrunde",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 13528,
    "id": "kurz",
    "baseId": "bergpfad",
    "baseName": "Bergpfad",
    "waypointData": ""
  },
  {
    "variant": "Rundweg",
    "length": 10702,
    "id": "rundweg",
    "baseId": "bergpfad",
    "baseName": "Bergpfad",
    "waypointData": ""
  },
  {
    "variant": "Parkstrecke",
    "length": 16928,
    "id": "parkstrecke",
    "baseId": "bergpfad",
    "baseName": "Bergpfad",
    "waypointData": ""
  },
  {
    "variant": "Rundkurs",
    "length": 8409,
    "id": "rundkurs",
    "baseId": "hügelrunde",
    "baseName": "Hügelrunde",
    "waypointData": ""
  },
  {
    "variant": "Seeuferweg",
    "length": 5065,
    "id": "seeuferweg",
    "baseId": "hügelrunde",
    "baseName": "Hügelrunde",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 11589,
    "id": "kurz",
    "baseId": "hügelrunde",
    "baseName": "Hügelrunde",
    "waypointData": ""
  },
  {
    "variant": "Stadtrunde",
    "length": 7056,
    "id": "stadtrunde",
    "baseId": "parkparcours",
    "baseName": "Parkparcours",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 12890,
    "id": "hügelvariante",
    "baseId": "parkparcours",
    "baseName": "Parkparcours",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 13192,
    "id": "kurz",
    "baseId": "parkparcours",
    "baseName": "Parkparcours",
    "waypointData": ""
  },
  {
    "variant": "lang",
    "length": 7868,
    "id": "lang",
    "baseId": "stadtparklauf",
    "baseName": "Stadtparklauf",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 14665,
    "id": "hügelvariante",
    "baseId": "stadtparklauf",
    "baseName": "Stadtparklauf",
    "waypointData": ""
  },
  {
    "variant": "Rundweg",
    "length": 17951,
    "id": "rundweg",
    "baseId": "stadtparklauf",
    "baseName": "Stadtparklauf",
    "waypointData": ""
  },
  {
    "variant": "Seeuferweg",
    "length": 15332,
    "id": "seeuferweg",
    "baseId": "waldweg",
    "baseName": "Waldweg",
    "waypointData": ""
  },
  {
    "variant": "Rundweg",
    "length": 17467,
    "id": "rundweg",
    "baseId": "waldweg",
    "baseName": "Waldweg",
    "waypointData": ""
  },
  {
    "variant": "Parkstrecke",
    "length": 11892,
    "id": "parkstrecke",
    "baseId": "waldweg",
    "baseName": "Waldweg",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 12082,
    "id": "hügelvariante",
    "baseId": "waldweg",
    "baseName": "Waldweg",
    "waypointData": ""
  },
  {
    "variant": "Rundkurs",
    "length": 17464,
    "id": "rundkurs",
    "baseId": "flussuferweg",
    "baseName": "Flussuferweg",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 12301,
    "id": "hügelvariante",
    "baseId": "flussuferweg",
    "baseName": "Flussuferweg",
    "waypointData": ""
  },
  {
    "variant": "Stadtrunde",
    "length": 16527,
    "id": "stadtrunde",
    "baseId": "strandlauf",
    "baseName": "Strandlauf",
    "waypointData": ""
  },
  {
    "variant": "Seeuferweg",
    "length": 5812,
    "id": "seeuferweg",
    "baseId": "strandlauf",
    "baseName": "Strandlauf",
    "waypointData": ""
  },
  {
    "variant": "lang",
    "length": 6990,
    "id": "lang",
    "baseId": "seeumrundung",
    "baseName": "Seeumrundung",
    "waypointData": ""
  },
  {
    "variant": "Hügelvariante",
    "length": 19736,
    "id": "hügelvariante",
    "baseId": "seeumrundung",
    "baseName": "Seeumrundung",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 9268,
    "id": "kurz",
    "baseId": "seeumrundung",
    "baseName": "Seeumrundung",
    "waypointData": ""
  },
  {
    "variant": "kurz",
    "length": 17263,
    "id": "kurz",
    "baseId": "feldweglauf",
    "baseName": "Feldweglauf",
    "waypointData": ""
  }
]

export function useTracksApi() {
  function getTracks (): Promise<backend.Track[]> {
    return Promise.resolve(mockTracks)
  }
  return {getTracks}
}
