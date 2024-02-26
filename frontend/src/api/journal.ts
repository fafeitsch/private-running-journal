import { backend } from "../../wailsjs/go/models";

const mockEntries = [
  {
    date: "2023-10-14T05:20:13.185Z",
    trackBaseName: "Waldrunde",
    trackVariant: "Seeuferweg",
    length: 16007,
    id: "1",
  },
  {
    date: "2022-10-01T03:37:57.295Z",
    trackBaseName: "Bergpfad",
    trackVariant: "kurz",
    length: 13528,
    id: "2",
  },
  {
    date: "2022-05-29T14:27:19.436Z",
    trackBaseName: "Hügelrunde",
    trackVariant: "Rundkurs",
    length: 8409,
    id: "3",
  },
  {
    date: "2023-06-02T13:45:09.937Z",
    trackBaseName: "Waldrunde",
    trackVariant: "kurz",
    length: 7864,
    id: "4",
  },
  {
    date: "2022-05-29T16:14:19.186Z",
    trackBaseName: "Parkparcours",
    trackVariant: "Stadtrunde",
    length: 7056,
    id: "5",
  },
  {
    date: "2022-02-27T11:01:12.552Z",
    trackBaseName: "Stadtparklauf",
    trackVariant: "lang",
    length: 7868,
    id: "6",
  },
  {
    date: "2023-08-28T19:21:05.137Z",
    trackBaseName: "Bergpfad",
    trackVariant: "Rundweg",
    length: 10702,
    id: "7",
  },
  {
    date: "2022-12-20T14:15:41.784Z",
    trackBaseName: "Waldweg",
    trackVariant: "Seeuferweg",
    length: 15332,
    id: "8",
  },
  {
    date: "2023-07-23T19:32:20.109Z",
    trackBaseName: "Parkparcours",
    trackVariant: "Hügelvariante",
    length: 12890,
    id: "9",
  },
  {
    date: "2023-10-04T21:54:16.578Z",
    trackBaseName: "Flussuferweg",
    trackVariant: "Rundkurs",
    length: 17464,
    id: "10",
  },
  {
    date: "2023-05-01T00:38:28.316Z",
    trackBaseName: "Waldrunde",
    trackVariant: "Hügelvariante",
    length: 17153,
    id: "11",
  },
  {
    date: "2022-11-10T04:24:31.028Z",
    trackBaseName: "Strandlauf",
    trackVariant: "Stadtrunde",
    length: 16527,
    id: "12",
  },
  {
    date: "2023-04-20T00:57:05.740Z",
    trackBaseName: "Hügelrunde",
    trackVariant: "Seeuferweg",
    length: 5065,
    id: "13",
  },
  {
    date: "2023-04-02T06:32:33.757Z",
    trackBaseName: "Parkparcours",
    trackVariant: "kurz",
    length: 13192,
    id: "14",
  },
  {
    date: "2022-10-31T00:54:18.674Z",
    trackBaseName: "Seeumrundung",
    trackVariant: "lang",
    length: 6990,
    id: "15",
  },
  {
    date: "2023-01-16T16:35:07.042Z",
    trackBaseName: "Waldweg",
    trackVariant: "Rundweg",
    length: 17467,
    id: "16",
  },
  {
    date: "2022-07-16T15:24:26.567Z",
    trackBaseName: "Stadtparklauf",
    trackVariant: "Hügelvariante",
    length: 14665,
    id: "17",
  },
  {
    date: "2023-11-12T06:05:43.922Z",
    trackBaseName: "Bergpfad",
    trackVariant: "Parkstrecke",
    length: 16928,
    id: "18",
  },
  {
    date: "2022-07-07T08:23:29.823Z",
    trackBaseName: "Flussuferweg",
    trackVariant: "Hügelvariante",
    length: 12301,
    id: "19",
  },
  {
    date: "2022-02-19T08:22:38.235Z",
    trackBaseName: "Waldrunde",
    trackVariant: "Rundkurs",
    length: 18890,
    id: "20",
  },
  {
    date: "2023-01-05T19:46:42.392Z",
    trackBaseName: "Bergpfad",
    trackVariant: "kurz",
    length: 20994,
    id: "21",
  },
  {
    date: "2022-10-19T03:43:33.298Z",
    trackBaseName: "Seeumrundung",
    trackVariant: "Hügelvariante",
    length: 19736,
    id: "22",
  },
  {
    date: "2023-02-22T14:01:23.983Z",
    trackBaseName: "Waldweg",
    trackVariant: "Parkstrecke",
    length: 11892,
    id: "23",
  },
  {
    date: "2022-05-26T11:13:04.015Z",
    trackBaseName: "Hügelrunde",
    trackVariant: "kurz",
    length: 11589,
    id: "24",
  },
  {
    date: "2023-12-02T23:57:11.741Z",
    trackBaseName: "Stadtparklauf",
    trackVariant: "Rundweg",
    length: 17951,
    id: "25",
  },
  {
    date: "2023-08-14T03:27:10.064Z",
    trackBaseName: "Strandlauf",
    trackVariant: "Seeuferweg",
    length: 5812,
    id: "26",
  },
  {
    date: "2022-06-19T11:21:08.149Z",
    trackBaseName: "Feldweglauf",
    trackVariant: "kurz",
    length: 17263,
    id: "27",
  },
  {
    date: "2022-10-12T08:01:36.238Z",
    trackBaseName: "Parkparcours",
    trackVariant: "Stadtrunde",
    length: 16039,
    id: "28",
  },
  {
    date: "2023-05-17T06:11:13.938Z",
    trackBaseName: "Waldweg",
    trackVariant: "Hügelvariante",
    length: 10000,
    id: "29",
  },
  {
    date: "2022-09-23T21:54:14.129Z",
    trackBaseName: "Seeumrundung",
    trackVariant: "kurz",
    length: 10000,
    id: "30",
  },
];

export function useJournalApi() {
  async function getJournalEntries(): Promise<backend.JournalListEntry[]> {
    return Promise.resolve(mockEntries);
  }
  async function getJournalEntry(id: string): Promise<backend.JournalEntry> {
    const listEntry = mockEntries.find((entry) => entry.id === id);
    if (!listEntry) {
      throw new Error("not found");
    }
    const result = new backend.JournalEntry(
      JSON.stringify({
        id,
        date: listEntry.date,
        comment: "Rainy day",
        time: "01:45:47",
        track: {
          baseName: listEntry.trackBaseName,
          baseId: listEntry.trackBaseName.toLowerCase(),
          variant: listEntry.trackVariant,
          length: listEntry.length,
          id: listEntry.trackVariant.toLowerCase(),
          waypointData: "",
        },
      }),
    );
    return Promise.resolve(result);
  }
  return { getListEntries: getJournalEntries, getListEntry: getJournalEntry };
}
