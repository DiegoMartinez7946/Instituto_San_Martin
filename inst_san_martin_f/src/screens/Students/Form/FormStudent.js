import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import { NIVELES_ORDENADOS, jerarquiaDeNivel, etiquetaNivel } from '../../../constant/nivelesAcademicos';
import { validateDNI } from '../../../utils/dni';
import {
  validateCorreoElectronicoOpcional,
  validateTelefonoOpcional,
  sanitizeTelefonoInput
} from '../../../utils/contact';

const degId = (d) => (d && (d.id || d.ID)) || '';

const FormStudent = ({ dataEntry, degrees, saveData }) => {

  const [dniError, setDniError] = useState('');
  const [emailError, setEmailError] = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [data, setData] = useState({
    id: '',
    name: '',
    email: '',
    phone: '',
    dni: '',
    address: '',
    nivelAprobado: '',
    degreeIds: []
  });

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : null;
    setData({
      id: (entry && entry.id) || '',
      name: (entry && entry.name) || '',
      email: (entry && entry.email) || '',
      phone: sanitizeTelefonoInput((entry && entry.phone) || ''),
      dni: (entry && entry.dni) || '',
      address: (entry && entry.address) || '',
      nivelAprobado: (entry && entry.nivelAprobado) ? String(entry.nivelAprobado).toLowerCase() : '',
      degreeIds: entry && Array.isArray(entry.degreeIds) ? [...entry.degreeIds] : []
    });
  }, [dataEntry]);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    if (name === 'dni') {
      setDniError('');
    }
    if (name === 'email') {
      setEmailError('');
    }
    setData((prev) => {
      const next = { ...prev, [name]: value };
      if (name === 'nivelAprobado' && (degrees || []).length) {
        const ha = jerarquiaDeNivel(value);
        const kept = (prev.degreeIds || []).filter((hexId) => {
          const d = (degrees || []).find((x) => degId(x) === hexId);
          if (!d) return false;
          const hn = jerarquiaDeNivel((d.nivel || '').toLowerCase().trim());
          if (ha == null || hn == null) return false;
          return ha >= hn - 1;
        });
        next.degreeIds = kept;
      }
      return next;
    });
  };

  const handlePhoneChange = (event) => {
    setPhoneError('');
    const digits = sanitizeTelefonoInput(event.target.value);
    setData((prev) => ({ ...prev, phone: digits }));
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

  const haAlumno = jerarquiaDeNivel(data.nivelAprobado);

  const degreeCheckboxDisabled = (d) => {
    if (!data.nivelAprobado) return true;
    const nivelCarrera = (d.nivel || '').toLowerCase().trim();
    if (!nivelCarrera) return true;
    const hn = jerarquiaDeNivel(nivelCarrera);
    if (hn == null || haAlumno == null) return true;
    return haAlumno < hn - 1;
  };

  const degreeCheckboxTitle = (d) => {
    if (!data.nivelAprobado) return 'Seleccione primero el nivel aprobado del alumno';
    if (!(d.nivel || '').trim()) return 'La carrera no tiene nivel configurado; edítela en Carreras';
    if (degreeCheckboxDisabled(d)) return 'El nivel aprobado no alcanza para inscribir en esta carrera';
    return '';
  };

  const sendData = (e) => {
    e.preventDefault();
    setEmailError('');
    setPhoneError('');
    const emailCheck = validateCorreoElectronicoOpcional(data.email);
    if (!emailCheck.ok) {
      setEmailError(emailCheck.message);
      return;
    }
    const phoneCheck = validateTelefonoOpcional(data.phone);
    if (!phoneCheck.ok) {
      setPhoneError(phoneCheck.message);
      return;
    }
    const dniCheck = validateDNI(data.dni);
    if (!dniCheck.ok) {
      setDniError(dniCheck.message);
      return;
    }
    setDniError('');
    saveData({
      ...data,
      dni: String(data.dni).trim(),
      email: String(data.email).trim(),
      phone: data.phone
    });
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
          isInvalid={!!emailError}
          placeholder="nombre@ejemplo.com"
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentPhone">
        <Form.Label>Número de teléfono</Form.Label>
        <Form.Control
          type="text"
          name="phone"
          inputMode="numeric"
          autoComplete="tel"
          maxLength={15}
          value={data.phone}
          onChange={handlePhoneChange}
          isInvalid={!!phoneError}
        />
        <Form.Control.Feedback type="invalid">{phoneError}</Form.Control.Feedback>
        <Form.Text className="text-muted">Solo números, 7 a 15 dígitos. Deje vacío si no aplica.</Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentDni">
        <Form.Label>DNI</Form.Label>
        <Form.Control
          type="text"
          name="dni"
          inputMode="numeric"
          autoComplete="off"
          maxLength={8}
          value={data.dni}
          onChange={handleInputChange}
          isInvalid={!!dniError}
          required
        />
        <Form.Control.Feedback type="invalid">{dniError}</Form.Control.Feedback>
        <Form.Text className="text-muted">7 u 8 dígitos numéricos, sin puntos.</Form.Text>
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

      <Form.Group className="mb-3" controlId="studentNivelAprobado">
        <Form.Label>Nivel aprobado</Form.Label>
        <div className="border rounded p-2">
          {NIVELES_ORDENADOS.map((n) => (
            <Form.Check
              key={n}
              type="radio"
              id={`nivel-aprob-${n}`}
              name="nivelAprobado"
              value={n}
              label={etiquetaNivel(n)}
              checked={data.nivelAprobado === n}
              onChange={handleInputChange}
            />
          ))}
        </div>
        <Form.Text className="text-muted">
          Indica el máximo nivel cursado/aprobado; los niveles inferiores se consideran aprobados.
        </Form.Text>
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
              const dis = degreeCheckboxDisabled(d);
              return (
                <Form.Check
                  key={id}
                  type="checkbox"
                  id={`deg-${id}`}
                  label={d.name + (d.nivel ? ` (${etiquetaNivel(d.nivel)})` : '')}
                  title={degreeCheckboxTitle(d)}
                  checked={(data.degreeIds || []).includes(id)}
                  disabled={dis}
                  onChange={() => { if (!dis) toggleDegree(id); }}
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
