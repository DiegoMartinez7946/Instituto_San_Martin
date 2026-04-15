import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

import { validateCorreoElectronicoRequerido } from '../../utils/contact';

const Login = ({ login }) => {

  const [emailError, setEmailError] = useState('');
  const [data, setData] = useState({
    email: '',
    password: ''
  });

  const handleInputChange = (event) => {
    if (event.target.name === 'email') {
      setEmailError('');
    }
    setData({
      ...data,          
      [event.target.name] : event.target.value
    });
  };

  const sendData = (e) => {
    e.preventDefault();
    const emailCheck = validateCorreoElectronicoRequerido(data.email);
    if (!emailCheck.ok) {
      setEmailError(emailCheck.message);
      return;
    }
    setEmailError('');
    login({ ...data, email: String(data.email).trim() }); 
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Label>Email</Form.Label>
        <Form.Control
          type="email"
          placeholder="Ingrese email"
          name="email"
          value={data.email}
          onChange={handleInputChange}
          isInvalid={!!emailError}
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
        <Form.Text className="text-muted">
          Nunca divulgaremos su email con nadie.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPassword">
        <Form.Label>Password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Ingrese password"
          name="password"
          value={data.password}
          onChange={handleInputChange}
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        Ingresar
      </Button>
    </Form>
  )
}

export default Login;
