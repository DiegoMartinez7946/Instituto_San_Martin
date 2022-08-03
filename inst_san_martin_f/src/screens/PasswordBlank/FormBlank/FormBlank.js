import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

const FormBlank = ({ saveData }) => {

  const [data, setData] = useState({
    email: '',
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
        />
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
