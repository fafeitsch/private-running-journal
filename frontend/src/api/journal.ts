import {backend} from '../../wailsjs/go/models';

export function useJournalApi() {
  async function getListEntries(): Promise<backend.JournalListEntry[]> {
    const result =[
      {
        date: '2023-10-14T05:20:13.185Z',
        trackName: 'Waldrunde',
        trackVariant: 'Seeuferweg',
        length: 16007,
        id: '1'
      },
      {
        date: '2022-10-01T03:37:57.295Z',
        trackName: 'Bergpfad',
        trackVariant: 'kurz',
        length: 13528,
        id: '2'
      },
      {
        date: '2022-05-29T14:27:19.436Z',
        trackName: 'Hügelrunde',
        trackVariant: 'Rundkurs',
        length: 8409,
        id: '3'
      },
      {
        date: '2023-06-02T13:45:09.937Z',
        trackName: 'Waldrunde',
        trackVariant: 'kurz',
        length: 7864,
        id: '4'
      },
      {
        date: '2022-05-29T16:14:19.186Z',
        trackName: 'Parkparcours',
        trackVariant: 'Stadtrunde',
        length: 7056,
        id: '5'
      },
      {
        date: '2022-02-27T11:01:12.552Z',
        trackName: 'Stadtparklauf',
        trackVariant: 'lang',
        length: 7868,
        id: '6'
      },
      {
        date: '2023-08-28T19:21:05.137Z',
        trackName: 'Bergpfad',
        trackVariant: 'Rundweg',
        length: 10702,
        id: '7'
      },
      {
        date: '2022-12-20T14:15:41.784Z',
        trackName: 'Waldweg',
        trackVariant: 'Seeuferweg',
        length: 15332,
        id: '8'
      },
      {
        date: '2023-07-23T19:32:20.109Z',
        trackName: 'Parkparcours',
        trackVariant: 'Hügelvariante',
        length: 12890,
        id: '9'
      },
      {
        date: '2023-10-04T21:54:16.578Z',
        trackName: 'Flussuferweg',
        trackVariant: 'Rundkurs',
        length: 17464,
        id: '10'
      },
      {
        date: '2023-05-01T00:38:28.316Z',
        trackName: 'Waldrunde',
        trackVariant: 'Hügelvariante',
        length: 17153,
        id: '11'
      },
      {
        date: '2022-11-10T04:24:31.028Z',
        trackName: 'Strandlauf',
        trackVariant: 'Stadtrunde',
        length: 16527,
        id: '12'
      },
      {
        date: '2023-04-20T00:57:05.740Z',
        trackName: 'Hügelrunde',
        trackVariant: 'Seeuferweg',
        length: 5065,
        id: '13'
      },
      {
        date: '2023-04-02T06:32:33.757Z',
        trackName: 'Parkparcours',
        trackVariant: 'kurz',
        length: 13192,
        id: '14'
      },
      {
        date: '2022-10-31T00:54:18.674Z',
        trackName: 'Seeumrundung',
        trackVariant: 'lang',
        length: 6990,
        id: '15'
      },
      {
        date: '2023-01-16T16:35:07.042Z',
        trackName: 'Waldweg',
        trackVariant: 'Rundweg',
        length: 17467,
        id: '16'
      },
      {
        date: '2022-07-16T15:24:26.567Z',
        trackName: 'Stadtparklauf',
        trackVariant: 'Hügelvariante',
        length: 14665,
        id: '17'
      },
      {
        date: '2023-11-12T06:05:43.922Z',
        trackName: 'Bergpfad',
        trackVariant: 'Parkstrecke',
        length: 16928,
        id: '18'
      },
      {
        date: '2022-07-07T08:23:29.823Z',
        trackName: 'Flussuferweg',
        trackVariant: 'Hügelvariante',
        length: 12301,
        id: '19'
      },
      {
        date: '2022-02-19T08:22:38.235Z',
        trackName: 'Waldrunde',
        trackVariant: 'Rundkurs',
        length: 18890,
        id: '20'
      },
      {
        date: '2023-01-05T19:46:42.392Z',
        trackName: 'Bergpfad',
        trackVariant: 'kurz',
        length: 20994,
        id: '21'
      },
      {
        date: '2022-10-19T03:43:33.298Z',
        trackName: 'Seeumrundung',
        trackVariant: 'Hügelvariante',
        length: 19736,
        id: '22'
      },
      {
        date: '2023-02-22T14:01:23.983Z',
        trackName: 'Waldweg',
        trackVariant: 'Parkstrecke',
        length: 11892,
        id: '23'
      },
      {
        date: '2022-05-26T11:13:04.015Z',
        trackName: 'Hügelrunde',
        trackVariant: 'kurz',
        length: 11589,
        id: '24'
      },
      {
        date: '2023-12-02T23:57:11.741Z',
        trackName: 'Stadtparklauf',
        trackVariant: 'Rundweg',
        length: 17951,
        id: '25'
      },
      {
        date: '2023-08-14T03:27:10.064Z',
        trackName: 'Strandlauf',
        trackVariant: 'Seeuferweg',
        length: 5812,
        id: '26'
      },
      {
        date: '2022-06-19T11:21:08.149Z',
        trackName: 'Feldweglauf',
        trackVariant: 'kurz',
        length: 17263,
        id: '27'
      },
      {
        date: '2022-10-12T08:01:36.238Z',
        trackName: 'Parkparcours',
        trackVariant: 'Stadtrunde',
        length: 16039,
        id: '28'
      },
      {
        date: '2023-05-17T06:11:13.938Z',
        trackName: 'Waldweg',
        trackVariant: 'Hügelvariante',
        length: 12082,
        id: '29'
      },
      {
        date: '2022-09-23T21:54:14.129Z',
        trackName: 'Seeumrundung',
        trackVariant: 'kurz',
        length: 9268,
        id: '30'
      }
    ]
    return Promise.resolve(result)
  }
  return {getListEntries}
}
