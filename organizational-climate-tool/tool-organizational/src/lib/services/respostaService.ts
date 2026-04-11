import { apiGet, apiPost, apiDelete } from '@/lib/api';
import type { Resposta, SubmitRespostaRequest } from '@/lib/types';

export const respostaService = {
  submit(data: SubmitRespostaRequest): Promise<void> {
    return apiPost('/respostas/submit', data);
  },

  getStatsByPesquisa(pesquisaId: number): Promise<unknown> {
    return apiGet(`/pesquisas/${pesquisaId}/respostas/stats`);
  },

  getAggregatedByPesquisa(pesquisaId: number): Promise<Resposta[]> {
    return apiGet<Resposta[]>(`/pesquisas/${pesquisaId}/respostas/aggregated`);
  },

  countByPesquisa(pesquisaId: number): Promise<number> {
    return apiGet<number>(`/pesquisas/${pesquisaId}/respostas/count`);
  },

  getByDateRange(pesquisaId: number, inicio: string, fim: string): Promise<Resposta[]> {
    return apiGet<Resposta[]>(`/pesquisas/${pesquisaId}/respostas/by-date`, { inicio, fim });
  },

  deleteByPesquisa(pesquisaId: number): Promise<void> {
    return apiDelete(`/pesquisas/${pesquisaId}/respostas`);
  },
};
