import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import { NIVELES_ORDENADOS } from '../../../constant/nivelesAcademicos';

const FormDegree = ({ dataEntry, saveData }) => {

  const [data, setData] = useState({
    id: '',
    name: '',
    nivel: '',
    resolucionId: ''
  });

  useEffect(() => {
      setData({
        id: dataEntry.id || dataEntry.ID || '',
        name: dataEntry.name || '',
        nivel: (dataEntry.nivel || '').toLowerCase(),
        resolucionId: dataEntry.resolucionId != null ? String(dataEntry.resolucionId) : ''
      }); 
  }, [dataEntry]);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setData((prev) => ({
      ...prev,
      [name]: value
    }));
  };

  const sendData = (e) => {
    e.preventDefault();
    saveData(data); 
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicID">
        <Form.Control type="hidden" disabled name="id" value={data ? data.id : ''}/>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicName">
        <Form.Label>Nombre</Form.Label>
        <Form.Control type="text" placeholder="Ingrese nombre" name="name" onChange={handleInputChange} value={data ? data.name : ''} autoFocus />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formDegreeNivel">
        <Form.Label>Nivel</Form.Label>
        <Form.Select name="nivel" value={data.nivel} onChange={handleInputChange} required>
          <option value="">Seleccione nivel</option>
          {NIVELES_ORDENADOS.map((n) => (
            <option key={n} value={n}>
              {n.charAt(0).toUpperCase() + n.slice(1)}
            </option>
          ))}
        </Form.Select>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formDegreeResolucion">
        <Form.Label>Resolución ID</Form.Label>
        <Form.Control
          type="text"
          name="resolucionId"
          placeholder="Texto alfanumérico (puede incluir caracteres especiales)"
          onChange={handleInputChange}
          value={data.resolucionId}
          required
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100">
        { (dataEntry && (dataEntry.id || dataEntry.ID)) ? 'Actualizar' : 'Guardar' }
      </Button>
    </Form>
  )
}

export default FormDegree;
