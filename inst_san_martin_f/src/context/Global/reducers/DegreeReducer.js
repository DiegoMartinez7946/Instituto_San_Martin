export default (state, action) => {
  switch(action.type) {
    case 'ADD_DEGREE':
      const degree = action.payload;
      const prevDegrees = Array.isArray(state.degrees) ? state.degrees : [];
      return {
        ...state,
        degrees: [...prevDegrees, degree]
      };
    case 'UPDATE_DEGREE':
      const updatedDegrees = (state.degrees || []).map(d => {
        const did = d.id || d.ID;
        const pid = action.payload.id || action.payload.ID;
        if (did === pid) {
          d.name = action.payload.name;
          d.nivel = action.payload.nivel;
          d.resolucionId = action.payload.resolucionId;
        }
        return d;
      });
      return {
        ...state,
        degrees: updatedDegrees 
      };
    case 'CHANGE_ACTIVE_DEGREE':
      const changedDegrees = (state.degrees || []).map(d => {
        const did = d.id || d.ID;
        const pid = action.payload.id || action.payload.ID;
        if (did === pid) {
          d.active = action.payload.active;
        }
        return d;
      });
      return {
        ...state,
        degrees: changedDegrees 
      };
    case 'GET_DEGREE':
      return {
        ...state,
        degrees: Array.isArray(action.payload) ? action.payload : []
      };
    default:
      return state;
  }
};
