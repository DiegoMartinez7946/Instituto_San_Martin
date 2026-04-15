const asList = (items) => (Array.isArray(items) ? items : []);

export default (state, action) => {
  const list = asList(state.pursuetypes);

  switch (action.type) {
    case 'ADD_PURSUE_TYPE':
      return {
        ...state,
        pursuetypes: [...list, action.payload]
      };
    case 'UPDATE_PURSUE_TYPE':
      return {
        ...state,
        pursuetypes: list.map((p) =>
          p.ID === action.payload.ID ? { ...p, type: action.payload.type } : p
        )
      };
    case 'DELETE_PURSUE_TYPE':
      return {
        ...state,
        pursuetypes: list.filter((f) => f.ID !== action.payload.ID)
      };
    case 'GET_PURSUE_TYPE':
      return {
        ...state,
        pursuetypes: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
