import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalRole from './Modal/ModalRole';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getRole, addRole, deleteRole, updateRole } from '../../context/Global/actions/RoleActions';

import styles from './Roles.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const Roles = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataRoles, setDataRoles] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [error, setError] = useState(null);

  const showError = (message, type) => {
    message ?
    setError(
      <Notification 
        message={message}
        type={type}
        show={showError}
      />  
    ) :
    setError(null);
  };

  const getAllRoles = async () => {
    const result = await getRole(globalDispatch);
    setDataRoles(result);
    return; 
  };

  const buildNotification = (result) => {
    switch (result.code) {
      case 199:
        showError(result.message, 'warning');
        break;
      case 200:
      case 201:
        showError(result.message, 'success');
        break;
      case 400:
        showError(result.message, 'danger');
        break; 
      default:
        break;
    }
  };

  useEffect(() => {
    setDataRoles(getAllRoles(globalDispatch));
  }, []);

  useEffect(() => {
    setDataRoles(globalState.roles);
  }, [globalState]);

  const saveEventHandler = async (e) => {
    const result =
      e.ID === '' ? await addRole(globalDispatch, e) : await updateRole(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    setDataRoles(getAllRoles(globalDispatch));
    setShow((current) => !current);
    setDataRow('');
    return result;
  };

  const addRoleEvent = () => {
    setDataRow('');
    setShow((current) => !current);
  };

  const closeRoleEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => { 
    if (e === 'edit') {
      setDataRow(d);
      setShow(current => !current);
    }
    if (e === 'delete') {
      const result = await deleteRole(globalDispatch, d);
      buildNotification(result); 
    }
  }
  const filteredRoles = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return dataRoles;
    return (dataRoles || []).filter((r) => String(r.type || r.name || '').toLowerCase().includes(q));
  }, [dataRoles, appliedSearch]);

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Roles</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addRoleEvent()}
              >
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
        <Row>
          <Col xs lg="2"></Col>
          <Col className="min-w-0">
            <Table
              key={'role'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filteredRoles}
              actions={'ed'}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalRole 
        show={show}
        handleClose={closeRoleEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default Roles;
