import React, { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAngleDown, faAngleRight } from '@fortawesome/free-solid-svg-icons';

import { useGlobal } from "./../../../context/Global/GlobalProvider";
import { decodeTokenFromCookie, normalizeRoleFromToken } from '../../../utils/jwt';

import styles from './SideBar.module.css';

const SideBar = ({ show, click }) => {

  const [globalState] = useGlobal();
  const [icon, setIcon] = useState(null);
  const [iconPwdMenu, setIconPwdMenu] = useState(null);
  const [iconAdmMenu, setIconAdmMenu] = useState(null);
  const { userLogin } = globalState;

  useEffect(() => {
    const down = <FontAwesomeIcon icon={faAngleDown} />;
    setIcon(down);
    setIconPwdMenu(down);
    setIconAdmMenu(down);
  }, [globalState]);

  const sessionRole =
    normalizeRoleFromToken(userLogin) || normalizeRoleFromToken(decodeTokenFromCookie());

  return (
    <nav className={(show ? [styles.side_bar, styles.side_bar__open].join(" ") : styles.side_bar)}>
      <div className={styles.side_bar__close}>
        <span onClick={click}>X</span>
      </div>
      {(() => {
        switch (sessionRole) {
          case "ALUMNO": return (
            <ul>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/passwordchange">
                  Cambiar Password
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
          case "DOCENTE": return (
            <ul>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/passwordchange">
                  Cambiar Password
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
          case "ADMINISTRATIVO": return (
            <ul>
              <li>
                <div
                  className={styles.sidebar__dropdown_container}
                  onMouseEnter={() => setIcon(<FontAwesomeIcon icon={faAngleRight} />)}
                  onMouseLeave={() => setIcon(<FontAwesomeIcon icon={faAngleDown} />)}
                >
                  <span>Administrar {icon}</span>
                  <ul>
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/students"
                      >
                        Alumnos
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/teachers"
                      >
                        Docentes
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/degree"
                      >
                        Carreras
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/studyplan"
                      >
                        Plan de estudio
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/shift"
                      >
                        Turnos
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/testtype"
                      >
                        Tipo Examen
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/pursuetype"
                      >
                        Modalidad
                      </NavLink>
                    </li>
                  </ul>
                </div>
              </li>
              <li>
                <NavLink
                  className={styles.side_bar__link}
                  to="/passwordchange">
                  Cambiar Password
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
                <div 
                  className={styles.sidebar__dropdown_container}
                  onMouseEnter={() => setIconPwdMenu(<FontAwesomeIcon icon={faAngleRight} />)}
                  onMouseLeave={() => setIconPwdMenu(<FontAwesomeIcon icon={faAngleDown} />)}
                >
                  <span>Contraseñas {iconPwdMenu}</span>
                  <ul>
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/passwordblank">
                        Blanquear usuario (staff)
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/passwordchange">
                        Cambiar mi contraseña
                      </NavLink>
                    </li>
                  </ul>
                </div>
              </li>
              <li>
                <div 
                  className={styles.sidebar__dropdown_container}
                  onMouseEnter={() => setIconAdmMenu(<FontAwesomeIcon icon={faAngleRight} />)}
                  onMouseLeave={() => setIconAdmMenu(<FontAwesomeIcon icon={faAngleDown} />)}
                >
                  <span>Administrar {iconAdmMenu}</span>
                  <ul>
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/users">
                        Usuarios
                      </NavLink>
                    </li>
                    <hr />
                    <li>
                      <NavLink
                        className={styles.side_bar__link_submenu}
                        to="/roles">
                        Roles
                      </NavLink>
                    </li>
                  </ul>
                </div>
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
