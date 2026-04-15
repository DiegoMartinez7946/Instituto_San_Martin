import React, { useState, useEffect, useMemo } from 'react';
import { Form, Button } from 'react-bootstrap';

import { NIVELES_ORDENADOS, etiquetaNivel } from '../../../constant/nivelesAcademicos';
import { validateDNI } from '../../../utils/dni';
import {
  validateCorreoElectronicoOpcional,
  validateTelefonoOpcional,
  sanitizeTelefonoInput
} from '../../../utils/contact';

const degId = (d) => (d && (d.id || d.ID)) || '';

const FormTeacher = ({
  dataEntry,
  degrees,
  titulosHabilitantes,
  modalidades,
  saveData
}) => {
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
    enseniaEn: [],
    degreeIds: [],
    tituloHabilitanteId: '',
    modalidadId: ''
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
      enseniaEn: entry && Array.isArray(entry.enseniaEn) ? [...entry.enseniaEn] : [],
      degreeIds: entry && Array.isArray(entry.degreeIds) ? [...entry.degreeIds] : [],
      tituloHabilitanteId: (entry && entry.tituloHabilitanteId) || '',
      modalidadId: (entry && entry.modalidadId) || ''
    });
  }, [dataEntry]);

  const tituloEsSi = useMemo(() => {
    const row = (titulosHabilitantes || []).find(
      (x) => (x.id || x.ID) === data.tituloHabilitanteId
    );
    return row && String(row.codigo).toUpperCase() === 'SI';
  }, [titulosHabilitantes, data.tituloHabilitanteId]);

  useEffect(() => {
    if (tituloEsSi) return;
    const mod = (modalidades || []).find((m) => (m.id || m.ID) === data.modalidadId);
    if (mod && String(mod.codigo).toUpperCase() === 'TITULAR') {
      const fallback = (modalidades || []).find((m) => String(m.codigo).toUpperCase() !== 'TITULAR');
      if (fallback) {
        setData((prev) => ({ ...prev, modalidadId: fallback.id || fallback.ID }));
      }
    }
  }, [tituloEsSi, modalidades, data.modalidadId]);

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    if (name === 'dni') {
      setDniError('');
    }
    if (name === 'email') {
      setEmailError('');
    }
    setData((prev) => ({ ...prev, [name]: value }));
  };

  const handlePhoneChange = (event) => {
    setPhoneError('');
    const digits = sanitizeTelefonoInput(event.target.value);
    setData((prev) => ({ ...prev, phone: digits }));
  };

  const toggleNivel = (nivel) => {
    setData((prev) => {
      const set = new Set(prev.enseniaEn || []);
      if (set.has(nivel)) set.delete(nivel);
      else set.add(nivel);
      return { ...prev, enseniaEn: [...set] };
    });
  };

  const toggleDegree = (hexId) => {
    setData((prev) => {
      const set = new Set(prev.degreeIds || []);
      if (set.has(hexId)) set.delete(hexId);
      else set.add(hexId);
      return { ...prev, degreeIds: [...set] };
    });
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

  const selectableDegrees = (degrees || []).filter((d) => {
    const hid = degId(d);
    if (d.active === true) return true;
    if ((data.degreeIds || []).includes(hid)) return true;
    return false;
  });

  const modalidadDisabled = (m) => {
    return String(m.codigo).toUpperCase() === 'TITULAR' && !tituloEsSi;
  };

  return (
    <Form onSubmit={sendData}>
      <Form.Control type="hidden" name="id" value={data.id} readOnly />

      <Form.Group className="mb-3" controlId="teacherId">
        <Form.Label>ID docente</Form.Label>
        <Form.Control type="text" readOnly plaintext value={data.id ? data.id : 'Se asigna al guardar'} />
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherName">
        <Form.Label>Nombre</Form.Label>
        <Form.Control type="text" name="name" value={data.name} onChange={handleInputChange} required autoFocus />
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherEmail">
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

      <Form.Group className="mb-3" controlId="teacherPhone">
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

      <Form.Group className="mb-3" controlId="teacherDni">
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

      <Form.Group className="mb-3" controlId="teacherAddress">
        <Form.Label>Dirección</Form.Label>
        <Form.Control type="text" name="address" value={data.address} onChange={handleInputChange} />
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherEnsenia">
        <Form.Label>Enseña en</Form.Label>
        <div className="border rounded p-2">
          {NIVELES_ORDENADOS.map((n) => (
            <Form.Check
              key={n}
              type="checkbox"
              id={`ensenia-${n}`}
              label={etiquetaNivel(n)}
              checked={(data.enseniaEn || []).includes(n)}
              onChange={() => toggleNivel(n)}
            />
          ))}
        </div>
        <Form.Text className="text-muted">Puede marcar más de un nivel.</Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherTitulo">
        <Form.Label>Título habilitante</Form.Label>
        <div className="border rounded p-2">
          {(titulosHabilitantes || []).map((t) => {
            const tid = t.id || t.ID;
            if (!tid) return null;
            return (
              <Form.Check
                key={tid}
                type="radio"
                name="tituloHabilitanteId"
                id={`titulo-${tid}`}
                value={tid}
                label={String(t.codigo).toUpperCase() === 'SI' ? 'Sí' : 'No'}
                checked={data.tituloHabilitanteId === tid}
                onChange={handleInputChange}
              />
            );
          })}
        </div>
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherModalidad">
        <Form.Label>Modalidad</Form.Label>
        <div className="border rounded p-2">
          {(modalidades || []).map((m) => {
            const mid = m.id || m.ID;
            if (!mid) return null;
            const dis = modalidadDisabled(m);
            return (
              <Form.Check
                key={mid}
                type="radio"
                name="modalidadId"
                id={`mod-${mid}`}
                value={mid}
                label={etiquetaModalidad(m.codigo)}
                checked={data.modalidadId === mid}
                disabled={dis}
                title={dis ? 'Titular solo con título habilitante Sí' : ''}
                onChange={handleInputChange}
              />
            );
          })}
        </div>
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherDegrees">
        <Form.Label>Carreras</Form.Label>
        <div className="border rounded p-2" style={{ maxHeight: '200px', overflowY: 'auto' }}>
          {selectableDegrees.length === 0 ? (
            <span className="text-muted">No hay carreras activas. Cree carreras primero.</span>
          ) : (
            selectableDegrees.map((d) => {
              const id = degId(d);
              if (!id) return null;
              return (
                <Form.Check
                  key={id}
                  type="checkbox"
                  id={`deg-${id}`}
                  label={d.name + (d.nivel ? ` (${etiquetaNivel(d.nivel)})` : '')}
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

function etiquetaModalidad(codigo) {
  const c = String(codigo || '').toLowerCase();
  if (c === 'titular') return 'Titular';
  if (c === 'provisional') return 'Provisional';
  if (c === 'suplente') return 'Suplente';
  return codigo || '';
}

export default FormTeacher;
