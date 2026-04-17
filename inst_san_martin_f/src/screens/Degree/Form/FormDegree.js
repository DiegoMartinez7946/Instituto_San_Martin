import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';
import { NIVELES_ORDENADOS } from '../../../constant/nivelesAcademicos';

const FormDegree = ({ dataEntry, saveData, changeActive }) => {
  const [data, setData] = useState({
    id: '',
    name: '',
    nivel: '',
    resolucionId: '',
    active: true
  });
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
      resolucionId: entry.resolucionId != null ? String(entry.resolucionId) : '',
      active: entry.active !== false
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

  const sendData = async (e) => {
    e.preventDefault();
    const res = await saveData(data);
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

      <Form.Group className="mb-3" controlId="formDegreeResolucion">
        <Form.Label>Resolución ID</Form.Label>
        <Form.Control
          type="text"
          name="resolucionId"
          placeholder="Texto alfanumérico (puede incluir caracteres especiales)"
          onChange={handleInputChange}
          value={data.resolucionId}
          required
          readOnly={readOnly}
        />
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
