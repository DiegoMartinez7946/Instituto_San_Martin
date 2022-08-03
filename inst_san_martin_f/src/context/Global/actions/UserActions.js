import clientAxios from "../../../config/axios";

export const login = async (dispatch, item) => {
  const result = await clientAxios.post('/login', item);
  
  dispatch({
    type: 'LOGIN',
    payload: result.data
  });

  return result;
};

export const logout = async (dispatch) => {
  dispatch({
    type: 'LOGOUT',
    payload: {},
  });
};

export const getUsersByRole = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const result = await clientAxios.get('/user', options);

  dispatch({
    type: 'GETUSERSBYROLE',
    payload: result.data.data
  });

  return result.data.data;
};

export const blankPassword = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const result = await clientAxios.put('/user/password/blank',
    { email: item.email, newPassword: item.newPassword }, options);

  dispatch({
    type: 'BLANK_PASSWORD',
    payload: item
  });

  return result.data;
};

export const changePassword = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const result = await clientAxios.put('/user/password', {
    email: item.email,
    currentPassword: item.currentPassword,
    newPassword: item.newPassword
  }, options);

  dispatch({
    type: 'CHANGE_PASSWORD',
    payload: item
  });

  return result.data;
};

export const addUser = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const result = await clientAxios.post('/user', {
    email: item.email,
    password: item.password,
    userType: item.userType
  }, options);

  delete item.password;

  dispatch({
    type: 'ADD_USER',
    payload: item
  });

  return result.data;
};
