import React, { useContext, useReducer } from 'react';
import { GlobalContext } from './GlobalContext';
import  combineReducers from './reducers';
import  RoleReducer from './reducers/RoleReducer.js';
import  UserReducer  from './reducers/UserReducer';
import { GlobalState } from './GlobalState';

export const useGlobal = () => {
  const { state, dispatch } = useContext(GlobalContext);
  return [state, dispatch];
};

export const GlobalProvider = ({ children }) => {

  const GlobalReducer = combineReducers({
    roles: RoleReducer,
    users: UserReducer
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
