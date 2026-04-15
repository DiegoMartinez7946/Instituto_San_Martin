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
    if (result == null || result.code === undefined) {
      showError('No se recibió respuesta del servidor. Revise la conexión o el token.', 'danger');
      return;
    }
    const code = Number(result.code);
    switch (code) {
      case 199:
        showError(result.message || 'Advertencia', 'warning');
        break;
      case 200:
      case 201:
        showError(result.message || 'Operación correcta', 'success');
        break;
      case 400:
      case 404:
        showError(result.message || 'Error', 'danger');
        break;
      default:
        showError(result.message || 'Respuesta no reconocida', 'warning');
        break;
    }
  };

  const eventHandler = async (e) => {
    let result;
    try {
      result = await blankPassword(globalDispatch, e);
    } catch (err) {
      showError(
        (err && err.response && err.response.data && err.response.data.message) ||
          'Error de red al contactar el servidor.',
        'danger'
      );
      return null;
    }
    buildNotification(result);
    return result;
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
