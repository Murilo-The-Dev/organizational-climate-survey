import { apiGet, apiPost, apiPut, apiDelete } from '@/lib/api';
import type { UsuarioAdministrador, CreateUsuarioRequest, StatusUsuario } from '@/lib/types';

export const usuarioService = {
  listByEmpresa(empresaId: number): Promise<UsuarioAdministrador[]> {
    return apiGet<UsuarioAdministrador[]>(`/empresas/${empresaId}/usuarios-administradores`);
  },

  getById(id: number): Promise<UsuarioAdministrador> {
    return apiGet<UsuarioAdministrador>(`/usuarios-administradores/${id}`);
  },

  getByEmail(email: string): Promise<UsuarioAdministrador> {
    return apiGet<UsuarioAdministrador>(`/usuarios-administradores/email/${email}`);
  },

  create(data: CreateUsuarioRequest & { id_empresa: number }): Promise<UsuarioAdministrador> {
    return apiPost<UsuarioAdministrador>('/usuarios-administradores', data);
  },

  update(id: number, data: Partial<CreateUsuarioRequest>): Promise<UsuarioAdministrador> {
    return apiPut<UsuarioAdministrador>(`/usuarios-administradores/${id}`, data);
  },

  updateStatus(id: number, status: StatusUsuario): Promise<UsuarioAdministrador> {
    return apiPut<UsuarioAdministrador>(`/usuarios-administradores/${id}/status`, { status });
  },

  updatePassword(id: number, novaSenha: string): Promise<void> {
    return apiPut(`/usuarios-administradores/${id}/password`, { nova_senha: novaSenha });
  },

  delete(id: number): Promise<void> {
    return apiDelete(`/usuarios-administradores/${id}`);
  },
};
