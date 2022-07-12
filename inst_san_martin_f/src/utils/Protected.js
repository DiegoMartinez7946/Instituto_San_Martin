import React from 'react';
import lodash from 'lodash';
import { Navigate } from "react-router-dom";
import decodeToken from './jwt';

import { ScreenPermission } from '../constant/constant';

const Protected = ({ children }) => {

  const isLoggedIn = !lodash.isEmpty(document.cookie);
  const user = isLoggedIn ? decodeToken(document.cookie) : null;
  const type = user ? lodash.lowerCase(user.type) : null;

  switch(type) {
    case 'alumno':
      return isLoggedIn && lodash.includes(ScreenPermission['alumno'], children.type.name) ? children : <Navigate to='/main' replace />;
    case 'administrativo':
      return isLoggedIn && lodash.includes(ScreenPermission['administrativo'], children.type.name) ? children : <Navigate to='/main' replace />;
    case 'administrador':
      return isLoggedIn && lodash.includes(ScreenPermission['administrador'], children.type.name) ? children : <Navigate to='/main' replace />;
    default:
      return isLoggedIn ? children : <Navigate to='/' replace />;
  }

};

export default Protected;
