import { dashboard } from "../../wailsjs/go/models";
import { LoadDashboard } from "../../wailsjs/go/dashboard/Assembler";
import DashboardDto = dashboard.DashboardDto;

export function useDashboardApi() {
  function loadDashboardApi(from: Date, to: Date, topTracks: number): Promise<DashboardDto> {
    return LoadDashboard({ from, to, topTracks } as any);
  }

  return { loadDashboardApi };
}
