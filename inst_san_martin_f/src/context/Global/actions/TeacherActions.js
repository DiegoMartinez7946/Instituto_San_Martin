import clientAxios from '../../../config/axios';

const authHeaders = () => {
  const access_token = document.cookie.replace('token=', '');
  return {
    headers: {
      Authorization: `Bearer${access_token}`
    }
  };
};

const buildTeacherBody = (item) => ({
  name: item.name != null ? String(item.name) : '',
  email: item.email !== undefined && item.email !== null ? String(item.email).trim() : '',
  phone: item.phone !== undefined && item.phone !== null ? String(item.phone) : '',
  dni: item.dni != null ? String(item.dni) : '',
  address: item.address !== undefined && item.address != null ? String(item.address) : '',
  enseniaEn: Array.isArray(item.enseniaEn) ? item.enseniaEn : [],
  degreeIds: Array.isArray(item.degreeIds) ? item.degreeIds : [],
  tituloHabilitanteId: item.tituloHabilitanteId != null ? String(item.tituloHabilitanteId) : '',
  modalidadId: item.modalidadId != null ? String(item.modalidadId) : ''
});

export const getTitulosHabilitantes = async () => {
  const res = await clientAxios.get('/titulo-habilitante', authHeaders());
  return res.data.data || [];
};

export const getModalidadesDocente = async () => {
  const res = await clientAxios.get('/modalidad-docente', authHeaders());
  return res.data.data || [];
};

export const addTeacher = async (dispatch, item) => {
  const result = await clientAxios.post('/teacher', buildTeacherBody(item), authHeaders());
  return result.data;
};

export const updateTeacher = async (dispatch, item) => {
  const body = {
    ...buildTeacherBody(item),
    id: item.id != null ? String(item.id) : ''
  };
  const result = await clientAxios.put('/teacher', body, authHeaders());
  return result.data;
};

export const getTeachers = async (dispatch) => {
  const result = await clientAxios.get('/teacher', authHeaders());
  const list = result.data.data || [];
  const normalized = list.map((t) => ({
    id: t.id,
    name: t.name,
    email: t.email != null ? t.email : '',
    phone: t.phone || '',
    dni: t.dni,
    address: t.address || '',
    enseniaEn: t.enseniaEn || [],
    degreeIds: t.degreeIds || [],
    tituloHabilitanteId: t.tituloHabilitanteId || '',
    modalidadId: t.modalidadId || ''
  }));
  dispatch({ type: 'GET_TEACHERS', payload: normalized });
  return normalized;
};
