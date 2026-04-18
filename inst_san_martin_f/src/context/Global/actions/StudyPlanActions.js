import clientAxios from '../../../config/axios';

const fail = (message, code = 400) => ({ message, code, data: null });

export const getStudyPlans = async (dispatch) => {
  const access_token = document.cookie.replace('token=', '');
  try {
    const result = await clientAxios.get('/studyplan', {
      headers: { Authorization: `Bearer${access_token}` }
    });
    const body = result.data || {};
    const rows = body.data;
    const list = Array.isArray(rows) ? rows : [];
    dispatch({ type: 'GET_STUDY_PLAN', payload: list });
    return list;
  } catch (err) {
    const data = err.response && err.response.data;
    dispatch({ type: 'GET_STUDY_PLAN', payload: [] });
    return data || fail(err.message || 'Error de red');
  }
};

export const addStudyPlan = async (dispatch, item) => {
  const access_token = document.cookie.replace('token=', '');
  const options = { headers: { Authorization: `Bearer${access_token}` } };
  try {
    const result = await clientAxios.post(
      '/studyplan',
      {
        nombre: item.nombre,
        numeroResolucion: item.numeroResolucion,
        cohorte: Number(item.cohorte),
        extension: Number(item.extension),
        anioAprobacion: Number(item.anioAprobacion),
        anioCaducidad: Number(item.anioCaducidad),
        cargaHoraria: Number(item.cargaHoraria),
        active: true
      },
      options
    );
    const body = result.data || {};
    const codeNum = Number(body.code);
    if (codeNum !== 201) {
      return body;
    }
    const id = body.data || '';
    dispatch({
      type: 'ADD_STUDY_PLAN',
      payload: { ...item, id, ID: id, active: true }
    });
    return body;
  } catch (err) {
    const data = err.response && err.response.data;
    if (data) return data;
    return fail(err.message || 'Error de red');
  }
};

export const updateStudyPlan = async (dispatch, item) => {
  const access_token = document.cookie.replace('token=', '');
  const options = { headers: { Authorization: `Bearer${access_token}` } };
  try {
    const result = await clientAxios.put(
      '/studyplan',
      {
        id: item.id || item.ID,
        nombre: item.nombre,
        numeroResolucion: item.numeroResolucion,
        cohorte: Number(item.cohorte),
        extension: Number(item.extension),
        anioAprobacion: Number(item.anioAprobacion),
        anioCaducidad: Number(item.anioCaducidad),
        cargaHoraria: Number(item.cargaHoraria),
        active: item.active !== false
      },
      options
    );
    const body = result.data || {};
    if (Number(body.code) !== 200) {
      return body;
    }
    dispatch({ type: 'UPDATE_STUDY_PLAN', payload: item });
    return body;
  } catch (err) {
    const data = err.response && err.response.data;
    if (data) return data;
    return fail(err.message || 'Error de red');
  }
};

export const changeActiveStudyPlan = async (dispatch, item) => {
  const access_token = document.cookie.replace('token=', '');
  const options = { headers: { Authorization: `Bearer${access_token}` } };
  try {
    const result = await clientAxios.put(
      '/studyplan/active',
      { id: item.id || item.ID, active: item.active },
      options
    );
    const body = result.data || {};
    if (Number(body.code) !== 200) return body;
    dispatch({ type: 'CHANGE_ACTIVE_STUDY_PLAN', payload: item });
    return body;
  } catch (err) {
    const data = err.response && err.response.data;
    if (data) return data;
    return fail(err.message || 'Error de red');
  }
};
