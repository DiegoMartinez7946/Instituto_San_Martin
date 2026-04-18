import clientAxios from '../../../config/axios';

/** Lista de carreras: API suele mandar degreeIds; si solo hay degreeShifts (Mongo sin degreeids), se deduce aquí. */
const degreeIdsFromApiStudent = (s, normalizedShifts) => {
  const fromIds = Array.isArray(s.degreeIds) ? s.degreeIds.map((x) => String(x)).filter(Boolean) : [];
  if (fromIds.length) return fromIds;
  return (normalizedShifts || []).map((r) => r.degreeId).filter(Boolean);
};

const authHeaders = () => {
  const access_token = document.cookie.replace('token=', '');
  return {
    headers: {
      Authorization: `Bearer${access_token}`
    }
  };
};

const buildStudentBody = (item) => {
  const email =
    item.email !== undefined && item.email !== null
      ? String(item.email).trim()
      : '';
  const np = item.newPassword != null ? String(item.newPassword).trim() : '';
  const degreeShifts = Array.isArray(item.degreeShifts)
    ? item.degreeShifts
        .map((r) => ({
          degreeId: r.degreeId != null ? String(r.degreeId).trim() : '',
          shiftId: r.shiftId != null ? String(r.shiftId).trim() : ''
        }))
        .filter((r) => r.degreeId && r.shiftId)
    : [];
  let degreeIds = Array.isArray(item.degreeIds)
    ? item.degreeIds.map((x) => String(x).trim()).filter(Boolean)
    : [];
  if (!degreeIds.length && degreeShifts.length) {
    degreeIds = degreeShifts.map((r) => r.degreeId);
  }
  const body = {
    name: item.name != null ? String(item.name) : '',
    email,
    phone: item.phone !== undefined && item.phone !== null ? String(item.phone) : '',
    dni: item.dni != null ? String(item.dni) : '',
    address: item.address !== undefined && item.address !== null ? String(item.address) : '',
    nivelAprobado: item.nivelAprobado != null ? String(item.nivelAprobado).trim() : '',
    modalidad: item.modalidad != null ? String(item.modalidad).trim() : '',
    condicion: item.condicion != null ? String(item.condicion).trim() : '',
    degreeIds,
    degreeShifts,
    active: item.active === false ? false : true
  };
  if (np.length >= 6) {
    body.newPassword = np;
  }
  return body;
};

export const addStudent = async (dispatch, item) => {
  const result = await clientAxios.post('/student', buildStudentBody(item), authHeaders());
  return result.data;
};

export const updateStudent = async (dispatch, item) => {
  const body = {
    ...buildStudentBody(item),
    id: item.id != null ? String(item.id) : ''
  };
  const result = await clientAxios.put('/student', body, authHeaders());
  return result.data;
};

export const changeActiveStudent = async (dispatch, item) => {
  const result = await clientAxios.put(
    '/student/active',
    { id: item.id != null ? String(item.id) : '', active: item.active === true },
    authHeaders()
  );
  if (Number(result.data.code) === 200) {
    await getStudents(dispatch);
  }
  return result.data;
};

export const getStudents = async (dispatch) => {
  const result = await clientAxios.get('/student', authHeaders());
  const list = result.data.data || [];
  const normalized = list.map((s) => {
    const rawShifts = Array.isArray(s.degreeShifts)
      ? s.degreeShifts
      : Array.isArray(s.degree_shifts)
        ? s.degree_shifts
        : [];
    const degreeShifts = rawShifts
      .map((r) => ({
        degreeId: String(r.degreeId != null ? r.degreeId : r.degreeid || '').trim(),
        shiftId: String(r.shiftId != null ? r.shiftId : r.shiftid || '').trim()
      }))
      .filter((r) => r.degreeId && r.shiftId);
    const degreeIds = degreeIdsFromApiStudent(s, degreeShifts);
    return {
      id: s.id != null ? String(s.id) : '',
      name: s.name,
      email: s.email != null ? s.email : s.Email || '',
      phone: s.phone || '',
      dni: s.dni,
      address: s.address || '',
      nivelAprobado: s.nivelAprobado || '',
      modalidad: (s.modalidad != null && s.modalidad !== '' ? s.modalidad : s.Modalidad) || '',
      condicion: (s.condicion != null && s.condicion !== '' ? s.condicion : s.Condicion) || '',
      degreeIds,
      degreeShifts,
      active: s.active !== false
    };
  });
  dispatch({
    type: 'GET_STUDENTS',
    payload: normalized
  });
  return normalized;
};
