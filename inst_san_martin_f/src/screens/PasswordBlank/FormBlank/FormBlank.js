import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

import { validateCorreoElectronicoRequerido } from '../../../utils/contact';

const FormBlank = ({ saveData }) => {

  const [emailError, setEmailError] = useState('');
  const [data, setData] = useState({
    email: '',
    newPassword: '',
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
    saveData({ ...data, email: String(data.email).trim() }); 
    setData({
      email: '',
      newPassword: ''
    });
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
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPassword">
        <Form.Label>Password</Form.Label>
        <Form.Control 
          type="text" 
          placeholder="Ingrese password" 
          name="newPassword" 
          onChange={handleInputChange} 
          value={data.newPassword}
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        Actualizar
      </Button>
    </Form>
  )
}

export default FormBlank;
