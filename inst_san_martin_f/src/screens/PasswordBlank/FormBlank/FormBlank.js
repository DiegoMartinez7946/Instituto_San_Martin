import React, { useState } from 'react';
import { Form, Button, InputGroup } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

import { isSaveSuccess } from '../../../utils/saveResult';
import { validateCorreoElectronicoRequerido } from '../../../utils/contact';

const FormBlank = ({ saveData }) => {
  const [emailError, setEmailError] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [showNewPassword, setShowNewPassword] = useState(false);
  const [data, setData] = useState({
    email: '',
    newPassword: ''
  });

  const handleInputChange = (event) => {
    if (event.target.name === 'email') {
      setEmailError('');
    }
    if (event.target.name === 'newPassword') {
      setPasswordError('');
    }
    setData({
      ...data,
      [event.target.name]: event.target.value
    });
  };

  const sendData = async (e) => {
    e.preventDefault();
    const emailCheck = validateCorreoElectronicoRequerido(data.email);
    if (!emailCheck.ok) {
      setEmailError(emailCheck.message);
      return;
    }
    setEmailError('');
    const pwd = String(data.newPassword || '');
    if (pwd.length < 6) {
      setPasswordError('La contraseña debe tener al menos 6 caracteres');
      return;
    }
    setPasswordError('');
    const res = await saveData({ ...data, email: String(data.email).trim(), newPassword: pwd });
    if (isSaveSuccess(res)) {
      setData({
        email: '',
        newPassword: ''
      });
      setShowNewPassword(false);
    }
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Label>Email</Form.Label>
        <Form.Control
          type="email"
          placeholder="Ingrese email"
          name="email"
          autoFocus
          onChange={handleInputChange}
          value={data.email}
          isInvalid={!!emailError}
          autoComplete="username"
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPassword">
        <Form.Label>Nueva contraseña</Form.Label>
        <InputGroup hasValidation>
          <Form.Control
            type={showNewPassword ? 'text' : 'password'}
            placeholder="Mínimo 6 caracteres"
            name="newPassword"
            onChange={handleInputChange}
            value={data.newPassword}
            isInvalid={!!passwordError}
            autoComplete="new-password"
          />
          <Button
            variant="outline-secondary"
            type="button"
            title={showNewPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            aria-label={showNewPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            onClick={() => setShowNewPassword((v) => !v)}
            tabIndex={-1}
          >
            <FontAwesomeIcon icon={showNewPassword ? faEyeSlash : faEye} />
          </Button>
          <Form.Control.Feedback type="invalid">{passwordError}</Form.Control.Feedback>
        </InputGroup>
        <Form.Text className="text-muted">Se guarda cifrada en la base de datos.</Form.Text>
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        Actualizar
      </Button>
    </Form>
  );
};

export default FormBlank;
