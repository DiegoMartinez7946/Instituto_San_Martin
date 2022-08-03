import React, { useState, useEffect } from 'react';
import { Button, Container, Row, Col } from 'react-bootstrap';
import ModalDegree from './Modal/ModalDegree';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree, addDegree, changeActiveDegree, updateDegree } from '../../context/Global/actions/DegreeActions';

import styles from './Degree.module.css';

const Degree = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [dataDegrees, setDataDegrees] = useState([]);
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
    const result = e.ID === '' ?
      await addDegree(globalDispatch, e) :
      await updateDegree(globalDispatch, e);
   
    buildNotification(result); 
    setDataDegrees(getAllDegrees(globalDispatch));
    setShow(current => !current);
    setDataRow('');
  };

  const addDegreeEvent = () => {
    setShow(current => !current);
  };

  const closeDegreeEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  const tableEvents = async (e, d) => { 
    if (e === 'edit') {
      setDataRow(d);
      setShow(current => !current);
    }
    if (e === 'check') {
      d.active = !d.active;
      const result = await changeActiveDegree(globalDispatch, d);
      buildNotification(result); 
    }
  }

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Carreras</h1>
        </Row>
        <hr />
        <br />
        <Row className="justify-content-center">
          <Col xs lg="4">
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addDegreeEvent()}
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
              key={'degree'}
              tableEvents={(e, d) => tableEvents(e, d)}
              data={dataDegrees}
              actions={'ec'}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalDegree 
        show={show}
        handleClose={closeDegreeEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default Degree;
