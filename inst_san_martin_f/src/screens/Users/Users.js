import React, { useState, useEffect, useMemo } from 'react';
import { Container, Row, Col, Button, Form, InputGroup } from 'react-bootstrap';
import Modal from './Modal/Modal';
import ModalTeacher from '../Teachers/Modal/ModalTeacher';
import ModalStudent from '../Students/Modal/ModalStudent';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getUsersByRole, addUser, updateUser, changeActiveUser } from '../../context/Global/actions/UserActions';
import { getRole } from '../../context/Global/actions/RoleActions';
import { getShift } from '../../context/Global/actions/ShiftActions';
import { getDegree } from '../../context/Global/actions/DegreeActions';
import {
  getTeachers,
  addTeacher,
  updateTeacher,
  changeActiveTeacher,
  getTitulosHabilitantes,
  getModalidadesDocente
} from '../../context/Global/actions/TeacherActions';
import {
  getStudents,
  addStudent,
  updateStudent,
  changeActiveStudent
} from '../../context/Global/actions/StudentActions';

import styles from './Users.module.css';
import listToolbar from '../common/ListToolbar.module.css';
import { decodeTokenFromCookie, normalizeRoleFromToken } from '../../utils/jwt';

/* Texto largo o varios turnos concatenados */
const USER_TABLE_WIDE_KEYS = ['nombre', 'email', 'rol', 'turno'];

/* Id interno por fila (no se muestra como columna: Table usa Object.keys). */
const USER_TABLE_ROW_ID = Symbol('userTableRowId');

const TEACHER_ID_PREFIX = 'teacher:';
const STUDENT_ID_PREFIX = 'student:';

/** Solo administrador o administrativo pueden editar filas de docente/alumno desde este listado. */
const canEditDocenteOAlumnoDesdeUsuarios = (sessionType) => {
  const t = String(sessionType || '').trim().toUpperCase();
  return t === 'ADMINISTRADOR' || t === 'ADMINISTRATIVO';
};

/**
 * Rol de la sesión: prioriza el JWT en la cookie (sigue válido tras recargar la página);
 * el contexto userLogin solo se rellena en el login explícito.
 */
const getSessionRoleType = (userLogin) =>
  normalizeRoleFromToken(decodeTokenFromCookie()) ||
  normalizeRoleFromToken(userLogin) ||
  '';

const Users = () => {
  const [globalState, globalDispatch] = useGlobal();
  const [dataUsers, setDataUsers] = useState([]);
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);

  const [showTeacherModal, setShowTeacherModal] = useState(false);
  const [teacherModalData, setTeacherModalData] = useState(null);
  const [showStudentModal, setShowStudentModal] = useState(false);
  const [studentModalData, setStudentModalData] = useState(null);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [titulosHabilitantes, setTitulosHabilitantes] = useState([]);
  const [modalidadesDocente, setModalidadesDocente] = useState([]);

  useEffect(() => {
    const load = async () => {
      await getRole(globalDispatch);
      await getShift(globalDispatch);
      await getDegree(globalDispatch);
      try {
        const [tit, mod] = await Promise.all([getTitulosHabilitantes(), getModalidadesDocente()]);
        setTitulosHabilitantes(Array.isArray(tit) ? tit : []);
        setModalidadesDocente(Array.isArray(mod) ? mod : []);
      } catch {
        setTitulosHabilitantes([]);
        setModalidadesDocente([]);
      }
      await getUsersByRole(globalDispatch);
    };
    load();
  }, [globalDispatch]);

  useEffect(() => {
    setDataUsers(globalState.users || []);
  }, [globalState]);

  const showError = (message, type) => {
    message
      ? setError(
          <Notification message={message} type={type} show={showError} />
        )
      : setError(null);
  };

  const buildNotification = (result) => {
    if (!result) return;
    switch (Number(result.code)) {
      case 199:
        showError(result.message, 'warning');
        break;
      case 200:
      case 201:
        showError(result.message, 'success');
        break;
      case 400:
      case 403:
      case 404:
        showError(result.message || 'Error', 'danger');
        break;
      default:
        break;
    }
  };

  const saveEventHandler = async (e) => {
    const result = e.id
      ? await updateUser(globalDispatch, e)
      : await addUser(globalDispatch, e);
    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    await getUsersByRole(globalDispatch);
    setShow(false);
    setDataRow(null);
    return result;
  };

  const addUserEvent = () => {
    setDataRow(null);
    setShow(true);
  };

  const closeUserEvent = () => {
    setDataRow(null);
    setShow(false);
  };

  const closeTeacherModal = () => {
    setTeacherModalData(null);
    setShowTeacherModal(false);
  };

  const closeStudentModal = () => {
    setStudentModalData(null);
    setShowStudentModal(false);
  };

  const saveTeacherFromUsers = async (e) => {
    const isNew = !e.id;
    try {
      const result = isNew
        ? await addTeacher(globalDispatch, e)
        : await updateTeacher(globalDispatch, e);
      buildNotification(result);
      const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
      if (codeNum !== 200 && codeNum !== 201) {
        return result;
      }
      await getTeachers(globalDispatch);
      await getUsersByRole(globalDispatch);
      closeTeacherModal();
      return result;
    } catch (err) {
      const data = err.response && err.response.data;
      const fail =
        data && data.message
          ? { ...data, code: data.code || err.response.status }
          : { message: err.message || 'Error de red', code: 400 };
      buildNotification(fail);
      return fail;
    }
  };

  const saveStudentFromUsers = async (e) => {
    try {
      const isNew = !e.id;
      const result = isNew
        ? await addStudent(globalDispatch, e)
        : await updateStudent(globalDispatch, e);
      buildNotification(result);
      const codeNum =
        result && result.code !== undefined && result.code !== null
          ? Number(result.code)
          : NaN;
      if (codeNum !== 200 && codeNum !== 201) {
        return result;
      }
      await getStudents(globalDispatch);
      await getUsersByRole(globalDispatch);
      closeStudentModal();
      return result;
    } catch (err) {
      const data = err.response && err.response.data;
      const fail =
        data && data.message
          ? { ...data, code: data.code || err.response.status }
          : { message: err.message || 'Error de red', code: 400 };
      buildNotification(fail);
      return fail;
    }
  };

  const tableEvents = async (ev, row) => {
    if (ev !== 'edit') return;

    const rowId = row[USER_TABLE_ROW_ID];
    const full = (dataUsers || []).find((u) => u.id === rowId);
    const sessionType = getSessionRoleType(globalState.userLogin);
    const roleUpper = String((full && full.role) || '').trim().toUpperCase();

    const openTeacherEdit = async () => {
      if (!canEditDocenteOAlumnoDesdeUsuarios(sessionType)) {
        showError(
          'Solo un administrador o usuario administrativo puede editar docentes desde este listado.',
          'warning'
        );
        return;
      }
      const raw = full && full.id ? String(full.id) : '';
      const hex = raw.startsWith(TEACHER_ID_PREFIX) ? raw.slice(TEACHER_ID_PREFIX.length) : '';
      if (!hex) {
        showError('Identificador de docente invalido.', 'danger');
        return;
      }
      try {
        const teachersList = await getTeachers(globalDispatch);
        const teacher = (teachersList || []).find((t) => String(t.id) === hex);
        if (!teacher) {
          showError('No se encontró el docente.', 'danger');
          return;
        }
        setTeacherModalData(teacher);
        setShowTeacherModal(true);
      } catch (err) {
        const msg =
          (err.response && err.response.data && err.response.data.message) ||
          err.message ||
          'No se pudo abrir la edición de docente';
        showError(msg, 'danger');
      }
    };

    const openStudentEdit = async () => {
      if (!canEditDocenteOAlumnoDesdeUsuarios(sessionType)) {
        showError(
          'Solo un administrador o usuario administrativo puede editar alumnos desde este listado.',
          'warning'
        );
        return;
      }
      const raw = full && full.id ? String(full.id) : '';
      const hex = raw.startsWith(STUDENT_ID_PREFIX) ? raw.slice(STUDENT_ID_PREFIX.length) : '';
      if (!hex) {
        showError('Identificador de alumno invalido.', 'danger');
        return;
      }
      try {
        const studentsList = await getStudents(globalDispatch);
        const student = (studentsList || []).find((s) => String(s.id) === hex);
        if (!student) {
          showError('No se encontró el alumno.', 'danger');
          return;
        }
        setStudentModalData(student);
        setShowStudentModal(true);
      } catch (err) {
        const msg =
          (err.response && err.response.data && err.response.data.message) ||
          err.message ||
          'No se pudo abrir la edición de alumno';
        showError(msg, 'danger');
      }
    };

    if (full && full.source === 'teacher') {
      await openTeacherEdit();
      return;
    }
    if (full && full.source === 'student') {
      await openStudentEdit();
      return;
    }

    /*
     * Misma persona puede tener registro en `user` (login) y ficha en `teacher` / `student`.
     * El backend no duplica por email: la tabla muestra la fila de `user` sin `source`.
     * En ese caso hay que abrir el modal de docente/alumno buscando por email.
     */
    if (full && !full.source && roleUpper === 'DOCENTE') {
      if (!canEditDocenteOAlumnoDesdeUsuarios(sessionType)) {
        showError(
          'Solo un administrador o usuario administrativo puede editar docentes desde este listado.',
          'warning'
        );
        return;
      }
      const teachersList = await getTeachers(globalDispatch);
      const em = String(full.email || '').trim().toLowerCase();
      const dni = String(full.dni || '').trim();
      if (em) {
        const teacher = (teachersList || []).find(
          (t) => String(t.email || '').trim().toLowerCase() === em
        );
        if (teacher) {
          setTeacherModalData(teacher);
          setShowTeacherModal(true);
          return;
        }
      }
      if (dni) {
        const teacherByDni = (teachersList || []).find(
          (t) => String(t.dni || '').trim() === dni
        );
        if (teacherByDni) {
          setTeacherModalData(teacherByDni);
          setShowTeacherModal(true);
          return;
        }
      }
      /* Sin ficha en colección teacher: igual se puede editar la cuenta en `user` (login). */
      showError(
        'No hay ficha en Docentes con el mismo email/DNI. Se abre la edición de la cuenta de usuario.',
        'warning'
      );
      setDataRow(full || null);
      setShow(true);
      return;
    }

    if (full && !full.source && roleUpper === 'ALUMNO') {
      if (!canEditDocenteOAlumnoDesdeUsuarios(sessionType)) {
        showError(
          'Solo un administrador o usuario administrativo puede editar alumnos desde este listado.',
          'warning'
        );
        return;
      }
      const studentsList = await getStudents(globalDispatch);
      const em = String(full.email || '').trim().toLowerCase();
      const dni = String(full.dni || '').trim();
      if (em) {
        const student = (studentsList || []).find(
          (s) => String(s.email || '').trim().toLowerCase() === em
        );
        if (student) {
          setStudentModalData(student);
          setShowStudentModal(true);
          return;
        }
      }
      if (dni) {
        const studentByDni = (studentsList || []).find(
          (s) => String(s.dni || '').trim() === dni
        );
        if (studentByDni) {
          setStudentModalData(studentByDni);
          setShowStudentModal(true);
          return;
        }
      }
      showError(
        'No hay ficha en Alumnos con el mismo email/DNI. Se abre la edición de la cuenta de usuario.',
        'warning'
      );
      setDataRow(full || null);
      setShow(true);
      return;
    }

    setDataRow(full || null);
    setShow(true);
  };

  const tableData = useMemo(
    () =>
      (dataUsers || []).map((u) => {
        const row = {
          nombre: u.name || '—',
          email: u.email || '',
          dni: u.dni || '',
          telefono: u.phone || '',
          direccion: u.address || '',
          rol: u.role || '',
          turno: u.shiftType || '—',
          estado: u.active !== false
        };
        row[USER_TABLE_ROW_ID] = u.id;
        return row;
      }),
    [dataUsers]
  );
  const filteredTableData = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return tableData;
    return tableData.filter((r) => String(r.nombre || '').toLowerCase().includes(q));
  }, [tableData, appliedSearch]);

  return (
    <React.Fragment>
      <Container fluid className={styles.container}>
        {error}
        <Row>
          <h1>Usuarios</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
            <Button className="w-100" variant="primary" size="xs" onClick={() => addUserEvent()}>
              Agregar
            </Button>
          </div>
          <div className={listToolbar.toolbarHalf}>
            <InputGroup className="w-100">
              <Form.Control
                type="text"
                value={searchText}
                onChange={(e) => {
                  const v = e.target.value;
                  setSearchText(v);
                  if (!String(v).trim()) setAppliedSearch('');
                }}
                placeholder="Buscar por nombre"
              />
              <Button type="button" className={listToolbar.buscarBtn} onClick={() => setAppliedSearch(searchText)}>
                Buscar
              </Button>
            </InputGroup>
          </div>
        </div>
        <br />
        <Row className={styles.tableRowWrap}>
          <Col xs={12} className="px-2 px-md-3 min-w-0">
            {filteredTableData.length > 0 ? (
              <Table
                key="user"
                tableEvents={(e, d) => tableEvents(e, d)}
                data={filteredTableData}
                actions="e"
                wideKeys={USER_TABLE_WIDE_KEYS}
              />
            ) : (
              <p className="text-muted">No hay usuarios en esta lista.</p>
            )}
          </Col>
        </Row>
      </Container>
      <Modal
        show={show}
        handleClose={closeUserEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow}
        roles={globalState.roles || []}
        shifts={globalState.shifts || []}
        degrees={globalState.degrees || []}
        changeActive={async (payload) => {
          const result = await changeActiveUser(globalDispatch, payload);
          buildNotification(result);
          return result;
        }}
      />

      <ModalTeacher
        show={showTeacherModal}
        handleClose={closeTeacherModal}
        saveEvent={(e) => saveTeacherFromUsers(e)}
        data={teacherModalData}
        degrees={globalState.degrees || []}
        titulosHabilitantes={titulosHabilitantes}
        modalidades={modalidadesDocente}
        changeActive={async (payload) => {
          const result = await changeActiveTeacher(globalDispatch, payload);
          buildNotification(result);
          if (Number(result.code) === 200) {
            await getUsersByRole(globalDispatch);
          }
          return result;
        }}
      />

      <ModalStudent
        show={showStudentModal}
        handleClose={closeStudentModal}
        saveEvent={(e) => saveStudentFromUsers(e)}
        data={studentModalData}
        degrees={globalState.degrees || []}
        changeActive={async (payload) => {
          const result = await changeActiveStudent(globalDispatch, payload);
          buildNotification(result);
          if (Number(result.code) === 200) {
            await getUsersByRole(globalDispatch);
          }
          return result;
        }}
      />
    </React.Fragment>
  );
};

export default Users;
