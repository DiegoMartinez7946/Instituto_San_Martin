export default (state, action) => {
  switch(action.type) {
    case 'ADD_DEGREE':
      const degree = action.payload;
      const newDegrees = state.degrees;
      newDegrees.push(degree);
      return {
        ...state,
        degress: newDegrees
      };
    case 'UPDATE_DEGREE':
      const updatedDegrees = state.degrees.map(d => {
        if (d.ID === action.payload.ID) {
          d.name = action.payload.name;
        }
        return d;
      });
      return {
        ...state,
        degrees: updatedDegrees 
      };
    case 'CHANGE_ACTIVE_DEGREE':
      const changedDegrees = state.degrees.map(d => {
        if (d.ID === action.payload.ID) {
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
        degrees: action.payload
      };
    default:
      return state;
  }
};
