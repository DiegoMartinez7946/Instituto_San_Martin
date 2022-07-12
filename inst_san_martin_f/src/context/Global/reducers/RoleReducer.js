export default (state, action) => {
  switch(action.type) {
    case 'ADD_ROLE':
      return state;
    case 'UPDATE_ROLE':
      return state;
    case 'DELETE_ROLE':
      return state;
    case 'GET_ROLE':
      return {
        ...state,
        roles: action.payload
      };
    default:
      return state;
  }
};
