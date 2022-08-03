import clientAxios from '../../../config/axios';

export const addDegree = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  }
  const result = await clientAxios.post('/degree', 
    { name: item.name, active: true }, options);
  item.ID = result.data.data || '';
 
  dispatch({
    type: 'ADD_DEGREE',
    payload: item
  });

  return result.data;
};

export const updateDegree = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/degree', 
    { id: item.ID, name: item.name }, options);

  dispatch({
    type: 'UPDATE_DEGREE',
    payload: item
  });
  
  return result.data;
};

export const changeActiveDegree = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/degree/active', 
    { id: item.ID, active: item.active }, options);

  if (result.data.code !== 200) return result.data;

  dispatch({
    type: 'CHANGE_ACTIVE_DEGREE',
    payload: item
  });

  return result.data;
};

export const getDegree = async (dispatch) => {
  
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.get('/degree', {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });
  
  dispatch({
    type: 'GET_DEGREE',
    payload: result.data.data
  });

  return result.data.data;
};
