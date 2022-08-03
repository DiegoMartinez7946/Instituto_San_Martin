import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

const FormChange = ({ email, saveData }) => {

  const [data, setData] = useState({
    email,
    currentPassword: '',
    newPassword: '',
  });

  const handleInputChange = (event) => {
    setData({
      ...data,          
      [event.target.name] : event.target.value
    });
  };

  const sendData = (e) => {
    e.preventDefault();
    saveData(data); 
    setData({
      email: '',
      currentPassword: '',
      newPassword: ''
    });
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Control 
          type="hidden" 
          name="email" 
          value={data.email}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPasswordCurrent">
        <Form.Label>Password Actual</Form.Label>
        <Form.Control 
          autoFocus
          type="text" 
          placeholder="Ingrese password actual" 
          name="currentPassword" 
          onChange={handleInputChange} 
          value={data.currentPassword}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPasswordNew">
        <Form.Label>Password Nueva</Form.Label>
        <Form.Control 
          type="text" 
          placeholder="Ingrese nueva password" 
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

export default FormChange;
