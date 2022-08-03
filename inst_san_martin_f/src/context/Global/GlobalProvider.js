import React, { useContext, useReducer } from 'react';
import { GlobalContext } from './GlobalContext';
import  combineReducers from './reducers';
import  DegreeReducer from './reducers/DegreeReducer';
import  PursueTypeReducer from './reducers/PursueTypeReducer';
import  RoleReducer from './reducers/RoleReducer';
import  ShiftReducer from './reducers/ShiftReducer';
import  TestTypeReducer from './reducers/TestTypeReducer';
import  UserReducer  from './reducers/UserReducer';
import { GlobalState } from './GlobalState';

export const useGlobal = () => {
  const { state, dispatch } = useContext(GlobalContext);
  return [state, dispatch];
};

export const GlobalProvider = ({ children }) => {

  const GlobalReducer = combineReducers({
    degrees: DegreeReducer,
    pursuetypes: PursueTypeReducer,
    roles: RoleReducer,
    shifts: ShiftReducer,
    testtypes: TestTypeReducer,
    users: UserReducer,
  });

  const [state, dispatch] = useReducer(GlobalReducer, GlobalState);

  return (
    <GlobalContext.Provider value = {{
      state: state,
      dispatch: dispatch
    }}>
      { children }
    </GlobalContext.Provider>
  );  
};
