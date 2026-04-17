const asList = (shifts) => (Array.isArray(shifts) ? shifts : []);

export default (state, action) => {
  const list = asList(state.shifts);

  switch (action.type) {
    case 'ADD_SHIFT':
      return {
        ...state,
        shifts: [...list, action.payload]
      };
    case 'UPDATE_SHIFT':
      return {
        ...state,
        shifts: list.map((s) =>
          s.ID === action.payload.ID ? { ...s, type: action.payload.type } : s
        )
      };
    case 'DELETE_SHIFT':
      return {
        ...state,
        shifts: list.filter((f) => f.ID !== action.payload.ID)
      };
    case 'GET_SHIFT':
      return {
        ...state,
        shifts: Array.isArray(action.payload)
          ? action.payload.map((s) => ({
              ...s,
              ID:
                s.ID ||
                s.id ||
                s._id ||
                (s._id && typeof s._id === 'object' && s._id.$oid ? s._id.$oid : '') ||
                ''
            }))
          : []
      };
    default:
      return state;
  }
};
