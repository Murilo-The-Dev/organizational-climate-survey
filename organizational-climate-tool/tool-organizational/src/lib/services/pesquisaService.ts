import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api';
import type { Pesquisa, CreatePesquisaRequest, StatusPesquisa } from '@/lib/types';

export const pesquisaService = {
  listByEmpresa(empresaId: number): Promise<Pesquisa[]> {
    return apiGet<Pesquisa[]>(`/empresas/${empresaId}/pesquisas`);
  },

  getById(id: number): Promise<Pesquisa> {
    return apiGet<Pesquisa>(`/pesquisas/${id}`);
  },

  getByLink(link: string): Promise<Pesquisa> {
    return apiGet<Pesquisa>(`/pesquisas/link/${link}`);
  },

  getFormularioPublico(link: string): Promise<{ pesquisa: Pesquisa; perguntas: import('@/lib/types').Pergunta[] }> {
    return apiGet(`/pesquisas/link/${link}/formulario`);
  },

  create(empresaId: number, data: CreatePesquisaRequest): Promise<Pesquisa> {
    return apiPost<Pesquisa>(`/empresas/${empresaId}/pesquisas`, data);
  },

  update(id: number, data: Partial<CreatePesquisaRequest>): Promise<Pesquisa> {
    return apiPut<Pesquisa>(`/pesquisas/${id}`, data);
  },

  updateStatus(id: number, status: StatusPesquisa): Promise<Pesquisa> {
    return apiPut<Pesquisa>(`/pesquisas/${id}/status`, { status });
  },

  delete(id: number): Promise<void> {
    return apiDelete(`/pesquisas/${id}`);
  },
};
