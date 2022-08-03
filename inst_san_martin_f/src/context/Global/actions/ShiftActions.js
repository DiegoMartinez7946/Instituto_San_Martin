import clientAxios from '../../../config/axios';

export const addShift = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  }
  const result = await clientAxios.post('/shift', 
    { type: item.type }, options);
  item.ID = result.data.data || '';
 
  dispatch({
    type: 'ADD_SHIFT',
    payload: item
  });

  return result.data;
};

export const updateShift = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/shift', 
    { id: item.ID, type: item.type }, options);

  dispatch({
    type: 'UPDATE_SHIFT',
    payload: item
  });
  
  return result.data;
};

export const deleteShift = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.delete(`/shift?id=${item.ID}`, {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });

  if (result.data.code !== 200) return result.data;

  dispatch({
    type: 'DELETE_SHIFT',
    payload: item
  });

  return result.data;
};

export const getShift = async (dispatch) => {
  
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.get('/shift', {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });
  
  dispatch({
    type: 'GET_SHIFT',
    payload: result.data.data
  });

  return result.data.data;
};
