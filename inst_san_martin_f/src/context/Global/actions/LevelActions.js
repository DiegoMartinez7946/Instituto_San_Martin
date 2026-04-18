import clientAxios from '../../../config/axios';

export const getLevels = async (dispatch) => {
  const access_token = document.cookie.replace('token=', '');
  const result = await clientAxios.get('/level', {
    headers: { Authorization: `Bearer${access_token}` }
  });
  dispatch({ type: 'GET_LEVELS', payload: result.data.data });
  return result.data.data;
};
