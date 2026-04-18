import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { useGlobal } from '../../../context/Global/GlobalProvider';
import { isSaveSuccess } from '../../../utils/saveResult';
import { NIVELES_ORDENADOS } from '../../../constant/nivelesAcademicos';

/** Códigos persistidos en API/DB (minúsculas). */
const TURNOS_CARRERA = [
  { code: 'manana', label: 'Mañana' },
  { code: 'tarde', label: 'Tarde' },
  { code: 'noche', label: 'Noche' }
];

const defaultTurnosCarrera = () => TURNOS_CARRERA.map((t) => t.code);

const normalizeTurnosFromEntry = (entry) => {
  const allowed = new Set(TURNOS_CARRERA.map((t) => t.code));
  const raw =
    entry && Array.isArray(entry.turnos)
      ? entry.turnos
      : entry && Array.isArray(entry.Turnos)
        ? entry.Turnos
        : [];
  const pickedSet = new Set(
    raw.map((x) => String(x || '').toLowerCase().trim()).filter((c) => allowed.has(c))
  );
  if (pickedSet.size) {
    return TURNOS_CARRERA.map((t) => t.code).filter((c) => pickedSet.has(c));
  }
  return defaultTurnosCarrera();
};

const FormDegree = ({ dataEntry, saveData, changeActive }) => {
  const [globalState] = useGlobal();
  const [data, setData] = useState({
    id: '',
    name: '',
    nivel: '',
    studyPlanId: '',
    turnos: defaultTurnosCarrera(),
    active: true
  });
  const [turnosError, setTurnosError] = useState('');
  const [activeFormConfirm, setActiveFormConfirm] = useState(null);
  const [activeFormSaving, setActiveFormSaving] = useState(false);

  const lockEntityKey =
    dataEntry && typeof dataEntry === 'object' ? dataEntry.id || dataEntry.ID || '' : '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : {};
    setData({
      id: entry.id || entry.ID || '',
      name: entry.name || '',
      nivel: (entry.nivel || '').toLowerCase(),
      studyPlanId: (() => {
        const x = entry.studyPlanId != null ? String(entry.studyPlanId).trim() : '';
        if (!x || x === '000000000000000000000000' || /^[a-f0-9]{24}$/i.test(x)) {
          return '';
        }
        return x;
      })(),
      turnos: normalizeTurnosFromEntry(entry),
      active: entry.active !== false
    });
    setTurnosError('');
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
      setData((prev) => ({ ...prev, active: activeFormConfirm.toActive }));
      setActiveFormConfirm(null);
    } else if (code === 199) {
      setActiveFormConfirm(null);
    }
  };

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    setData((prev) => ({
      ...prev,
      [name]: value
    }));
  };

  const toggleTurno = (code) => {
    setTurnosError('');
    setData((prev) => {
      const set = new Set(prev.turnos || []);
      if (set.has(code)) {
        set.delete(code);
      } else {
        set.add(code);
      }
      return { ...prev, turnos: TURNOS_CARRERA.map((t) => t.code).filter((c) => set.has(c)) };
    });
  };

  const sendData = async (e) => {
    e.preventDefault();
    if (!(data.turnos && data.turnos.length)) {
      setTurnosError('Seleccione al menos un turno.');
      return;
    }
    setTurnosError('');
    const payload = {
      ...data,
      id: data.id != null && data.id !== '' ? String(data.id) : data.id,
      turnos: Array.isArray(data.turnos) ? data.turnos : defaultTurnosCarrera(),
      studyPlanId: data.studyPlanId != null && String(data.studyPlanId).trim() !== '' ? String(data.studyPlanId).trim() : ''
    };
    const res = await saveData(payload);
    if (isSaveSuccess(res)) {
      armLockAfterSave();
    }
  };

  const hasId = !!(data.id || data.ID);

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicID">
        <Form.Control type="hidden" disabled name="id" value={data ? data.id : ''} />
      </Form.Group>

      <FormEditLockBanner
        entityKey={lockEntityKey}
        estadoActivo={lockEntityKey ? data.active !== false : undefined}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
        onCancelUnlock={() => setUnlocked(false)}
        unlockVariant="warning"
      />

      <Form.Group className="mb-3" controlId="formBasicName">
        <Form.Label>Nombre</Form.Label>
        <Form.Control
          type="text"
          placeholder="Ingrese nombre"
          name="name"
          onChange={handleInputChange}
          value={data ? data.name : ''}
          autoFocus
          readOnly={readOnly}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formDegreeNivel">
        <Form.Label>Nivel</Form.Label>
        <Form.Select
          name="nivel"
          value={data.nivel}
          onChange={handleInputChange}
          required
          disabled={readOnly}
        >
          <option value="">Seleccione nivel</option>
          {NIVELES_ORDENADOS.map((n) => (
            <option key={n} value={n}>
              {n.charAt(0).toUpperCase() + n.slice(1)}
            </option>
          ))}
        </Form.Select>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formDegreeStudyPlan">
        <Form.Label>Plan de estudio</Form.Label>
        <Form.Select
          name="studyPlanId"
          value={data.studyPlanId || ''}
          onChange={handleInputChange}
          disabled={readOnly}
        >
          <option value="">Sin plan asociado</option>
          {(Array.isArray(globalState.studyPlans) ? globalState.studyPlans : [])
            .filter((p) => p.active !== false && String(p.numeroResolucion || '').trim())
            .map((p) => {
              const nr = String(p.numeroResolucion).trim();
              const label = [p.nombre, nr].filter(Boolean).join(' — ');
              return (
                <option key={p.id || p.ID} value={nr}>
                  {label || nr}
                </option>
              );
            })}
        </Form.Select>
        <Form.Text className="text-muted">
          Solo planes activos con número de resolución; ese valor queda vinculado a la carrera.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formDegreeTurnos">
        <Form.Label>Turnos en que se dicta</Form.Label>
        <div className="border rounded p-2">
          {TURNOS_CARRERA.map(({ code, label }) => (
            <Form.Check
              key={code}
              type="checkbox"
              id={`degree-turno-${code}`}
              label={label}
              checked={(data.turnos || []).includes(code)}
              disabled={readOnly}
              onChange={() => toggleTurno(code)}
            />
          ))}
        </div>
        {turnosError ? <div className="text-danger small mt-1">{turnosError}</div> : null}
        <Form.Text className="text-muted">
          Puede marcar más de uno. Indica en qué franjas horarias aplica la carrera.
        </Form.Text>
      </Form.Group>

      {hasId && unlocked ? (
        <Form.Group className="mb-3" controlId="degreeEstado">
          <Form.Label>Estado</Form.Label>
          <div>
            <Form.Check
              inline
              type="radio"
              id="degree-active-si"
              name="degreeEstadoActivo"
              label="Activo"
              checked={data.active !== false}
              onChange={() => openActiveConfirm(true)}
            />
            <Form.Check
              inline
              type="radio"
              id="degree-active-no"
              name="degreeEstadoActivo"
              label="Inactivo"
              checked={data.active === false}
              onChange={() => openActiveConfirm(false)}
            />
          </div>
          <Form.Text className="text-muted">El cambio requiere confirmación.</Form.Text>
        </Form.Group>
      ) : null}

      <br />
      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {hasId ? 'Actualizar' : 'Guardar'}
      </Button>

      <ConfirmChangeEstadoModal
        show={!!activeFormConfirm}
        onHide={() => !activeFormSaving && setActiveFormConfirm(null)}
        kind="carrera"
        itemName={data.name || ''}
        fromActive={activeFormConfirm ? activeFormConfirm.fromActive : true}
        toActive={activeFormConfirm ? activeFormConfirm.toActive : true}
        onConfirm={confirmActiveForm}
        confirming={activeFormSaving}
      />
    </Form>
  );
};

export default FormDegree;
