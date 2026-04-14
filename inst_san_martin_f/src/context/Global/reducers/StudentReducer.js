export default (state, action) => {
  switch (action.type) {
    case 'GET_STUDENTS':
      return {
        ...state,
        students: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
