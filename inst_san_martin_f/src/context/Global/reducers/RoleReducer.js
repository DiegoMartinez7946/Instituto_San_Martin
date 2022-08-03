export default (state, action) => {
  switch(action.type) {
    case 'ADD_ROLE':
      const rol = action.payload;
      const newRoles = state.roles;
      newRoles.push(rol);
      return {
        ...state,
        roles: newRoles
      };
    case 'UPDATE_ROLE':
      const updatedRoles = state.roles.map(r => {
        if (r.ID === action.payload.ID) {
          r.type = action.payload.type;
        }
        return r;
      });
      return {
        ...state,
        roles: updatedRoles 
      };
    case 'DELETE_ROLE':
      const deletedRoles = state.roles.filter(f => f.ID !== action.payload.ID)
      return {
        ...state,
        roles: deletedRoles
      };
    case 'GET_ROLE':
      return {
        ...state,
        roles: action.payload
      };
    default:
      return state;
  }
};
