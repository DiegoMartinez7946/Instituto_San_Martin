import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalDegree from './Modal/ModalDegree';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';
import ConfirmChangeEstadoModal from '../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree, addDegree, changeActiveDegree, updateDegree } from '../../context/Global/actions/DegreeActions';

import styles from './Degree.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const Degree = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataDegrees, setDataDegrees] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [error, setError] = useState(null);
  const [estadoConfirm, setEstadoConfirm] = useState(null);
  const [estadoConfirmLoading, setEstadoConfirmLoading] = useState(false);

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

  const getAllDegrees = async () => {
    const result = await getDegree(globalDispatch);
    setDataDegrees(result);
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
    setDataDegrees(getAllDegrees(globalDispatch));
  }, []);

  useEffect(() => {
    setDataDegrees(globalState.degrees);
  }, [globalState]);

  const saveEventHandler = async (e) => {
    const result = !(e.id || e.ID)
      ? await addDegree(globalDispatch, e)
      : await updateDegree(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    setDataDegrees(getAllDegrees(globalDispatch));
    setShow((current) => !current);
    setDataRow('');
    return result;
  };

  const addDegreeEvent = () => {
    setDataRow('');
    setShow((current) => !current);
  };

  const closeDegreeEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => {
    if (e === 'edit') {
      setDataRow(d);
      setShow((current) => !current);
    }
    if (e === 'check') {
      const cur = d.active === true;
      setEstadoConfirm({ row: d, fromActive: cur, toActive: !cur });
    }
  };

  const confirmEstadoToggle = async () => {
    if (!estadoConfirm) return;
    setEstadoConfirmLoading(true);
    const item = { ...estadoConfirm.row, active: estadoConfirm.toActive };
    const result = await changeActiveDegree(globalDispatch, item);
    buildNotification(result);
    setEstadoConfirmLoading(false);
    if (Number(result.code) === 200) {
      setEstadoConfirm(null);
    }
  };
  const filteredDegrees = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return dataDegrees;
    return (dataDegrees || []).filter((d) => String(d.name || '').toLowerCase().includes(q));
  }, [dataDegrees, appliedSearch]);

  return (
    <React.Fragment>
      <Container fluid className={styles.container}>
        { error }
        <Row>
          <h1>Carreras</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addDegreeEvent()}
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
        <Row className={styles.tableRowWrap}>
          <Col xs={12} className="px-2 px-md-3 min-w-0">
            <Table
              key={'degree'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filteredDegrees}
              actions={'ec'}
              wideKeys={[]}
            />
          </Col>
        </Row>
      </Container>
      <ModalDegree
        show={show}
        handleClose={closeDegreeEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
        changeActive={async (payload) => {
          const result = await changeActiveDegree(globalDispatch, payload);
          buildNotification(result);
          return result;
        }}
      />
      <ConfirmChangeEstadoModal
        show={!!estadoConfirm}
        onHide={() => !estadoConfirmLoading && setEstadoConfirm(null)}
        kind="carrera"
        itemName={estadoConfirm && estadoConfirm.row ? estadoConfirm.row.name : ''}
        fromActive={estadoConfirm ? estadoConfirm.fromActive : true}
        toActive={estadoConfirm ? estadoConfirm.toActive : true}
        onConfirm={confirmEstadoToggle}
        confirming={estadoConfirmLoading}
      />
    </React.Fragment>
  );
};

export default Degree;
