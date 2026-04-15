import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import { validateCorreoElectronicoRequerido } from '../../../utils/contact';

const FormUser = ({ saveData }) => {

  const [emailError, setEmailError] = useState('');
  const [data, setData] = useState({
    email: '',
    userType: 'ADMINISTRATIVO',
    password: '',
  });

  useEffect(() => {
      setData({
        email: '',
        userType: 'ADMINISTRATIVO',
        password: ''
      }); 
  }, []);

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
      userType: 'ADMINISTRATIVO',
      password: ''
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
          value={data ? data.email : ''} 
          autoFocus
          onChange={handleInputChange}
          isInvalid={!!emailError}
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPassword">
        <Form.Label>Password</Form.Label>
        <Form.Control 
          type="text" 
          placeholder="Ingrese password" 
          name="password" 
          onChange={handleInputChange} 
          value={data ? data.password : ''} 
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        Guardar
      </Button>
    </Form>
  )
}

export default FormUser;
