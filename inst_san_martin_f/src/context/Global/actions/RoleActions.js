import clientAxios from '../../../config/axios';

export const addRole = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  }
  const result = await clientAxios.post('/user/role', 
    { type: item.type }, options);
  item.ID = result.data.data || '';
 
  dispatch({
    type: 'ADD_ROLE',
    payload: item
  });

  return result.data;
};

export const updateRole = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/user/role', 
    { id: item.ID, type: item.type }, options);

  dispatch({
    type: 'UPDATE_ROLE',
    payload: item
  });
  
  return result.data;
};

export const deleteRole = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.delete(`/user/role?id=${item.ID}`, {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });

  if (result.data.code !== 200) return result.data;

  dispatch({
    type: 'DELETE_ROLE',
    payload: item
  });

  return result.data;
};

export const getRole = async (dispatch) => {
  
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.get('/user/role', {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });
  
  dispatch({
    type: 'GET_ROLE',
    payload: result.data.data
  });

  return result.data.data;
};
