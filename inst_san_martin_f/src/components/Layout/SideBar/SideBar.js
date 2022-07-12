import React from 'react';
import { NavLink } from 'react-router-dom';

import { useGlobal } from "./../../../context/Global/GlobalProvider";

import styles from './SideBar.module.css';

const SideBar = ({ show, click }) => {

  const [globalState] = useGlobal();
  const { userLogin } = globalState;

  return (
    <nav className={(show ? [styles.side_bar, styles.side_bar__open].join(" ") : styles.side_bar)}>
      <div className={styles.side_bar__close}>
        <span onClick={click}>X</span>
      </div>
      {(() => {
        switch(userLogin.type) {
          case "ALUMNO": return (
            <ul>
              <li>
                <NavLink 
                  className={styles.side_bar__link}
                  to="/logout">
                  Logout
                </NavLink>
              </li>
            </ul> 
          );
          case "ADMINISTRATIVO": return (
            <ul>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/students">
                  Alumnos
                </NavLink>
              </li>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/logout">
                  Logout
                </NavLink>
              </li>
            </ul>
          );
          case "ADMINISTRADOR": return(
            <ul>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/roles">
                  Roles
                </NavLink>
              </li>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/users">
                  Administrativos
                </NavLink>
              </li>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/logout">
                  Logout
                </NavLink>
              </li>
            </ul>
          );
          default: return (
            <ul>
              <li>
                <NavLink 
                  className={styles.side_bar__link}
                  to="/login">
                  Login
                </NavLink>
              </li>
            </ul> 
          );
        }
      })()}
    </nav>
  );
};


export default SideBar;
