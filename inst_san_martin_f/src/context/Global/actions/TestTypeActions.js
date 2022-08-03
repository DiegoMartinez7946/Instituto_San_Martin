import clientAxios from '../../../config/axios';

export const addTestType = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  }
  const result = await clientAxios.post('/testtype', 
    { type: item.type }, options);
  item.ID = result.data.data || '';
 
  dispatch({
    type: 'ADD_TEST_TYPE',
    payload: item
  });

  return result.data;
};

export const updateTestType = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const options = {
    headers: {
      'Authorization': `Bearer${access_token}`
    },
  };
  const result = await clientAxios.put('/testtype', 
    { id: item.ID, type: item.type }, options);

  dispatch({
    type: 'UPDATE_TEST_TYPE',
    payload: item
  });
  
  return result.data;
};

export const deleteTestType = async (dispatch, item) => {

  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.delete(`/testtype?id=${item.ID}`, {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });

  if (result.data.code !== 200) return result.data;

  dispatch({
    type: 'DELETE_TEST_TYPE',
    payload: item
  });

  return result.data;
};

export const getTestType = async (dispatch) => {
  
  const access_token = document.cookie.replace("token=", "");
  const result = await clientAxios.get('/testtype', {
    headers: {
      'Authorization': `Bearer${access_token}`
    }
  });
  
  dispatch({
    type: 'GET_TEST_TYPE',
    payload: result.data.data
  });

  return result.data.data;
};
