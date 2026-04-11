import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api';
import type { Empresa, CreateEmpresaRequest } from '@/lib/types';

export const empresaService = {
  list(): Promise<Empresa[]> {
    return apiGet<Empresa[]>('/empresas');
  },

  getById(id: number): Promise<Empresa> {
    return apiGet<Empresa>(`/empresas/${id}`);
  },

  getByCnpj(cnpj: string): Promise<Empresa> {
    return apiGet<Empresa>(`/empresas/cnpj/${cnpj}`);
  },

  create(data: CreateEmpresaRequest): Promise<Empresa> {
    return apiPost<Empresa>('/empresas', data);
  },

  update(id: number, data: Partial<CreateEmpresaRequest>): Promise<Empresa> {
    return apiPut<Empresa>(`/empresas/${id}`, data);
  },

  delete(id: number): Promise<void> {
    return apiDelete(`/empresas/${id}`);
  },
};
