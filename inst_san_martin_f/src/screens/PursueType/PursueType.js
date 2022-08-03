import React, { useState, useEffect } from 'react';
import { Button, Container, Row, Col } from 'react-bootstrap';
import ModalPursueType from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getPursueType, addPursueType, deletePursueType, updatePursueType } from '../../context/Global/actions/PursueTypeActions';

import styles from './PursueType.module.css';

const PursueType = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataPursueTypes, setDataPursueTypes] = useState([]);
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
    const result = e.ID === '' ?
      await addPursueType(globalDispatch, e) :
      await updatePursueType(globalDispatch, e);
   
    buildNotification(result); 
    setDataPursueTypes(getAllPursueTypes(globalDispatch));
    setShow(current => !current);
    setDataRow('');
  };

  const addPursueTypeEvent = () => {
    setShow(current => !current);
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

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Modalidades</h1>
        </Row>
        <hr />
        <br />
        <Row className="justify-content-center">
          <Col xs lg="4">
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addPursueTypeEvent()}
              >
                Agregar
              </Button>
          </Col>
        </Row>
        <br />  
        <Row>
          <Col xs lg="2"></Col>
          <Col>
            <Table
              key={'pursuetype'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={dataPursueTypes}
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
