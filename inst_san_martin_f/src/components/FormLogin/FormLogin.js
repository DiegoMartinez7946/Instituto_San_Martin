import React, { useState } from 'react';
import { Form, Button } from 'react-bootstrap';

const Login = ({ login }) => {

  const [data, setData] = useState({
    email: '',
    password: ''
  });

  const handleInputChange = (event) => {
    setData({
      ...data,          
      [event.target.name] : event.target.value
    });
  };

  const sendData = (e) => {
    e.preventDefault();
    login(data); 
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Label>Email</Form.Label>
        <Form.Control type="email" placeholder="Ingrese email" name="email" onChange={handleInputChange}/>
        <Form.Text className="text-muted">
          Nunca divulgaremos su email con nadie.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPassword">
        <Form.Label>Password</Form.Label>
        <Form.Control type="password" placeholder="Ingrese password" name="password" onChange={handleInputChange}/>
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        Ingresar
      </Button>
    </Form>
  )
}

export default Login;
