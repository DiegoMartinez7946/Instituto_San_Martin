import React, { useState, useEffect, useMemo } from 'react';
import { Container, Row, Col, Button, Form, InputGroup } from 'react-bootstrap';
import Modal from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getUsersByRole, addUser, updateUser, changeActiveUser } from '../../context/Global/actions/UserActions';
import { getRole } from '../../context/Global/actions/RoleActions';
import { getShift } from '../../context/Global/actions/ShiftActions';
import { getDegree } from '../../context/Global/actions/DegreeActions';

import styles from './Users.module.css';
import listToolbar from '../common/ListToolbar.module.css';

/* Texto largo o varios turnos concatenados */
const USER_TABLE_WIDE_KEYS = ['nombre', 'email', 'rol', 'turno'];

/* Id interno por fila (no se muestra como columna: Table usa Object.keys). */
const USER_TABLE_ROW_ID = Symbol('userTableRowId');

const Users = () => {
  const [globalState, globalDispatch] = useGlobal();
  const [dataUsers, setDataUsers] = useState([]);
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);

  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');

  useEffect(() => {
    const load = async () => {
      await getRole(globalDispatch);
      await getShift(globalDispatch);
      await getDegree(globalDispatch);
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

  const tableEvents = async (ev, row) => {
    if (ev !== 'edit') return;
    const rowId = row[USER_TABLE_ROW_ID];
    const full = (dataUsers || []).find((u) => u.id === rowId);
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
    </React.Fragment>
  );
};

export default Users;
