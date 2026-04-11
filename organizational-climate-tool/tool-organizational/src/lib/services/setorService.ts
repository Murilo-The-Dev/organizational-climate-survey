import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api';
import type { Setor, CreateSetorRequest } from '@/lib/types';

export const setorService = {
  listByEmpresa(empresaId: number): Promise<Setor[]> {
    return apiGet<Setor[]>(`/empresas/${empresaId}/setores`);
  },

  getById(id: number): Promise<Setor> {
    return apiGet<Setor>(`/setores/${id}`);
  },

  create(data: CreateSetorRequest): Promise<Setor> {
    return apiPost<Setor>('/setores', data);
  },

  update(id: number, data: Partial<CreateSetorRequest>): Promise<Setor> {
    return apiPut<Setor>(`/setores/${id}`, data);
  },

  delete(id: number): Promise<void> {
    return apiDelete(`/setores/${id}`);
  },
};
