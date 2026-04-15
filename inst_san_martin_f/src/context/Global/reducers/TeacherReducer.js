export default (state, action) => {
  switch (action.type) {
    case 'GET_TEACHERS':
      return {
        ...state,
        teachers: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
