import React from 'react';
import lodash from 'lodash';
import { Navigate } from "react-router-dom";
import { decodeTokenFromCookie, getAuthTokenFromCookie } from './jwt';

import { ScreenPermission } from '../constant/constant';

const Protected = ({ children }) => {

  const token = getAuthTokenFromCookie();
  const isLoggedIn = Boolean(token);
  const user = isLoggedIn ? decodeTokenFromCookie() : null;
  const type = user ? lodash.lowerCase(user.type) : null;

  switch(type) {
    case 'alumno':
      return isLoggedIn && lodash.includes(ScreenPermission['alumno'], children.type.name) ? children : <Navigate to='/home' replace />;
    case 'administrativo':
      return isLoggedIn && lodash.includes(ScreenPermission['administrativo'], children.type.name) ? children : <Navigate to='/home' replace />;
    case 'administrador':
      return isLoggedIn && lodash.includes(ScreenPermission['administrador'], children.type.name) ? children : <Navigate to='/home' replace />;
    default:
      return isLoggedIn ? children : <Navigate to='/' replace />;
  }

};

export default Protected;
