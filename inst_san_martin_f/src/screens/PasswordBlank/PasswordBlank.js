import React, { useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import FormBlank from './FormBlank/FormBlank';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { blankPassword } from '../../context/Global/actions/UserActions';

import styles from './PasswordBlank.module.css';

const PasswordBlank = () => {

  const [, globalDispatch] = useGlobal();
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

  const eventHandler = async (e) => {
    const result = await blankPassword(globalDispatch, e);
    buildNotification(result);
  };

  return (
    <Container className={styles.container}>
      { error }
      <Row>
        <h1>Blanqueo de Password</h1>
      </Row>
      <hr />
      <br />
      <Row className="justify-content-center">
        <Col xs lg="4">
          <FormBlank
            saveData={(e) => eventHandler(e)}
          />
        </Col>
      </Row>
    </Container>
  );
};

export default PasswordBlank;
