import React, { useState, useEffect } from 'react';
import { Container, Row, Col, Button } from 'react-bootstrap';
import Modal from './Modal/Modal';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getUsersByRole, addUser } from '../../context/Global/actions/UserActions';

import styles from './Users.module.css';

const Users = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [dataUsers, setDataUsers] = useState([]);
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [error, setError] = useState(null);

  useEffect(() => {
    setDataUsers(getUsersByRole(globalDispatch));
  }, []);

  useEffect(() => {
    setDataUsers(globalState.users);
  }, [globalState]);

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

  const saveEventHandler = async (e) => {
    const result = await addUser(globalDispatch, e);
   
    console.log('result', result);
    buildNotification(result); 
    setDataUsers(getUsersByRole(globalDispatch));
    setShow(current => !current);
    setDataRow('');
  };

  const addUserEvent = () => {
    setShow(current => !current);
  };

  const closeUserEvent = () => {
    setDataRow('');
    setShow(current => !current);
  };

  return (
    <React.Fragment>
      <Container className={styles.container}>
        { error }
        <Row>
          <h1>Usuarios</h1>
        </Row>
        <hr />
        <br />
        <Row className="justify-content-center">
          <Col xs lg="4">
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addUserEvent()}
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
              key={'user'}
              tableEvents={null}
              data={dataUsers}
              actions={''}
            /> 
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <Modal 
        show={show}
        handleClose={closeUserEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
      />
    </React.Fragment>
  );
};

export default Users;
