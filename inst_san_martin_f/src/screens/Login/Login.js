import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Container, Row, Col, Toast, ToastContainer } from 'react-bootstrap';

import FormLogin from '../../components/FormLogin/FormLogin';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { login } from '../../context/Global/actions/UserActions';

import styles from './Login.module.css';

const Login = () => {

  const [, globalDispatch] = useGlobal();
  const [notification, setNotification] = useState(null);
  const navigate = useNavigate();

  const showError = (message) => {
    setNotification(
      <ToastContainer position={'top-end'} >
        <Toast
            className="d-inline-block m-1"
            bg={'danger'}
            onClose={() => setNotification(null)} delay={3000} autohide
          >
            <Toast.Header>
              <strong className="me-auto">Error Login</strong>
            </Toast.Header>
            <Toast.Body className={'text-white'}>
              { message }
            </Toast.Body>
          </Toast>
      </ToastContainer>);
  };

  const loginEvent = async (e) => {
    const result = await login(globalDispatch, e);
    result.data.code === 400 ? showError(result.data.message) : navigate('/main');
  };

  return (
    <Container className={styles.container}>
      { notification }
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
