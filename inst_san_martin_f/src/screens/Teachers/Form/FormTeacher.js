import React, { useState, useEffect, useMemo, useCallback } from 'react';
import { Form, Button } from 'react-bootstrap';

import clientAxios from '../../../config/axios';
import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';
import { NIVELES_ORDENADOS, etiquetaNivel } from '../../../constant/nivelesAcademicos';
import { validateDNI } from '../../../utils/dni';
import {
  validateCorreoElectronicoOpcional,
  validateTelefonoOpcional,
  sanitizeTelefonoInput
} from '../../../utils/contact';

const degId = (d) => (d && (d.id || d.ID)) || '';
const shiftId = (s) => {
  if (!s) return '';
  if (typeof s === 'string') {
    const str = s.trim();
    const objectIdMatch = str.match(/ObjectId\(["']?([a-fA-F0-9]{24})["']?\)/);
    if (objectIdMatch) return objectIdMatch[1];
    return str;
  }
  const raw = s.id || s.ID || s._id || s.shiftId || s.shiftid || '';
  if (raw && typeof raw === 'object') {
    if (raw.$oid) return String(raw.$oid).trim();
    if (raw.id) return String(raw.id).trim();
    if (raw.ID) return String(raw.ID).trim();
    if (raw._id) return shiftId(raw._id);
    if (raw.hex) return String(raw.hex).trim();
    if (raw.Hex) return String(raw.Hex).trim();
    const keys = Object.keys(raw);
    if (keys.length === 1) return shiftId(raw[keys[0]]);
  }
  const str = String(raw || '').trim();
  const objectIdMatch = str.match(/ObjectId\(["']?([a-fA-F0-9]{24})["']?\)/);
  if (objectIdMatch) return objectIdMatch[1];
  return str === '[object Object]' ? '' : str;
};

const normalizeCareerRow = (c) => ({
  degreeId: c.degreeId != null ? String(c.degreeId) : '',
  tituloHabilitanteId: c.tituloHabilitanteId != null ? String(c.tituloHabilitanteId) : '',
  modalidadId: c.modalidadId != null ? String(c.modalidadId) : '',
  shiftId: shiftId(c.shiftId || c.shiftid || c.ShiftID || c.shiftID)
});

/** Carreras cuyo nivel de la carrera está incluido en enseniaEn (lista de niveles del docente). */
const filterCareersByEnsenia = (careerRows, enseniaEn, degreesList) => {
  const niveles = new Set((enseniaEn || []).map((n) => String(n).toLowerCase().trim()));
  return (careerRows || []).filter((row) => {
    const did = String(row.degreeId || '');
    if (!did) return false;
    const d = (degreesList || []).find((x) => degId(x) === did);
    if (!d) return false;
    const n = String(d.nivel || '').toLowerCase().trim();
    return n && niveles.has(n);
  });
};

const catalogId = (rows, code) => {
  const u = String(code).toUpperCase();
  const r = (rows || []).find((x) => String(x.codigo).toUpperCase() === u);
  return r ? r.id || r.ID : '';
};

const FormTeacher = ({
  dataEntry,
  degrees,
  shifts,
  titulosHabilitantes,
  modalidades,
  saveData,
  changeActive
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
    active: true
  });
  const [careers, setCareers] = useState([]);
  const [fallbackShifts, setFallbackShifts] = useState([]);
  const [careerError, setCareerError] = useState('');
  const [activeFormConfirm, setActiveFormConfirm] = useState(null);
  const [activeFormSaving, setActiveFormSaving] = useState(false);

  const lockEntityKey = (dataEntry && dataEntry.id) || '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  const ids = useMemo(
    () => ({
      tituloSi: catalogId(titulosHabilitantes, 'SI'),
      tituloNo: catalogId(titulosHabilitantes, 'NO'),
      titular: catalogId(modalidades, 'TITULAR'),
      provisional: catalogId(modalidades, 'PROVISIONAL'),
      suplente: catalogId(modalidades, 'SUPLENTE')
    }),
    [titulosHabilitantes, modalidades]
  );
  const shiftOptions = useMemo(() => {
    const merged = [...(Array.isArray(shifts) ? shifts : []), ...(Array.isArray(fallbackShifts) ? fallbackShifts : [])];
    const byId = new Map();
    merged.forEach((s) => {
      const sid = shiftId(s);
      if (!sid) return;
      if (!byId.has(sid)) byId.set(sid, s);
    });
    return Array.from(byId.values());
  }, [shifts, fallbackShifts]);

  const degreeNivelInEnsenia = useCallback(
    (d) => {
      const n = String(d.nivel || '').toLowerCase().trim();
      return n && (data.enseniaEn || []).includes(n);
    },
    [data.enseniaEn]
  );

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : null;
    if (!entry) {
      setData({
        id: '',
        name: '',
        email: '',
        phone: '',
        dni: '',
        address: '',
        enseniaEn: [],
        active: true
      });
      setCareers([]);
      setActiveFormConfirm(null);
      return;
    }
    const enseniaEn = Array.isArray(entry.enseniaEn) ? [...entry.enseniaEn] : [];
    setData({
      id: entry.id || '',
      name: entry.name || '',
      email: entry.email || '',
      phone: sanitizeTelefonoInput(entry.phone || ''),
      dni: entry.dni || '',
      address: entry.address || '',
      enseniaEn,
      active: entry.active !== false
    });
    const raw = Array.isArray(entry.careers) ? entry.careers.map(normalizeCareerRow) : [];
    setCareers(filterCareersByEnsenia(raw, enseniaEn, degrees));
    setActiveFormConfirm(null);
  }, [dataEntry, degrees]);

  useEffect(() => {
    if (Array.isArray(shifts) && shifts.some((s) => shiftId(s))) {
      setFallbackShifts([]);
      return;
    }
    let mounted = true;
    const loadShifts = async () => {
      try {
        const access_token = document.cookie.replace('token=', '');
        const result = await clientAxios.get('/shift', {
          headers: { Authorization: `Bearer${access_token}` }
        });
        if (!mounted) return;
        const rows =
          result &&
          result.data &&
          Array.isArray(result.data.data)
            ? result.data.data
            : [];
        setFallbackShifts(rows);
      } catch (e) {
        if (mounted) setFallbackShifts([]);
      }
    };
    loadShifts();
    return () => {
      mounted = false;
    };
  }, [shifts]);

  const openActiveConfirm = (toActive) => {
    if (!data.id || !changeActive) return;
    const from = data.active !== false;
    if (toActive === from) return;
    setActiveFormConfirm({ fromActive: from, toActive });
  };

  const confirmActiveForm = async () => {
    if (!activeFormConfirm || !data.id || !changeActive) return;
    setActiveFormSaving(true);
    const res = await changeActive({ id: data.id, active: activeFormConfirm.toActive });
    setActiveFormSaving(false);
    const code = Number(res.code);
    if (code === 200) {
      setData((prev) => ({ ...prev, active: activeFormConfirm.toActive }));
      setActiveFormConfirm(null);
    } else if (code === 199) {
      setActiveFormConfirm(null);
    }
  };

  const tituloEsSi = (tituloHabilitanteId) => {
    const row = (titulosHabilitantes || []).find((x) => (x.id || x.ID) === tituloHabilitanteId);
    return row && String(row.codigo).toUpperCase() === 'SI';
  };

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
    setCareerError('');
    setData((prev) => {
      const set = new Set(prev.enseniaEn || []);
      if (set.has(nivel)) set.delete(nivel);
      else set.add(nivel);
      const enseniaEn = [...set];
      setCareers((careersPrev) => filterCareersByEnsenia(careersPrev, enseniaEn, degrees));
      return { ...prev, enseniaEn };
    });
  };

  const toggleDegree = (hexId) => {
    setCareerError('');
    if (!ids.tituloNo || !ids.provisional) return;
    setCareers((prev) => {
      const exists = prev.some((c) => c.degreeId === hexId);
      if (exists) {
        return prev.filter((c) => c.degreeId !== hexId);
      }
      return [
        ...prev,
        {
          degreeId: hexId,
          tituloHabilitanteId: ids.tituloNo,
          modalidadId: ids.provisional,
          shiftId: shiftOptions.length ? String(shiftId(shiftOptions[0])) : ''
        }
      ];
    });
  };

  const setTituloForCareer = (degreeId, tituloHabilitanteId) => {
    setCareers((prev) =>
      prev.map((c) => {
        if (c.degreeId !== degreeId) return c;
        let modalidadId = c.modalidadId;
        if (!tituloEsSi(tituloHabilitanteId) && c.modalidadId === ids.titular) {
          modalidadId = ids.provisional || modalidadId;
        }
        return { ...c, tituloHabilitanteId, modalidadId };
      })
    );
  };

  const setModalidadForCareer = (degreeId, modalidadId) => {
    setCareers((prev) =>
      prev.map((c) => (c.degreeId === degreeId ? { ...c, modalidadId } : c))
    );
  };

  const setShiftForCareer = (degreeId, selectedShiftId) => {
    setCareers((prev) =>
      prev.map((c) => (c.degreeId === degreeId ? { ...c, shiftId: selectedShiftId } : c))
    );
  };

  const sendData = async (e) => {
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
    if (!careers.length) {
      setCareerError('Debe seleccionar al menos una carrera y completar título y modalidad por carrera.');
      return;
    }
    if (careers.some((c) => !String(c.shiftId || '').trim())) {
      setCareerError('Debe indicar un turno (mañana, tarde o vespertino) por cada carrera.');
      return;
    }
    setCareerError('');
    const res = await saveData({
      ...data,
      active: data.active !== false,
      dni: String(data.dni).trim(),
      email: String(data.email).trim(),
      phone: data.phone,
      careers
    });
    if (isSaveSuccess(res)) {
      armLockAfterSave();
    }
  };

  const selectableDegrees = (degrees || []).filter((d) => {
    const hid = degId(d);
    if (!hid) return false;
    if (d.active === true) return true;
    if (careers.some((c) => c.degreeId === hid)) return true;
    return false;
  });

  const careerForDegree = (hexId) => careers.find((c) => c.degreeId === hexId);

  return (
    <Form onSubmit={sendData}>
      <Form.Control type="hidden" name="id" value={data.id} readOnly />

      <FormEditLockBanner
        entityKey={lockEntityKey}
        estadoActivo={lockEntityKey ? data.active !== false : undefined}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
        onCancelUnlock={() => setUnlocked(false)}
        unlockVariant="warning"
      />

      <Form.Group className="mb-3" controlId="teacherId">
        <Form.Label>ID docente</Form.Label>
        <Form.Control type="text" readOnly plaintext value={data.id ? data.id : 'Se asigna al guardar'} />
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherName">
        <Form.Label>Nombre</Form.Label>
        <Form.Control
          type="text"
          name="name"
          value={data.name}
          onChange={handleInputChange}
          required
          autoFocus
          readOnly={readOnly}
        />
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
          readOnly={readOnly}
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
          readOnly={readOnly}
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
          readOnly={readOnly}
        />
        <Form.Control.Feedback type="invalid">{dniError}</Form.Control.Feedback>
        <Form.Text className="text-muted">7 u 8 dígitos numéricos, sin puntos.</Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="teacherAddress">
        <Form.Label>Dirección</Form.Label>
        <Form.Control
          type="text"
          name="address"
          value={data.address}
          onChange={handleInputChange}
          readOnly={readOnly}
        />
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
              disabled={readOnly}
            />
          ))}
        </div>
        <Form.Text className="text-muted">Puede marcar más de un nivel. Las carreras deben coincidir con estos niveles.</Form.Text>
      </Form.Group>

      {data.id && unlocked ? (
        <Form.Group className="mb-3" controlId="teacherEstado">
          <Form.Label>Estado</Form.Label>
          <div>
            <Form.Check
              inline
              type="radio"
              id="teacher-active-si"
              name="teacherEstadoActivo"
              label="Activo"
              checked={data.active !== false}
              onChange={() => openActiveConfirm(true)}
            />
            <Form.Check
              inline
              type="radio"
              id="teacher-active-no"
              name="teacherEstadoActivo"
              label="Inactivo"
              checked={data.active === false}
              onChange={() => openActiveConfirm(false)}
            />
          </div>
          <Form.Text className="text-muted">El cambio requiere confirmación.</Form.Text>
        </Form.Group>
      ) : null}

      <Form.Group className="mb-3" controlId="teacherDegrees">
        <Form.Label>Carreras y condición laboral</Form.Label>
        {careerError ? (
          <div className="text-danger small mb-2">{careerError}</div>
        ) : null}
        <div className="border rounded p-2" style={{ maxHeight: 'none', overflowY: 'visible' }}>
          {selectableDegrees.length === 0 ? (
            <span className="text-muted">No hay carreras activas. Cree carreras primero.</span>
          ) : (
            selectableDegrees.map((d) => {
              const id = degId(d);
              if (!id) return null;
              const allowed = degreeNivelInEnsenia(d);
              const checked = !!careerForDegree(id);
              const row = careerForDegree(id);
              const si = row && tituloEsSi(row.tituloHabilitanteId);
              return (
                <div key={id} className="mb-3 pb-2 border-bottom">
                  <Form.Check
                    type="checkbox"
                    id={`deg-${id}`}
                    label={d.name + (d.nivel ? ` (${etiquetaNivel(d.nivel)})` : '')}
                    checked={checked}
                    disabled={readOnly || !allowed}
                    title={!allowed ? 'Marque primero el nivel correspondiente en Enseña en' : ''}
                    onChange={() => {
                      if (allowed && !readOnly) toggleDegree(id);
                    }}
                  />
                  {checked && row && (
                    <div className="ms-4 mt-2 small">
                      <div className="mb-2">
                        <span className="fw-semibold d-block mb-1">Título habilitante (esta carrera)</span>
                        {(titulosHabilitantes || []).map((t) => {
                          const tid = t.id || t.ID;
                          if (!tid) return null;
                          return (
                            <Form.Check
                              key={tid}
                              inline
                              type="radio"
                              name={`titulo-${id}`}
                              id={`titulo-${id}-${tid}`}
                              label={String(t.codigo).toUpperCase() === 'SI' ? 'Sí' : 'No'}
                              checked={row.tituloHabilitanteId === tid}
                              disabled={readOnly}
                              onChange={() => setTituloForCareer(id, tid)}
                            />
                          );
                        })}
                      </div>
                      <div>
                        <span className="fw-semibold d-block mb-1">Modalidad (esta carrera)</span>
                        {(modalidades || []).map((m) => {
                          const mid = m.id || m.ID;
                          if (!mid) return null;
                          const cod = String(m.codigo || '').toUpperCase();
                          const disTitular = cod === 'TITULAR' && !si;
                          return (
                            <Form.Check
                              key={mid}
                              inline
                              type="radio"
                              name={`mod-${id}`}
                              id={`mod-${id}-${mid}`}
                              label={etiquetaModalidad(m.codigo)}
                              checked={row.modalidadId === mid}
                              disabled={readOnly || disTitular}
                              title={
                                disTitular ? 'Titular solo con título habilitante Sí' : ''
                              }
                              onChange={() => {
                                if (!disTitular && !readOnly) setModalidadForCareer(id, mid);
                              }}
                            />
                          );
                        })}
                      </div>
                      <div className="mt-2">
                        <span className="fw-semibold d-block mb-1">Turno (esta carrera)</span>
                        {shiftOptions.length === 0 ? (
                          <span className="text-muted small">No hay turnos cargados en Administración → Turnos.</span>
                        ) : (
                          shiftOptions.map((s) => {
                            const sid = String(shiftId(s));
                            const selectedSid = String(shiftId({ shiftId: row.shiftId }));
                            return (
                              <Form.Check
                                key={sid}
                                inline
                                type="radio"
                                name={`shift-${id}`}
                                id={`shift-${id}-${sid}`}
                                label={s.type || sid}
                                checked={selectedSid === sid}
                                disabled={readOnly}
                                onChange={() => setShiftForCareer(id, sid)}
                              />
                            );
                          })
                        )}
                      </div>
                    </div>
                  )}
                </div>
              );
            })
          )}
        </div>
        <Form.Text className="text-muted">
          Por cada carrera indique título habilitante, modalidad y turno. Sin título habilitante puede ser provisional o
          suplente (una u otra, no titular).
        </Form.Text>
      </Form.Group>

      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {data.id ? 'Actualizar' : 'Guardar'}
      </Button>

      <ConfirmChangeEstadoModal
        show={!!activeFormConfirm}
        onHide={() => !activeFormSaving && setActiveFormConfirm(null)}
        kind="docente"
        itemName={data.name || ''}
        fromActive={activeFormConfirm ? activeFormConfirm.fromActive : true}
        toActive={activeFormConfirm ? activeFormConfirm.toActive : true}
        onConfirm={confirmActiveForm}
        confirming={activeFormSaving}
      />
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
