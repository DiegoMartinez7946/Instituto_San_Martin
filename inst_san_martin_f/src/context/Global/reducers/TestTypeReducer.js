const asList = (items) => (Array.isArray(items) ? items : []);

export default (state, action) => {
  const list = asList(state.testtypes);

  switch (action.type) {
    case 'ADD_TEST_TYPE':
      return {
        ...state,
        testtypes: [...list, action.payload]
      };
    case 'UPDATE_TEST_TYPE':
      return {
        ...state,
        testtypes: list.map((t) =>
          t.ID === action.payload.ID ? { ...t, type: action.payload.type } : t
        )
      };
    case 'DELETE_TEST_TYPE':
      return {
        ...state,
        testtypes: list.filter((f) => f.ID !== action.payload.ID)
      };
    case 'GET_TEST_TYPE':
      return {
        ...state,
        testtypes: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
