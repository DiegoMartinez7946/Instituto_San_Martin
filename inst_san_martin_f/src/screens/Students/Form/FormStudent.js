import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

const degId = (d) => (d && (d.id || d.ID)) || '';

const FormStudent = ({ dataEntry, degrees, saveData }) => {

  const [data, setData] = useState({
    id: '',
    name: '',
    email: '',
    phone: '',
    dni: '',
    address: '',
    degreeIds: []
  });

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : null;
    setData({
      id: (entry && entry.id) || '',
      name: (entry && entry.name) || '',
      email: (entry && entry.email) || '',
      phone: (entry && entry.phone) || '',
      dni: (entry && entry.dni) || '',
      address: (entry && entry.address) || '',
      degreeIds: entry && Array.isArray(entry.degreeIds) ? [...entry.degreeIds] : []
    });
  }, [dataEntry]);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setData((prev) => ({
      ...prev,
      [name]: value
    }));
  };

  const toggleDegree = (hexId) => {
    setData(prev => {
      const set = new Set(prev.degreeIds || []);
      if (set.has(hexId)) {
        set.delete(hexId);
      } else {
        set.add(hexId);
      }
      return { ...prev, degreeIds: [...set] };
    });
  };

  const sendData = (e) => {
    e.preventDefault();
    saveData(data);
  };

  const selectableDegrees = (degrees || []).filter(d => {
    const hid = degId(d);
    if (d.active === true) return true;
    if ((data.degreeIds || []).includes(hid)) return true;
    return false;
  });

  return (
    <Form onSubmit={sendData}>
      <Form.Control type="hidden" name="id" value={data.id} readOnly />

      <Form.Group className="mb-3" controlId="studentIdAlumno">
        <Form.Label>ID de alumno</Form.Label>
        <Form.Control
          type="text"
          readOnly
          plaintext
          value={data.id ? data.id : 'Se asigna al guardar'}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentName">
        <Form.Label>Nombre</Form.Label>
        <Form.Control
          type="text"
          name="name"
          value={data.name}
          onChange={handleInputChange}
          required
          autoFocus
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentEmail">
        <Form.Label>Correo electrónico</Form.Label>
        <Form.Control
          type="email"
          name="email"
          value={data.email}
          onChange={handleInputChange}
          placeholder="nombre@ejemplo.com"
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentPhone">
        <Form.Label>Número de teléfono</Form.Label>
        <Form.Control
          type="text"
          name="phone"
          value={data.phone}
          onChange={handleInputChange}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentDni">
        <Form.Label>DNI</Form.Label>
        <Form.Control
          type="text"
          name="dni"
          value={data.dni}
          onChange={handleInputChange}
          required
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentAddress">
        <Form.Label>Dirección</Form.Label>
        <Form.Control
          type="text"
          name="address"
          value={data.address}
          onChange={handleInputChange}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentDegrees">
        <Form.Label>Carreras</Form.Label>
        <div className="border rounded p-2" style={{ maxHeight: '200px', overflowY: 'auto' }}>
          {selectableDegrees.length === 0 ? (
            <span className="text-muted">No hay carreras activas. Cree carreras primero.</span>
          ) : (
            selectableDegrees.map(d => {
              const id = degId(d);
              if (!id) return null;
              return (
                <Form.Check
                  key={id}
                  type="checkbox"
                  id={`deg-${id}`}
                  label={d.name}
                  checked={(data.degreeIds || []).includes(id)}
                  onChange={() => toggleDegree(id)}
                />
              );
            })
          )}
        </div>
      </Form.Group>

      <Button variant="primary" type="submit" className="w-100">
        {data.id ? 'Actualizar' : 'Guardar'}
      </Button>
    </Form>
  );
};

export default FormStudent;
