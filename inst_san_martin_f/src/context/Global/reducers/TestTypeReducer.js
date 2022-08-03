export default (state, action) => {
  switch(action.type) {
    case 'ADD_TEST_TYPE':
      const testType = action.payload;
      const newTestTypes = state.testtypes;
      newTestTypes.push(testType);
      return {
        ...state,
        testtypes: newTestTypes
      };
    case 'UPDATE_TEST_TYPE':
      const updatedTestTypes = state.testtypes.map(t => {
        if (t.ID === action.payload.ID) {
          t.type = action.payload.type;
        }
        return t;
      });
      return {
        ...state,
        testtypes: updatedTestTypes 
      };
    case 'DELETE_TEST_TYPE':
      const deletedTestTypes = state.testtypes.filter(f => f.ID !== action.payload.ID)
      return {
        ...state,
        testtypes: deletedTestTypes
      };
    case 'GET_TEST_TYPE':
      return {
        ...state,
        testtypes: action.payload
      };
    default:
      return state;
  }
};
