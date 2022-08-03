import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

const FormUser = ({ saveData }) => {

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
        />
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
