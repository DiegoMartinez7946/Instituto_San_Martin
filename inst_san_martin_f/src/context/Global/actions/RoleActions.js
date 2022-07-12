import clientAxios from '../../../config/axios';

export const addRole = (dispatch, item) => {
  dispatch({
    type: 'ADD_ROLE',
    payload: item
  });
};

export const updateRole = (dispatch, item) => {
  dispatch({
    type: 'UPDATE_ROLE',
    payload: item
  });
};

export const deleteRole = (dispatch, item) => {
  dispatch({
    type: 'DELETE_ROLE',
    payload: item
  });
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
    payload: result.data
  });
};
