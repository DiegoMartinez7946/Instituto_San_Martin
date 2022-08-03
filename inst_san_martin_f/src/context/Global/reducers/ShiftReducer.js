
export default (state, action) => {
  switch(action.type) {
    case 'ADD_SHIFT':
      const shift = action.payload;
      const newShifts = state.shifts;
      newShifts.push(shift);
      return {
        ...state,
        shifts: newShifts
      };
    case 'UPDATE_SHIFT':
      const updatedShifts = state.shifts.map(s => {
        if (s.ID === action.payload.ID) {
          s.type = action.payload.type;
        }
        return s;
      });
      return {
        ...state,
        shifts: updatedShifts
      };
    case 'DELETE_SHIFT':
      const deletedShifts = state.shifts.filter(f => f.ID !== action.payload.ID)
      return {
        ...state,
        shifts: deletedShifts
      };
    case 'GET_SHIFT':
      return {
        ...state,
        shifts: action.payload
      };
    default:
      return state;
  }
};
