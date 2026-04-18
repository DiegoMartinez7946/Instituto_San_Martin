import React, { useContext, useReducer } from 'react';
import { GlobalContext } from './GlobalContext';
import  combineReducers from './reducers';
import  DegreeReducer from './reducers/DegreeReducer';
import  PursueTypeReducer from './reducers/PursueTypeReducer';
import  RoleReducer from './reducers/RoleReducer';
import  ShiftReducer from './reducers/ShiftReducer';
import  TestTypeReducer from './reducers/TestTypeReducer';
import  UserReducer  from './reducers/UserReducer';
import  StudentReducer from './reducers/StudentReducer';
import  TeacherReducer from './reducers/TeacherReducer';
import StudyPlanReducer from './reducers/StudyPlanReducer';
import LevelReducer from './reducers/LevelReducer';
import PortalTimeReducer from './reducers/PortalTimeReducer';
import { GlobalState } from './GlobalState';

export const useGlobal = () => {
  const { state, dispatch } = useContext(GlobalContext);
  return [state, dispatch];
};

export const GlobalProvider = ({ children }) => {

  const GlobalReducer = combineReducers({
    levels: LevelReducer,
    studyPlans: StudyPlanReducer,
    degrees: DegreeReducer,
    students: StudentReducer,
    teachers: TeacherReducer,
    pursuetypes: PursueTypeReducer,
    roles: RoleReducer,
    shifts: ShiftReducer,
    testtypes: TestTypeReducer,
    users: UserReducer,
    portalTime: PortalTimeReducer,
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
