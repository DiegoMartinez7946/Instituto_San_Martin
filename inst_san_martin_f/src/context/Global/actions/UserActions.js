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

  const shiftIds = Array.isArray(item.shiftIds)
    ? item.shiftIds.map((x) => String(x).trim()).filter(Boolean)
    : item.shiftId
      ? [String(item.shiftId).trim()].filter(Boolean)
      : [];

  const degreeIds = Array.isArray(item.degreeIds)
    ? item.degreeIds.map((x) => String(x).trim()).filter(Boolean)
    : [];

  const result = await clientAxios.post('/user', {
    email: item.email,
    password: item.password,
    userType: item.userType,
    name: item.name,
    dni: item.dni,
    address: item.address || '',
    phone: item.phone || '',
    shiftIds,
    degreeIds
  }, options);

  return result.data;
};

export const updateUser = async (dispatch, item) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const shiftIds = Array.isArray(item.shiftIds)
    ? item.shiftIds.map((x) => String(x).trim()).filter(Boolean)
    : item.shiftId
      ? [String(item.shiftId).trim()].filter(Boolean)
      : [];

  const degreeIds = Array.isArray(item.degreeIds)
    ? item.degreeIds.map((x) => String(x).trim()).filter(Boolean)
    : [];

  const result = await clientAxios.put('/user', {
    id: item.id,
    email: item.email,
    name: item.name,
    dni: item.dni,
    address: item.address || '',
    phone: item.phone || '',
    userType: item.userType || '',
    password: item.password || '',
    shiftIds,
    degreeIds,
    modalidad: item.modalidad || '',
    condicion: item.condicion || '',
    studentRecordId: item.studentRecordId != null ? String(item.studentRecordId).trim() : ''
  }, options);

  return result.data;
};

export const changeActiveUser = async (dispatch, payload) => {
  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };

  const result = await clientAxios.put('/user/active', {
    id: payload.id,
    active: payload.active === true
  }, options);

  return result.data;
};
