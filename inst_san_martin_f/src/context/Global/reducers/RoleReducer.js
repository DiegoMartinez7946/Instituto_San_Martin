const asList = (roles) => (Array.isArray(roles) ? roles : []);

export default (state, action) => {
  const list = asList(state.roles);

  switch (action.type) {
    case 'ADD_ROLE':
      return {
        ...state,
        roles: [...list, action.payload]
      };
    case 'UPDATE_ROLE':
      return {
        ...state,
        roles: list.map((r) =>
          r.ID === action.payload.ID ? { ...r, type: action.payload.type } : r
        )
      };
    case 'DELETE_ROLE':
      return {
        ...state,
        roles: list.filter((f) => f.ID !== action.payload.ID)
      };
    case 'GET_ROLE':
      return {
        ...state,
        roles: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
