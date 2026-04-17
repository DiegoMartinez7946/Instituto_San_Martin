import React, { useState } from 'react';
import { Container, Row, Col } from 'react-bootstrap';
import FormChange from './FormChange/FormChange';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { changePassword } from '../../context/Global/actions/UserActions';
import { decodeTokenFromCookie } from '../../utils/jwt';

import styles from './PasswordChange.module.css';

const PasswordChange = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [error, setError] = useState(null);

  const { userLogin } = globalState;

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
      case 403:
        showError(result.message, 'danger');
        break; 
      default:
        break;
    }
  };

  const eventHandler = async (e) => {
    const result = await changePassword(globalDispatch, e);
    buildNotification(result);
    return result;
  };

  return (
    <Container className={styles.container}>
      { error }
      <Row>
        <h1>Cambiar Password</h1>
      </Row>
      <hr />
      <br />
      <Row className="justify-content-center">
        <Col xs lg="4">
        {(() => {
          const decoded = decodeTokenFromCookie();
          const email =
            (userLogin && userLogin.email) ||
            (decoded && (decoded.email || decoded.Email)) ||
            '';
          return email ? (
            <FormChange saveData={(e) => eventHandler(e)} email={email} />
          ) : null;
        })()}
        </Col>
      </Row>
    </Container>
  );
};

export default PasswordChange;
