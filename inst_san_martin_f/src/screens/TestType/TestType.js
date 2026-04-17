import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalTestType from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getTestType, addTestType, deleteTestType, updateTestType } from '../../context/Global/actions/TestTypeActions';

import styles from './TestType.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const TestType = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataTestTypes, setDataTestTypes] = useState([]);
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

  const getAllTestTypes = async () => {
    const result = await getTestType(globalDispatch);
    setDataTestTypes(result);
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
    setDataTestTypes(getAllTestTypes(globalDispatch));
  }, []);

  useEffect(() => {
    setDataTestTypes(globalState.testtypes);
  }, [globalState]);

  const saveEventHandler = async (e) => {
    const result =
      e.ID === '' ? await addTestType(globalDispatch, e) : await updateTestType(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    setDataTestTypes(getAllTestTypes(globalDispatch));
    setShow((current) => !current);
    setDataRow('');
    return result;
  };

  const addTestTypeEvent = () => {
    setDataRow('');
    setShow((current) => !current);
  };

  const closeTestTypeEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => { 
    if (e === 'edit') {
      setDataRow(d);
      setShow(current => !current);
    }
    if (e === 'delete') {
      const result = await deleteTestType(globalDispatch, d);
      buildNotification(result); 
    }
  }
  const filteredTestTypes = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return dataTestTypes;
    return (dataTestTypes || []).filter((t) => String(t.type || t.name || '').toLowerCase().includes(q));
  }, [dataTestTypes, appliedSearch]);

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Tipo de Examen</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addTestTypeEvent()}
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
              key={'testtype'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filteredTestTypes}
              actions={'ed'}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalTestType 
        show={show}
        handleClose={closeTestTypeEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default TestType;
