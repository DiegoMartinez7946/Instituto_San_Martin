import React, { useEffect } from 'react';
import { Navigate } from "react-router-dom";

import { useGlobal } from "./../context/Global/GlobalProvider";
import { logout } from './../context/Global/actions/UserActions';

const Logout = () => {

  const [, globalDispatch] = useGlobal();

  useEffect(() => {
    logout(globalDispatch, {});   
  }, [])

  return <Navigate to ="/" replace />;

};

export default Logout;
