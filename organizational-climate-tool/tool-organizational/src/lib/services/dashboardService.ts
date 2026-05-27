import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api';
import type { Dashboard, DashboardData } from '@/lib/types';

export const dashboardService = {
  getByPesquisa(pesquisaId: number): Promise<Dashboard> {
    return apiGet<Dashboard>(`/pesquisas/${pesquisaId}/dashboard`);
  },

  getById(id: number): Promise<Dashboard> {
    return apiGet<Dashboard>(`/dashboards/${id}`);
  },

  listByEmpresa(empresaId: number): Promise<Dashboard[]> {
    return apiGet<Dashboard[]>(`/empresas/${empresaId}/dashboards`);
  },

  create(pesquisaId: number, titulo: string): Promise<Dashboard> {
    return apiPost<Dashboard>('/dashboards', { id_pesquisa: pesquisaId, titulo });
  },

  update(id: number, data: Partial<Dashboard>): Promise<Dashboard> {
    return apiPut<Dashboard>(`/dashboards/${id}`, data);
  },

  getData(id: number): Promise<DashboardData> {
    return apiGet<DashboardData>(`/dashboards/${id}/data`);
  },

  getMetrics(id: number): Promise<DashboardData> {
    return apiGet<DashboardData>(`/dashboards/${id}/metrics`);
  },

  delete(id: number): Promise<void> {
    return apiDelete(`/dashboards/${id}`);
  },
};
