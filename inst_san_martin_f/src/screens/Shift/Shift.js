import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalShift from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getShift, addShift, deleteShift, updateShift } from '../../context/Global/actions/ShiftActions';

import styles from './Shift.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const Shift = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataShifts, setDataShifts] = useState([]);
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

  const getAllShifts = async () => {
    const result = await getShift(globalDispatch);
    setDataShifts(result);
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
    setDataShifts(getAllShifts(globalDispatch));
  }, []);

  useEffect(() => {
    setDataShifts(globalState.shifts);
  }, [globalState]);

  const saveEventHandler = async (e) => {
    const result =
      e.ID === '' ? await addShift(globalDispatch, e) : await updateShift(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    setDataShifts(getAllShifts(globalDispatch));
    setShow((current) => !current);
    setDataRow('');
    return result;
  };

  const addShiftEvent = () => {
    setDataRow('');
    setShow((current) => !current);
  };

  const closeShiftEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => { 
    if (e === 'edit') {
      setDataRow(d);
      setShow(current => !current);
    }
    if (e === 'delete') {
      const result = await deleteShift(globalDispatch, d);
      buildNotification(result); 
    }
  }
  const filteredShifts = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return dataShifts;
    return (dataShifts || []).filter((s) => String(s.type || s.name || '').toLowerCase().includes(q));
  }, [dataShifts, appliedSearch]);

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Turnos de Cursado</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addShiftEvent()}
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
              key={'shift'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filteredShifts}
              actions={'ed'}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalShift 
        show={show}
        handleClose={closeShiftEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default Shift;
