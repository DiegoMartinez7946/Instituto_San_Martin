export default (state, action) => {
  switch (action.type) {
    case 'GET_LEVELS':
      return { ...state, levels: Array.isArray(action.payload) ? action.payload : [] };
    default:
      return state;
  }
};
