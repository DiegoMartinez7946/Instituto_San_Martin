import clientAxios from '../../../config/axios';

export const addPursueType = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  }
  const result = await clientAxios.post('/pursuetype', 
    { type: item.type }, options);
  item.ID = result.data.data || '';
 
  dispatch({
    type: 'ADD_PURSUE_TYPE',
    payload: item
  });

  return result.data;
};

export const updatePursueType = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/pursuetype', 
    { id: item.ID, type: item.type }, options);

  dispatch({
    type: 'UPDATE_PURSUE_TYPE',
    payload: item
  });
  
  return result.data;
};

export const deletePursueType = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.delete(`/pursuetype?id=${item.ID}`, {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });

  if (result.data.code !== 200) return result.data;

  dispatch({
    type: 'DELETE_PURSUE_TYPE',
    payload: item
  });

  return result.data;
};

export const getPursueType = async (dispatch) => {
  
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.get('/pursuetype', {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });
  
  dispatch({
    type: 'GET_PURSUE_TYPE',
    payload: result.data.data
  });

  return result.data.data;
};
