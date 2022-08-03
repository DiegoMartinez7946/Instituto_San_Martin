import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';
import lodash from 'lodash';

const FormPursueType = ({ dataEntry, saveData }) => {

  const [data, setData] = useState({
    ID: '',
    type: ''
  });

  useEffect(() => {
      setData({
        ID: dataEntry.ID || '',
        type: dataEntry.type || ''
      }); 
  }, []);

  const handleInputChange = (event) => {
    setData({
      ...data,          
      [event.target.name] : lodash.toUpper(event.target.value)
    });
  };

  const sendData = (e) => {
    e.preventDefault();
    saveData(data); 
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicID">
        <Form.Control type="hidden" disabled name="ID" value={data ? data.ID : ''}/>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicType">
        <Form.Label>Modalidad</Form.Label>
        <Form.Control type="text" placeholder="Ingrese modalidad" name="type" onChange={handleInputChange} value={data ? data.type : ''} autoFocus />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        { dataEntry ? 'Actualizar' : 'Guardar' }
      </Button>
    </Form>
  )
}

export default FormPursueType;
