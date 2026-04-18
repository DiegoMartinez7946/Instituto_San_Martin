import React, { useState, useEffect, useMemo } from 'react';
import { Form, Button, InputGroup, Row, Col } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';
import {
  NIVELES_ORDENADOS,
  jerarquiaDeNivel,
  etiquetaNivel,
  mapaOrdenDesdeLevels,
  ORDEN_POR_NIVEL_DEFAULT
} from '../../../constant/nivelesAcademicos';
import { useGlobal } from '../../../context/Global/GlobalProvider';
import { validateDNI } from '../../../utils/dni';
import {
  validateCorreoElectronicoOpcional,
  validateTelefonoOpcional,
  sanitizeTelefonoInput
} from '../../../utils/contact';

const degId = (d) => (d && (d.id || d.ID)) || '';

const shiftSid = (s) => {
  if (!s) return '';
  const raw = s.id != null ? s.id : s.ID;
  if (raw && typeof raw === 'object') {
    if (raw.$oid) return String(raw.$oid).trim();
    if (raw.hex) return String(raw.hex).trim();
  }
  return String(raw || '').trim();
};

const shiftTypeToTurno = (type) => {
  const u = String(type || '').toUpperCase();
  if (u === 'MAÑANA' || u === 'MANANA') return 'manana';
  if (u === 'TARDE') return 'tarde';
  if (u === 'VESPERTINO') return 'noche';
  return '';
};

const shiftsForDegree = (degree, shiftsList) => {
  const turnos = Array.isArray(degree && degree.turnos)
    ? degree.turnos.map((t) => String(t || '').toLowerCase().trim()).filter(Boolean)
    : [];
  return (shiftsList || []).filter((sh) => {
    if (!turnos.length) return true;
    const code = shiftTypeToTurno(sh.type || sh.Type);
    return code && turnos.includes(code);
  });
};

/** degreeIds en JSON; si el backend solo envía degreeShifts, deducimos la lista de carreras. */
const degreeIdsFromEntry = (entry) => {
  if (!entry) return [];
  const ids = Array.isArray(entry.degreeIds)
    ? entry.degreeIds.map((x) => String(x).trim()).filter(Boolean)
    : [];
  if (ids.length) return ids;
  const rows = Array.isArray(entry.degreeShifts) ? entry.degreeShifts : [];
  return rows
    .map((r) => String(r.degreeId != null ? r.degreeId : r.degreeid || '').trim())
    .filter(Boolean);
};

const buildDegreeShiftsMap = (entry, degreeIdList, degList, shiftsList) => {
  const map = {};
  const rows = entry && Array.isArray(entry.degreeShifts) ? entry.degreeShifts : [];
  rows.forEach((r) => {
    const did = String(r.degreeId != null ? r.degreeId : r.degreeid || '').trim();
    const sid = String(r.shiftId != null ? r.shiftId : r.shiftid || '').trim();
    if (did && sid) map[did] = sid;
  });
  (degreeIdList || []).forEach((did) => {
    if (map[did]) return;
    const deg = (degList || []).find((x) => degId(x) === did);
    const opts = shiftsForDegree(deg, shiftsList);
    const id0 = opts[0] ? shiftSid(opts[0]) : '';
    if (id0) map[did] = id0;
  });
  return map;
};

const FormStudent = ({ dataEntry, degrees, shifts = [], saveData, changeActive }) => {

  const [dniError, setDniError] = useState('');
  const [emailError, setEmailError] = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [formError, setFormError] = useState('');
  const [data, setData] = useState({
    id: '',
    name: '',
    email: '',
    phone: '',
    dni: '',
    address: '',
    nivelAprobado: '',
    modalidad: '',
    condicion: '',
    degreeIds: [],
    degreeShiftsByDegree: {},
    active: true
  });
  const [activeFormConfirm, setActiveFormConfirm] = useState(null);
  const [activeFormSaving, setActiveFormSaving] = useState(false);
  const [newPortalPassword, setNewPortalPassword] = useState('');
  const [portalPwdError, setPortalPwdError] = useState('');
  const [portalPwdRevealPress, setPortalPwdRevealPress] = useState(false);

  const lockEntityKey = (dataEntry && (dataEntry.id || dataEntry.idAlumno)) || '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  const [globalState] = useGlobal();
  const ordenMap = useMemo(
    () => mapaOrdenDesdeLevels(globalState.levels) || ORDEN_POR_NIVEL_DEFAULT,
    [globalState.levels]
  );

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : null;
    const isActive = entry ? entry.active !== false : true;
    const rawModalidad =
      entry && (entry.modalidad != null && entry.modalidad !== ''
        ? entry.modalidad
        : entry.Modalidad);
    const rawCondicion =
      entry && (entry.condicion != null && entry.condicion !== ''
        ? entry.condicion
        : entry.Condicion);
    const degreeIdsList = degreeIdsFromEntry(entry);
    setData({
      id: (entry && (entry.id || entry.idAlumno)) || '',
      name: (entry && entry.name) || '',
      email: (entry && entry.email) || '',
      phone: sanitizeTelefonoInput((entry && entry.phone) || ''),
      dni: (entry && entry.dni) || '',
      address: (entry && entry.address) || '',
      nivelAprobado: (entry && entry.nivelAprobado)
        ? String(entry.nivelAprobado).toLowerCase()
        : '',
      modalidad:
        isActive && rawModalidad ? String(rawModalidad).toUpperCase() : '',
      condicion:
        isActive && rawCondicion ? String(rawCondicion).toUpperCase() : '',
      degreeIds: degreeIdsList,
      degreeShiftsByDegree: buildDegreeShiftsMap(
        entry,
        degreeIdsList,
        degrees,
        shifts
      ),
      active: isActive
    });
    setActiveFormConfirm(null);
    setNewPortalPassword('');
    setPortalPwdError('');
    setPortalPwdRevealPress(false);
  }, [dataEntry]);

  useEffect(() => {
    if (!portalPwdRevealPress) return undefined;
    const release = () => setPortalPwdRevealPress(false);
    window.addEventListener('mouseup', release);
    window.addEventListener('touchend', release);
    window.addEventListener('touchcancel', release);
    return () => {
      window.removeEventListener('mouseup', release);
      window.removeEventListener('touchend', release);
      window.removeEventListener('touchcancel', release);
    };
  }, [portalPwdRevealPress]);

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
      setData((prev) => ({
        ...prev,
        active: activeFormConfirm.toActive,
        modalidad: activeFormConfirm.toActive ? prev.modalidad : '',
        condicion: activeFormConfirm.toActive ? prev.condicion : ''
      }));
      setActiveFormConfirm(null);
    } else if (code === 199) {
      setActiveFormConfirm(null);
    }
  };

  const handleInputChange = (event) => {
    const t = event.currentTarget || event.target;
    const name = t.name;
    const value = t.value;
    if (!name) {
      return;
    }
    if (name === 'dni') {
      setDniError('');
    }
    if (name === 'email') {
      setEmailError('');
    }
    setData((prev) => {
      const next = { ...prev, [name]: value };
      if (name === 'nivelAprobado' && (degrees || []).length) {
        const ha = jerarquiaDeNivel(value, ordenMap);
        const kept = (prev.degreeIds || []).filter((hexId) => {
          const d = (degrees || []).find((x) => degId(x) === hexId);
          if (!d) return false;
          const hn = jerarquiaDeNivel((d.nivel || '').toLowerCase().trim(), ordenMap);
          if (ha == null || hn == null) return false;
          return ha >= hn - 1;
        });
        next.degreeIds = kept;
        const nextShifts = { ...prev.degreeShiftsByDegree };
        Object.keys(nextShifts).forEach((k) => {
          if (!kept.includes(k)) delete nextShifts[k];
        });
        next.degreeShiftsByDegree = nextShifts;
      }
      return next;
    });
  };

  const handlePhoneChange = (event) => {
    setPhoneError('');
    setFormError('');
    const digits = sanitizeTelefonoInput(event.target.value);
    setData((prev) => ({ ...prev, phone: digits }));
  };

  const toggleDegree = (hexId) => {
    setData((prev) => {
      const set = new Set(prev.degreeIds || []);
      const nextShifts = { ...prev.degreeShiftsByDegree };
      if (set.has(hexId)) {
        set.delete(hexId);
        delete nextShifts[hexId];
      } else {
        set.add(hexId);
        const deg = (degrees || []).find((x) => degId(x) === hexId);
        const opts = shiftsForDegree(deg, shifts);
        const first = opts[0] ? shiftSid(opts[0]) : '';
        if (first) nextShifts[hexId] = first;
      }
      return { ...prev, degreeIds: [...set], degreeShiftsByDegree: nextShifts };
    });
  };

  const setShiftForDegree = (degreeHexId, shiftHexId) => {
    setData((prev) => ({
      ...prev,
      degreeShiftsByDegree: { ...prev.degreeShiftsByDegree, [degreeHexId]: shiftHexId }
    }));
  };

  const haAlumno = jerarquiaDeNivel(data.nivelAprobado, ordenMap);

  const degreeCheckboxDisabled = (d) => {
    if (!data.nivelAprobado) return true;
    const nivelCarrera = (d.nivel || '').toLowerCase().trim();
    if (!nivelCarrera) return true;
    const hn = jerarquiaDeNivel(nivelCarrera, ordenMap);
    if (hn == null || haAlumno == null) return true;
    return haAlumno < hn - 1;
  };

  const degreeCheckboxTitle = (d) => {
    if (!data.nivelAprobado) return 'Seleccione primero el nivel aprobado del alumno';
    if (!(d.nivel || '').trim()) return 'La carrera no tiene nivel configurado; edítela en Carreras';
    if (degreeCheckboxDisabled(d)) return 'El nivel aprobado no alcanza para inscribir en esta carrera';
    return '';
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
    if (data.active !== false) {
      if (!String(data.modalidad || '').trim()) {
        setFormError('Seleccione si el alumno es presencial o remoto');
        return;
      }
      if (!String(data.condicion || '').trim()) {
        setFormError('Seleccione si el alumno es regular o libre');
        return;
      }
    }
    const np = String(newPortalPassword || '').trim();
    if (np.length > 0 && np.length < 6) {
      setPortalPwdError('La contraseña de acceso al portal debe tener al menos 6 caracteres (o déjela vacía).');
      return;
    }
    setPortalPwdError('');
    const ids = data.degreeIds || [];
    const degreeShifts = [];
    for (let i = 0; i < ids.length; i += 1) {
      const did = ids[i];
      const sid = (data.degreeShiftsByDegree && data.degreeShiftsByDegree[did]) || '';
      if (!String(sid).trim()) {
        setFormError('Seleccione el turno de cursada para cada carrera marcada.');
        return;
      }
      degreeShifts.push({ degreeId: did, shiftId: String(sid).trim() });
    }
    setFormError('');
    const res = await saveData({
      ...data,
      dni: String(data.dni).trim(),
      email: String(data.email).trim(),
      phone: data.phone,
      modalidad: data.active === false ? '' : String(data.modalidad || '').toUpperCase(),
      condicion: data.active === false ? '' : String(data.condicion || '').toUpperCase(),
      degreeShifts,
      newPassword: np.length >= 6 ? np : ''
    });
    if (isSaveSuccess(res)) {
      setNewPortalPassword('');
      armLockAfterSave();
      setFormError('');
    }
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

      <FormEditLockBanner
        entityKey={lockEntityKey}
        estadoActivo={lockEntityKey ? data.active !== false : undefined}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
        onCancelUnlock={() => setUnlocked(false)}
        unlockVariant="warning"
      />

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
          readOnly={readOnly}
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
          readOnly={readOnly}
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
          readOnly={readOnly}
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
          readOnly={readOnly}
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
          readOnly={readOnly}
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
              disabled={readOnly}
            />
          ))}
        </div>
        <Form.Text className="text-muted">
          Indica el máximo nivel cursado/aprobado; los niveles inferiores se consideran aprobados.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentModalidadCondicion">
        <Form.Label>Rol en sistema</Form.Label>
        <div className="border rounded p-2">
          <Form.Label className="small text-muted mb-1">Modalidad del alumno</Form.Label>
          <Form.Check
            type="radio"
            id="student-modalidad-presencial"
            name="modalidad"
            value="PRESENCIAL"
            label="Presencial"
            checked={data.modalidad === 'PRESENCIAL'}
            onChange={() =>
              setData((prev) => ({ ...prev, modalidad: 'PRESENCIAL' }))
            }
            disabled={readOnly || data.active === false}
          />
          <Form.Check
            type="radio"
            id="student-modalidad-remoto"
            name="modalidad"
            value="REMOTO"
            label="Remoto"
            checked={data.modalidad === 'REMOTO'}
            onChange={() => setData((prev) => ({ ...prev, modalidad: 'REMOTO' }))}
            disabled={readOnly || data.active === false}
          />

          <hr className="my-2" />
          <Form.Label className="small text-muted mb-1">Condición académica</Form.Label>
          <Form.Check
            type="radio"
            id="student-condicion-regular"
            name="condicion"
            value="REGULAR"
            label="Regular"
            checked={data.condicion === 'REGULAR'}
            onChange={() =>
              setData((prev) => ({ ...prev, condicion: 'REGULAR' }))
            }
            disabled={readOnly || data.active === false}
          />
          <Form.Check
            type="radio"
            id="student-condicion-libre"
            name="condicion"
            value="LIBRE"
            label="Libre"
            checked={data.condicion === 'LIBRE'}
            onChange={() => setData((prev) => ({ ...prev, condicion: 'LIBRE' }))}
            disabled={readOnly || data.active === false}
          />
        </div>
        <Form.Text className="text-muted">
          Si el alumno está inactivo, estos campos quedan en blanco hasta reactivarlo.
        </Form.Text>
      </Form.Group>

      {data.id && unlocked ? (
        <Form.Group className="mb-3" controlId="studentEstado">
          <Form.Label>Estado</Form.Label>
          <div>
            <Form.Check
              inline
              type="radio"
              id="student-active-si"
              name="studentEstadoActivo"
              label="Activo"
              checked={data.active !== false}
              onChange={() => openActiveConfirm(true)}
            />
            <Form.Check
              inline
              type="radio"
              id="student-active-no"
              name="studentEstadoActivo"
              label="Inactivo"
              checked={data.active === false}
              onChange={() => openActiveConfirm(false)}
            />
          </div>
          <Form.Text className="text-muted">El cambio requiere confirmación.</Form.Text>
        </Form.Group>
      ) : null}

      <Form.Group className="mb-3" controlId="studentDegrees">
        <Form.Label>Carreras y turno de cursada</Form.Label>
        <Row className="g-2">
          <Col md={6}>
            <div className="border rounded p-2" style={{ maxHeight: '240px', overflowY: 'auto' }}>
              <div className="small text-muted mb-1">Carreras</div>
              {selectableDegrees.length === 0 ? (
                <span className="text-muted">No hay carreras activas. Cree carreras primero.</span>
              ) : (
                selectableDegrees.map((d) => {
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
                      disabled={readOnly || dis}
                      onChange={() => {
                        if (!dis && !readOnly) toggleDegree(id);
                      }}
                    />
                  );
                })
              )}
            </div>
          </Col>
          <Col md={6}>
            <div className="border rounded p-2" style={{ maxHeight: '240px', overflowY: 'auto' }}>
              <div className="small text-muted mb-1">Turno por carrera</div>
              {(data.degreeIds || []).length === 0 ? (
                <span className="text-muted small">Marque al menos una carrera a la izquierda.</span>
              ) : (
                (data.degreeIds || []).map((did) => {
                  const deg = (degrees || []).find((x) => degId(x) === did);
                  const label = deg
                    ? deg.name + (deg.nivel ? ` (${etiquetaNivel(deg.nivel)})` : '')
                    : did;
                  const opts = shiftsForDegree(deg, shifts);
                  const cur = (data.degreeShiftsByDegree && data.degreeShiftsByDegree[did]) || '';
                  return (
                    <div key={did} className="mb-3 pb-2 border-bottom">
                      <div className="small fw-semibold mb-1">{label}</div>
                      {opts.length === 0 ? (
                        <span className="text-danger small">
                          No hay turnos en Administración → Turnos compatibles con esta carrera (configure turnos en
                          la carrera y en Turnos).
                        </span>
                      ) : (
                        opts.map((sh) => {
                          const sid = shiftSid(sh);
                          if (!sid) return null;
                          return (
                            <Form.Check
                              key={sid}
                              inline
                              type="radio"
                              name={`student-shift-${did}`}
                              id={`student-shift-${did}-${sid}`}
                              label={sh.type || sh.Type || sid}
                              checked={cur === sid}
                              disabled={readOnly}
                              onChange={() => setShiftForDegree(did, sid)}
                            />
                          );
                        })
                      )}
                    </div>
                  );
                })
              )}
            </div>
          </Col>
        </Row>
        <Form.Text className="text-muted">
          Por cada carrera elegida indique el turno de cursada (mañana, tarde o vespertino según la colección Turnos).
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="studentPortalPassword">
        <Form.Label>Contraseña de acceso al portal (alumno)</Form.Label>
        <InputGroup>
          <Form.Control
            type={portalPwdRevealPress ? 'text' : 'password'}
            autoComplete="new-password"
            value={newPortalPassword}
            onChange={(e) => {
              setNewPortalPassword(e.target.value);
              setPortalPwdError('');
            }}
            placeholder={data.id ? 'Dejar vacío para no cambiar' : 'Mínimo 6 caracteres para permitir ingreso al sistema'}
            readOnly={readOnly}
          />
          <Button
            variant="outline-secondary"
            type="button"
            title="Mantener presionado para ver la contraseña"
            aria-label="Mostrar contraseña mientras se mantiene presionado"
            tabIndex={-1}
            disabled={readOnly}
            onMouseDown={(e) => {
              e.preventDefault();
              setPortalPwdRevealPress(true);
            }}
            onTouchStart={(e) => {
              e.preventDefault();
              setPortalPwdRevealPress(true);
            }}
          >
            <FontAwesomeIcon icon={portalPwdRevealPress ? faEyeSlash : faEye} />
          </Button>
        </InputGroup>
        {portalPwdError ? <div className="text-danger small mt-1">{portalPwdError}</div> : null}
        <Form.Text className="text-muted">
          El alumno podrá iniciar sesión con su correo y esta contraseña. Solo el administrativo puede asignarla o
          cambiarla desde aquí; el alumno puede cambiarla luego en Cambiar Password.
        </Form.Text>
      </Form.Group>

      {formError ? <div className="text-danger small mb-2">{formError}</div> : null}

      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {data.id ? 'Actualizar' : 'Guardar'}
      </Button>

      <ConfirmChangeEstadoModal
        show={!!activeFormConfirm}
        onHide={() => !activeFormSaving && setActiveFormConfirm(null)}
        kind="alumno"
        itemName={data.name || ''}
        fromActive={activeFormConfirm ? activeFormConfirm.fromActive : true}
        toActive={activeFormConfirm ? activeFormConfirm.toActive : true}
        onConfirm={confirmActiveForm}
        confirming={activeFormSaving}
      />
    </Form>
  );
};

export default FormStudent;
