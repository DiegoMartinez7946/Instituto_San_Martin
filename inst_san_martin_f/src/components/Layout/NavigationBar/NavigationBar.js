import React, { useState, useEffect } from 'react';
import { NavLink } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faAngleDown, faAngleRight, faHouse } from '@fortawesome/free-solid-svg-icons';
import lodash from 'lodash';
import { decodeTokenFromCookie, normalizeRoleFromToken } from '../../../utils/jwt';
import { useGlobal } from './../../../context/Global/GlobalProvider';
import { getArgentinaPortalTime } from '../../../context/Global/actions/TimeActions';

import styles from './NavigationBar.module.css';

import ToggleButton from '../SideBar/ToggleButton';

const permisoLegible = (type) => {
  if (!type) return '';
  const u = String(type).toUpperCase();
  const map = {
    ADMINISTRATIVO: 'administrativo',
    ADMINISTRADOR: 'administrador',
    ALUMNO: 'alumno',
    DOCENTE: 'docente'
  };
  return map[u] || String(type).toLowerCase();
};

const NavigationBar = ({ toggleClick }) => {

  const [user, setUser] = useState('');
  const [icon, setIcon] = useState(null);
 
  const [globalState, dispatch] = useGlobal();

  useEffect(() => {
    const decoded = decodeTokenFromCookie();
    setUser(decoded || {});
    setIcon(<FontAwesomeIcon icon={faAngleDown} />);
  }, [globalState]);

  useEffect(() => {
    const token = document.cookie.replace("token=", "");
    if (!token) {
      return undefined;
    }
    getArgentinaPortalTime(dispatch);
    const id = setInterval(() => getArgentinaPortalTime(dispatch), 30000);
    return () => clearInterval(id);
  }, [dispatch, globalState.userLogin]);

  return (
    <header className={styles.toolbar}>
        <nav className={styles.toolbar__navigation}>
            <div className={styles.toolbar__toggle_button}>
                <ToggleButton click={toggleClick}/>
            </div>
            <div className={styles.toolbar__logo}>
              <NavLink
                to={lodash.isEmpty(user) ? '/' : '/home'}
                className={styles.toolbar__home}
                title="Inicio"
              >
                <FontAwesomeIcon icon={faHouse} className={styles.toolbar__home_icon} aria-hidden />
                <span className={styles.toolbar__home_label}>home</span>
              </NavLink>
            </div>
            {!lodash.isEmpty(user) && normalizeRoleFromToken(user) ? (
              <div
                className={styles.sessionInfo}
                role="status"
                aria-live="polite"
                title="Información de la sesión (solo lectura)"
              >
                usuario: <span className={styles.sessionInfo__strong}>{user.email || '—'}</span>{' '}
                permisos: <span className={styles.sessionInfo__strong}>{permisoLegible(normalizeRoleFromToken(user))}</span>
              </div>
            ) : null}
            {!lodash.isEmpty(user) && globalState.portalTime ? (
              <div
                className={styles.argClock}
                role="timer"
                aria-live="polite"
                title="Hora de referencia (Argentina, America/Argentina/Buenos_Aires)"
              >
                <span className={styles.argClock__time}>{globalState.portalTime.hora}</span>
                <span className={styles.argClock__meta}>
                  {globalState.portalTime.dia} · {globalState.portalTime.fecha}
                </span>
              </div>
            ) : null}
            <div className={styles.spacer} />
            <div className={styles.toolbar__navigation_items}>
              {(() => {
                const role = normalizeRoleFromToken(user);
                switch (role) {
                  case "ALUMNO": return (
                    <ul>
                      <li>
                        <NavLink
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                          to="/passwordchange">
                          Cambiar Password
                        </NavLink>
                      </li>
                      <li>
                        <NavLink 
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
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
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                          to="/passwordchange">
                          Cambiar Password
                        </NavLink>
                      </li>
                      <li>
                        <NavLink
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
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
                          className={styles.toolbar__dropdown_container}
                          onMouseEnter={() => setIcon(<FontAwesomeIcon icon={faAngleRight} />)}
                          onMouseLeave={() => setIcon(<FontAwesomeIcon icon={faAngleDown} />)}
                        >
                          <span>Administrar {icon}</span>
                          <ul>
                            <li>
                              <NavLink
                                className={({ isActive }) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/students"
                              >
                                Alumnos
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({ isActive }) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/degree"
                              >
                                Carreras
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({ isActive }) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/teachers"
                              >
                                Docentes
                              </NavLink>
                            </li>
                          </ul>
                        </div>
                      </li>
                      <li>
                        <NavLink
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                          to="/passwordchange">
                          Cambiar Password
                        </NavLink>
                      </li>
                      <li>
                        <NavLink
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
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
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                          to="/passwordblank">
                          Blanqueo Password
                        </NavLink>
                      </li>
                      <li>
                        <div 
                          className={styles.toolbar__dropdown_container}
                          onMouseEnter={() => setIcon(<FontAwesomeIcon icon={faAngleRight} />)}
                          onMouseLeave={() => setIcon(<FontAwesomeIcon icon={faAngleDown} />)}
                        >
                          <span>Administrar {icon}</span>
                          <ul>
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/users">
                                Usuarios
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/degree">
                                Carreras
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/roles">
                                Roles
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/shift">
                                Turnos
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/testtype">
                                Tipo Examen
                              </NavLink>
                            </li>
                            <hr />
                            <li>
                              <NavLink
                                className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                                to="/pursuetype">
                                Modalidad
                              </NavLink>
                            </li>
                          </ul>
                        </div>
                      </li>
                      <li>
                        <NavLink
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
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
                          className={({isActive}) => (isActive ? styles.toolbar__link_active : styles.toolbar__link)}
                          to="/login">
                          Login
                        </NavLink>
                      </li>
                    </ul> 
                  );
                }
              })()}
            </div>
        </nav>
    </header>
  );
};

export default NavigationBar;
