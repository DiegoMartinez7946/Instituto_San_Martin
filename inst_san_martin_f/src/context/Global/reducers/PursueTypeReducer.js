export default (state, action) => {
  switch(action.type) {
    case 'ADD_PURSUE_TYPE':
      const pursueType = action.payload;
      const newPursueTypes = state.pursuetypes;
      newPursueTypes.push(pursueType);
      return {
        ...state,
        pursuetypes: newPursueTypes
      };
    case 'UPDATE_PURSUE_TYPE':
      const updatedPursueTypes = state.pursuetypes.map(p => {
        if (p.ID === action.payload.ID) {
          p.type = action.payload.type;
        }
        return p;
      });
      return {
        ...state,
        pursuetypes: updatedPursueTypes 
      };
    case 'DELETE_PURSUE_TYPE':
      const deletedPursueTypes = state.pursuetypes.filter(f => f.ID !== action.payload.ID)
      return {
        ...state,
        pursuetypes: deletedPursueTypes
      };
    case 'GET_PURSUE_TYPE':
      return {
        ...state,
        pursuetypes: action.payload
      };
    default:
      return state;
  }
};
