import React, { useState, useEffect } from 'react';
import { Form, Button, InputGroup } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';
import { validateDNI } from '../../../utils/dni';
import {
  validateCorreoElectronicoRequerido,
  validateTelefonoOpcional,
  sanitizeTelefonoInput
} from '../../../utils/contact';

const shiftRowHex = (s) => {
  if (!s) return '';
  const v = s.ID !== undefined && s.ID !== null && s.ID !== '' ? s.ID : s.id;
  return String(v !== undefined && v !== null ? v : '').trim();
};

const shiftIdsContain = (shiftIds, shift) => {
  const h = shiftRowHex(shift).toLowerCase();
  if (!h) return false;
  return (shiftIds || []).some((x) => String(x).trim().toLowerCase() === h);
};

const degRowId = (d) => (d && (d.id || d.ID)) || '';

const FormUser = ({ dataEntry, roles, shifts, degrees, saveData, changeActive }) => {
  const [emailError, setEmailError] = useState('');
  const [dniError, setDniError] = useState('');
  const [phoneError, setPhoneError] = useState('');
  const [formError, setFormError] = useState('');
  const [data, setData] = useState({
    id: '',
    name: '',
    email: '',
    phone: '',
    dni: '',
    address: '',
    userType: '',
    shiftIds: [],
    degreeId: '',
    modalidad: '',
    condicion: '',
    studentRecordId: '',
    password: '',
    active: true
  });
  const [activeFormConfirm, setActiveFormConfirm] = useState(null);
  const [activeFormSaving, setActiveFormSaving] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const lockEntityKey = (dataEntry && dataEntry.id) || '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  useEffect(() => {
    setShowPassword(false);
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : null;
    if (!entry || !entry.id) {
      setData({
        id: '',
        name: '',
        email: '',
        phone: '',
        dni: '',
        address: '',
        userType: '',
        shiftIds: [],
        degreeId: '',
        modalidad: '',
        condicion: '',
        studentRecordId: '',
        password: '',
        active: true
      });
      setActiveFormConfirm(null);
      return;
    }
    const firstDegree =
      Array.isArray(entry.degreeIds) && entry.degreeIds.length
        ? String(entry.degreeIds[0]).trim()
        : '';
    const isActive = entry.active !== false;
    const mRaw = entry.modalidad != null && entry.modalidad !== '' ? entry.modalidad : entry.Modalidad;
    const cRaw = entry.condicion != null && entry.condicion !== '' ? entry.condicion : entry.Condicion;
    setData({
      id: entry.id,
      name: entry.name || '',
      email: entry.email || '',
      phone: sanitizeTelefonoInput(entry.phone || ''),
      dni: entry.dni || '',
      address: entry.address || '',
      userType: (entry.role || '').toString().toUpperCase(),
      shiftIds: Array.isArray(entry.shiftIds)
        ? entry.shiftIds.map((x) => String(x).trim()).filter(Boolean)
        : entry.shiftId
          ? [String(entry.shiftId).trim()].filter(Boolean)
          : [],
      degreeId: firstDegree,
      modalidad: isActive ? String(mRaw != null ? mRaw : '').trim().toUpperCase() : '',
      condicion: isActive ? String(cRaw != null ? cRaw : '').trim().toUpperCase() : '',
      studentRecordId: entry.studentRecordId ? String(entry.studentRecordId).trim() : '',
      password: '',
      active: isActive
    });
    setActiveFormConfirm(null);
  }, [dataEntry]);

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
    const target = event.currentTarget || event.target;
    const name = target.name;
    const value = target.value;
    if (name === 'email') setEmailError('');
    if (name === 'dni') setDniError('');
    if (name === 'userType') {
      const upper = String(value).toUpperCase();
      setData((prev) => ({
        ...prev,
        userType: value,
        degreeId: upper === 'ALUMNO' || upper === 'DOCENTE' ? prev.degreeId : '',
        modalidad: upper === 'ALUMNO' ? prev.modalidad : '',
        condicion: upper === 'ALUMNO' ? prev.condicion : ''
      }));
      return;
    }
    setData((prev) => ({ ...prev, [name]: value }));
  };

  const handleDegreeChange = (event) => {
    setData((prev) => ({ ...prev, degreeId: event.target.value }));
  };

  const toggleShiftId = (shift) => {
    const canon = shiftRowHex(shift);
    if (!canon) return;
    const lower = canon.toLowerCase();
    setData((prev) => {
      const ids = [...(prev.shiftIds || []).map((x) => String(x).trim()).filter(Boolean)];
      const idx = ids.findIndex((x) => x.toLowerCase() === lower);
      if (idx >= 0) ids.splice(idx, 1);
      else ids.push(canon);
      return { ...prev, shiftIds: ids };
    });
  };

  const handlePhoneChange = (event) => {
    setPhoneError('');
    setData((prev) => ({ ...prev, phone: sanitizeTelefonoInput(event.target.value) }));
  };

  const sendData = async (e) => {
    e.preventDefault();
    setEmailError('');
    setDniError('');
    setPhoneError('');
    setFormError('');

    const emailCheck = validateCorreoElectronicoRequerido(data.email);
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

    if (!data.id && (!data.password || data.password.length < 6)) {
      setFormError('La contraseña debe tener al menos 6 caracteres');
      return;
    }

    if (!data.userType) {
      setFormError('Seleccione un rol');
      return;
    }

    const roleUpper = String(data.userType).toUpperCase();
    const degreeIds =
      roleUpper === 'ALUMNO' || roleUpper === 'DOCENTE'
        ? data.degreeId
          ? [String(data.degreeId).trim()]
          : []
        : [];

    const hasStudentFicha = !!(
      (dataEntry && dataEntry.studentLinked) ||
      (data.studentRecordId && String(data.studentRecordId).trim())
    );
    if (
      roleUpper === 'ALUMNO' &&
      data.active !== false &&
      hasStudentFicha &&
      (!String(data.modalidad || '').trim() || !String(data.condicion || '').trim())
    ) {
      setFormError(
        'Indique modalidad (presencial/remoto) y condición (regular/libre) para la ficha del alumno.'
      );
      return;
    }

    const payload = {
      id: data.id,
      name: String(data.name).trim(),
      email: String(data.email).trim(),
      phone: data.phone,
      dni: String(data.dni).trim(),
      address: String(data.address || '').trim(),
      userType: String(data.userType).toUpperCase(),
      shiftIds: (data.shiftIds || []).map((x) => String(x).trim()).filter(Boolean),
      degreeIds,
      modalidad: data.active === false ? '' : String(data.modalidad || '').toUpperCase(),
      condicion: data.active === false ? '' : String(data.condicion || '').toUpperCase(),
      studentRecordId: String(data.studentRecordId || '').trim(),
      password: data.password
    };

    const res = await saveData(payload);
    if (isSaveSuccess(res)) {
      armLockAfterSave();
      if (!data.id) {
        setData({
          id: '',
          name: '',
          email: '',
          phone: '',
          dni: '',
          address: '',
          userType: '',
          shiftIds: [],
          degreeId: '',
          modalidad: '',
          condicion: '',
          studentRecordId: '',
          password: '',
          active: true
        });
      } else {
        setData((prev) => ({ ...prev, password: '' }));
      }
    }
  };

  const roleOptions = (roles || []).filter((r) => r && r.type);
  const shiftOptions = (shifts || []).filter((s) => s && shiftRowHex(s));
  const degreeOptions = (degrees || []).filter((d) => d && degRowId(d) && d.active !== false);
  const roleUpper = String(data.userType || '').toUpperCase();
  const showInscriptoCarrera = roleUpper === 'ALUMNO' || roleUpper === 'DOCENTE';

  return (
    <Form onSubmit={sendData}>
      <FormEditLockBanner
        entityKey={lockEntityKey}
        estadoActivo={lockEntityKey ? data.active !== false : undefined}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
      />

      {data.id ? (
        <Form.Group className="mb-3" controlId="userId">
          <Form.Label>ID usuario</Form.Label>
          <Form.Control type="text" readOnly plaintext value={data.id} />
        </Form.Group>
      ) : null}

      <Form.Group className="mb-3" controlId="userName">
        <Form.Label>Nombre completo</Form.Label>
        <Form.Control
          type="text"
          name="name"
          value={data.name}
          onChange={handleInputChange}
          required
          autoFocus={!data.id}
          readOnly={readOnly}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="userEmail">
        <Form.Label>Correo electrónico</Form.Label>
        <Form.Control
          type="email"
          name="email"
          value={data.email}
          onChange={handleInputChange}
          isInvalid={!!emailError}
          readOnly={readOnly}
        />
        <Form.Control.Feedback type="invalid">{emailError}</Form.Control.Feedback>
      </Form.Group>

      <Form.Group className="mb-3" controlId="userPhone">
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

      <Form.Group className="mb-3" controlId="userDni">
        <Form.Label>DNI</Form.Label>
        <Form.Control
          type="text"
          name="dni"
          inputMode="numeric"
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

      <Form.Group className="mb-3" controlId="userAddress">
        <Form.Label>Dirección</Form.Label>
        <Form.Control
          type="text"
          name="address"
          value={data.address}
          onChange={handleInputChange}
          readOnly={readOnly}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="userRolTurno">
        <Form.Label className="d-block">Rol y turno de cursado</Form.Label>
        <div className="border rounded p-3 bg-white">
          <Form.Group className="mb-3 mb-md-0" controlId="userRole">
            <Form.Label className="small text-muted mb-1">Rol en el sistema</Form.Label>
            <Form.Select
              name="userType"
              value={data.userType}
              onChange={handleInputChange}
              required
              disabled={readOnly}
            >
              <option value="">Seleccione un rol</option>
              {roleOptions.length === 0 ? (
                <option value="">Sin roles cargados</option>
              ) : (
                roleOptions.map((r) => (
                  <option key={r.ID || r.id} value={String(r.type).toUpperCase()}>
                    {r.type}
                  </option>
                ))
              )}
            </Form.Select>
          </Form.Group>
          <Form.Group className="mt-3" controlId="userShifts">
            <Form.Label className="small text-muted mb-1">
              Turnos de cursado (puede elegir varios; mañana / tarde / vespertino)
            </Form.Label>
            <div className="d-flex flex-column gap-2">
              {shiftOptions.length === 0 ? (
                <span className="text-muted small">Sin turnos cargados</span>
              ) : (
                shiftOptions.map((s) => {
                  const sid = shiftRowHex(s);
                  if (!sid) return null;
                  const checked = shiftIdsContain(data.shiftIds, s);
                  return (
                    <Form.Check
                      key={sid}
                      type="checkbox"
                      id={`user-shift-${sid}`}
                      label={s.type}
                      checked={checked}
                      onChange={() => toggleShiftId(s)}
                      disabled={readOnly}
                    />
                  );
                })
              )}
            </div>
            <Form.Text className="text-muted">
              Opcional. Los turnos se administran en Administración → Turnos de cursado.
            </Form.Text>
          </Form.Group>

          {showInscriptoCarrera ? (
            <Form.Group className="mt-3" controlId="userDegrees">
              <Form.Label className="small text-muted mb-1">Inscripto en</Form.Label>
              <Form.Text className="text-muted d-block small mb-2">
                {roleUpper === 'DOCENTE'
                  ? 'Carrera en la que se desempeña como docente.'
                  : 'Carrera que cursa el alumno.'}
              </Form.Text>
              <div className="d-flex flex-column gap-2">
                <Form.Check
                  type="radio"
                  id="user-degree-none"
                  name="degreeId"
                  label="Sin carrera asignada"
                  value=""
                  checked={!data.degreeId}
                  onChange={handleDegreeChange}
                  disabled={readOnly}
                />
                {degreeOptions.length === 0 ? (
                  <span className="text-muted small">No hay carreras activas cargadas.</span>
                ) : (
                  degreeOptions.map((d) => {
                    const hid = degRowId(d);
                    if (!hid) return null;
                    return (
                      <Form.Check
                        key={hid}
                        type="radio"
                        id={`user-degree-${hid}`}
                        name="degreeId"
                        label={d.name || hid}
                        value={hid}
                        checked={String(data.degreeId) === String(hid)}
                        onChange={handleDegreeChange}
                        disabled={readOnly}
                      />
                    );
                  })
                )}
              </div>
              <Form.Text className="text-muted">Las carreras se administran en Administración → Carreras.</Form.Text>
            </Form.Group>
          ) : null}

          {roleUpper === 'ALUMNO' ? (
            <Form.Group className="mt-3" controlId="userModalidadCondicion">
              <Form.Label className="small text-muted mb-1">Ficha académica del alumno</Form.Label>
              <div className="border rounded p-2 mt-1">
                <Form.Label className="small text-muted mb-1">Modalidad</Form.Label>
                <Form.Check
                  type="radio"
                  id="user-modalidad-presencial"
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
                  id="user-modalidad-remoto"
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
                  id="user-condicion-regular"
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
                  id="user-condicion-libre"
                  name="condicion"
                  value="LIBRE"
                  label="Libre"
                  checked={data.condicion === 'LIBRE'}
                  onChange={() => setData((prev) => ({ ...prev, condicion: 'LIBRE' }))}
                  disabled={readOnly || data.active === false}
                />
              </div>
              <Form.Text className="text-muted">
                Se guardan en la ficha del alumno (colección alumnos). Si hay ficha vinculada por email, son obligatorios
                mientras el usuario esté activo.
              </Form.Text>
            </Form.Group>
          ) : null}
        </div>
      </Form.Group>

      {data.id && unlocked ? (
        <Form.Group className="mb-3" controlId="userEstado">
          <Form.Label>Estado</Form.Label>
          <div className="border rounded p-2">
            <Form.Check
              inline
              type="radio"
              id="user-active-si"
              name="userEstadoActivo"
              label="Activo"
              checked={data.active !== false}
              onChange={() => openActiveConfirm(true)}
            />
            <Form.Check
              inline
              type="radio"
              id="user-active-no"
              name="userEstadoActivo"
              label="Inactivo"
              checked={data.active === false}
              onChange={() => openActiveConfirm(false)}
            />
          </div>
          <Form.Text className="text-muted">El cambio requiere confirmación.</Form.Text>
        </Form.Group>
      ) : null}

      {formError ? <div className="text-danger small mb-2">{formError}</div> : null}

      <Form.Group className="mb-3" controlId="userPassword">
        <Form.Label>{data.id ? 'Nueva contraseña (opcional)' : 'Contraseña'}</Form.Label>
        <InputGroup>
          <Form.Control
            type={showPassword ? 'text' : 'password'}
            name="password"
            autoComplete="new-password"
            placeholder={data.id ? 'Vacío = no cambiar' : 'Mínimo 6 caracteres'}
            onChange={handleInputChange}
            value={data.password}
            readOnly={readOnly}
          />
          <Button
            variant="outline-secondary"
            type="button"
            title={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            aria-label={showPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            onClick={() => setShowPassword((v) => !v)}
            tabIndex={-1}
            disabled={readOnly}
          >
            <FontAwesomeIcon icon={showPassword ? faEyeSlash : faEye} />
          </Button>
        </InputGroup>
      </Form.Group>

      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {data.id ? 'Actualizar' : 'Guardar'}
      </Button>

      <ConfirmChangeEstadoModal
        show={!!activeFormConfirm}
        onHide={() => !activeFormSaving && setActiveFormConfirm(null)}
        kind="usuario"
        itemName={data.name || data.email || ''}
        fromActive={activeFormConfirm ? activeFormConfirm.fromActive : true}
        toActive={activeFormConfirm ? activeFormConfirm.toActive : true}
        onConfirm={confirmActiveForm}
        confirming={activeFormSaving}
      />
    </Form>
  );
};

export default FormUser;
