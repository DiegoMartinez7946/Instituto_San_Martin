import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalPursueType from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getPursueType, addPursueType, deletePursueType, updatePursueType } from '../../context/Global/actions/PursueTypeActions';

import styles from './PursueType.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const PursueType = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataPursueTypes, setDataPursueTypes] = useState([]);
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

  const getAllPursueTypes = async () => {
    const result = await getPursueType(globalDispatch);
    setDataPursueTypes(result);
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
    setDataPursueTypes(getAllPursueTypes(globalDispatch));
  }, []);

  useEffect(() => {
    setDataPursueTypes(globalState.pursuetypes);
  }, [globalState]);

  const saveEventHandler = async (e) => {
    const result =
      e.ID === '' ? await addPursueType(globalDispatch, e) : await updatePursueType(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    setDataPursueTypes(getAllPursueTypes(globalDispatch));
    setShow((current) => !current);
    setDataRow('');
    return result;
  };

  const addPursueTypeEvent = () => {
    setDataRow('');
    setShow((current) => !current);
  };

  const closePursueTypeEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => { 
    if (e === 'edit') {
      setDataRow(d);
      setShow(current => !current);
    }
    if (e === 'delete') {
      const result = await deletePursueType(globalDispatch, d);
      buildNotification(result); 
    }
  }
  const filteredPursueTypes = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return dataPursueTypes;
    return (dataPursueTypes || []).filter((t) => String(t.type || t.name || '').toLowerCase().includes(q));
  }, [dataPursueTypes, appliedSearch]);

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Modalidades</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addPursueTypeEvent()}
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
              key={'pursuetype'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filteredPursueTypes}
              actions={'ed'}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalPursueType 
        show={show}
        handleClose={closePursueTypeEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default PursueType;
