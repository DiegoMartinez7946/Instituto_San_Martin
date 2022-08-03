import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Row, Col } from 'react-bootstrap';

import FormLogin from '../../components/FormLogin/FormLogin';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { login } from '../../context/Global/actions/UserActions';

import styles from './Login.module.css';

const Login = () => {

  const [, globalDispatch] = useGlobal();
  const [error, setError] = useState(null);
  const navigate = useNavigate();

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


  const loginEvent = async (e) => {
    const result = await login(globalDispatch, e);
    result.data.code === 400 ? showError(result.data.message, 'danger') : navigate('/main');
  };

  return (
    <Container className={styles.container}>
      { error }
      <Row>
        <h1>Login</h1>
      </Row>
      <Row className="justify-content-center">
        <Col xs lg="4">
          <FormLogin 
            login={(e) => loginEvent(e)}
          />
        </Col>
      </Row>
    </Container>
  );
};

export default Login;
