export default (state, action) => {
  switch (action.type) {
    case "SET_PORTAL_TIME":
      return {
        ...state,
        portalTime: action.payload || null,
      };
    case "LOGOUT":
      return {
        ...state,
        portalTime: null,
      };
    default:
      return state;
  }
};
