import React, { useState, useEffect, useMemo } from 'react';
import { Form, Button } from 'react-bootstrap';

import ConfirmChangeEstadoModal from '../../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';
import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';

import styles from './FormStudyPlan.module.css';

const DIGIT_OPTS = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

const parseIntSafe = (v, def = 0) => {
  const n = parseInt(String(v).replace(/\D/g, ''), 10);
  return Number.isFinite(n) ? n : def;
};

/** Solo dígitos; máximo 4 + barra + 2 (ej. 6555/24). */
const maskNumeroResolucion = (raw) => {
  const d = String(raw || '')
    .replace(/\D/g, '')
    .slice(0, 6);
  if (d.length <= 4) return d;
  return `${d.slice(0, 4)}/${d.slice(4)}`;
};

/** Fin del año de caducidad (31/12 inclusive). */
const endOfCalendarYear = (year) => {
  if (!Number.isFinite(year) || year < 1000) return null;
  return new Date(year, 11, 31, 23, 59, 59, 999);
};

/**
 * Periodo de extensión: después del 31/12 del año (aprobación + cohorte) y hasta el fin del año de caducidad.
 * Sin años de extensión (0) nunca aplica.
 */
const isEnPeriodoExtension = (anioAprobacion, cohorte, extension, now = new Date()) => {
  const ext = Number(extension);
  if (!Number.isFinite(ext) || ext <= 0) return false;
  const apr = Number(anioAprobacion);
  const coh = Number(cohorte);
  if (!Number.isFinite(apr) || apr < 1000 || !Number.isFinite(coh)) return false;
  const finCohorte = endOfCalendarYear(apr + coh);
  const finTotal = endOfCalendarYear(apr + coh + ext);
  if (!finCohorte || !finTotal) return false;
  if (now > finTotal) return false;
  return now > finCohorte;
};

/** Tiempo aproximado hasta el 31/12 del año de caducidad. */
const yearsMonthsHastaCaducidad = (anioAprobacion, cohorte, extension, now = new Date()) => {
  const apr = Number(anioAprobacion);
  const coh = Number(cohorte);
  const ext = Number(extension);
  if (!Number.isFinite(apr) || apr < 1000 || !Number.isFinite(coh) || !Number.isFinite(ext)) {
    return { label: '—', expired: false };
  }
  const end = endOfCalendarYear(apr + coh + ext);
  if (!end) return { label: '—', expired: false };
  if (now > end) return { label: 'Plan caducado', expired: true };
  let rm = (end.getFullYear() - now.getFullYear()) * 12 + (end.getMonth() - now.getMonth());
  if (end.getDate() < now.getDate()) rm -= 1;
  if (rm < 0) return { label: 'Plan caducado', expired: true };
  const y = Math.floor(rm / 12);
  const mo = rm % 12;
  const parts = [];
  if (y > 0) parts.push(`${y} año${y === 1 ? '' : 's'}`);
  if (mo > 0) parts.push(`${mo} mes${mo === 1 ? '' : 'es'}`);
    if (!parts.length) parts.push('Menos de un mes');
  return { label: parts.join(' y '), expired: false };
};

const FormStudyPlan = ({ dataEntry, saveData, changeActive }) => {
  const [data, setData] = useState({
    id: '',
    nombre: '',
    numeroResolucion: '',
    cohorte: 0,
    extension: 0,
    anioAprobacion: new Date().getFullYear(),
    cargaHoraria: 1,
    active: true
  });
  const [activeFormConfirm, setActiveFormConfirm] = useState(null);
  const [activeFormSaving, setActiveFormSaving] = useState(false);

  const lockEntityKey =
    dataEntry && typeof dataEntry === 'object' ? dataEntry.id || dataEntry.ID || '' : '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  const anioCaducidad = useMemo(() => {
    const a = parseIntSafe(data.anioAprobacion, 0);
    const c = parseIntSafe(data.cohorte, 0);
    const e = parseIntSafe(data.extension, 0);
    if (a < 1000) return '—';
    return String(a + c + e);
  }, [data.anioAprobacion, data.cohorte, data.extension]);

  const tiempoHastaCaducidad = useMemo(
    () =>
      yearsMonthsHastaCaducidad(
        parseIntSafe(data.anioAprobacion, 0),
        parseIntSafe(data.cohorte, 0),
        parseIntSafe(data.extension, 0)
      ),
    [data.anioAprobacion, data.cohorte, data.extension]
  );

  const enPeriodoExtension = useMemo(
    () =>
      isEnPeriodoExtension(
        parseIntSafe(data.anioAprobacion, 0),
        parseIntSafe(data.cohorte, 0),
        parseIntSafe(data.extension, 0)
      ),
    [data.anioAprobacion, data.cohorte, data.extension]
  );

  useEffect(() => {
    const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : {};
    setData({
      id: entry.id || entry.ID || '',
      nombre: entry.nombre != null ? String(entry.nombre) : '',
      numeroResolucion: entry.numeroResolucion != null ? String(entry.numeroResolucion) : '',
      cohorte: entry.cohorte != null ? parseIntSafe(entry.cohorte, 0) : 0,
      extension: entry.extension != null ? parseIntSafe(entry.extension, 0) : 0,
      anioAprobacion: entry.anioAprobacion != null ? parseIntSafe(entry.anioAprobacion, new Date().getFullYear()) : new Date().getFullYear(),
      cargaHoraria: entry.cargaHoraria != null ? parseIntSafe(entry.cargaHoraria, 1) : 1,
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
    if (Number(res.code) === 200) {
      setData((prev) => ({ ...prev, active: activeFormConfirm.toActive }));
      setActiveFormConfirm(null);
    } else if (Number(res.code) === 199) {
      setActiveFormConfirm(null);
    }
  };

  const handleInputChange = (event) => {
    const { name, value } = event.target;
    if (name === 'numeroResolucion') {
      setData((prev) => ({ ...prev, numeroResolucion: maskNumeroResolucion(value) }));
      return;
    }
    setData((prev) => ({ ...prev, [name]: value }));
  };

  const handleDigitChange = (name) => (e) => {
    const v = e.target.value;
    setData((prev) => ({ ...prev, [name]: v === '' ? '' : parseIntSafe(v, 0) }));
  };

  const sendData = async (e) => {
    e.preventDefault();
    const cohorte = Math.min(9, Math.max(0, parseIntSafe(data.cohorte, 0)));
    const extension = Math.min(9, Math.max(0, parseIntSafe(data.extension, 0)));
    const anioAprobacion = parseIntSafe(data.anioAprobacion, 0);
    const cargaHoraria = parseIntSafe(data.cargaHoraria, 0);
    const anioCaducidadNum = anioAprobacion + cohorte + extension;
    const payload = {
      ...data,
      id: data.id != null && data.id !== '' ? String(data.id) : data.id,
      nombre: String(data.nombre || '').trim(),
      numeroResolucion: String(data.numeroResolucion || '').trim(),
      cohorte,
      extension,
      anioAprobacion,
      anioCaducidad: anioCaducidadNum,
      cargaHoraria
    };
    const res = await saveData(payload);
    if (isSaveSuccess(res)) {
      armLockAfterSave();
    }
  };

  const hasId = !!(data.id || data.ID);

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3">
        <Form.Control type="hidden" disabled name="id" value={data.id || ''} />
      </Form.Group>

      <FormEditLockBanner
        entityKey={lockEntityKey}
        estadoActivo={lockEntityKey ? data.active !== false : undefined}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
        onCancelUnlock={() => setUnlocked(false)}
        unlockVariant="warning"
      />

      <Form.Group className="mb-3" controlId="spNombrePlan">
        <Form.Label>Nombre del plan de estudio</Form.Label>
        <Form.Control
          type="text"
          name="nombre"
          placeholder="Ej. Técnico en Electromecánica (Actualización 2024)"
          value={data.nombre}
          onChange={handleInputChange}
          required
          readOnly={readOnly}
          maxLength={200}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="spNumeroResolucion">
        <Form.Label>Número de resolución</Form.Label>
        <Form.Control
          type="text"
          name="numeroResolucion"
          inputMode="numeric"
          autoComplete="off"
          placeholder="6555/24"
          value={data.numeroResolucion}
          onChange={handleInputChange}
          required
          readOnly={readOnly}
          maxLength={7}
          pattern="\d{4}/\d{2}"
          title="Cuatro dígitos, barra, dos dígitos (año de aprobación en dos cifras)"
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="spCohorte">
        <Form.Label>Cohorte (años de validez)</Form.Label>
        <Form.Select
          name="cohorte"
          value={String(data.cohorte)}
          onChange={(ev) => setData((p) => ({ ...p, cohorte: parseIntSafe(ev.target.value, 0) }))}
          disabled={readOnly}
        >
          {DIGIT_OPTS.map((d) => (
            <option key={d} value={String(d)}>
              {d}
            </option>
          ))}
        </Form.Select>
      </Form.Group>

      <Form.Group className="mb-3" controlId="spExtension">
        <Form.Label>Extensión (años adicionales)</Form.Label>
        <Form.Select
          name="extension"
          value={String(data.extension)}
          onChange={(ev) => setData((p) => ({ ...p, extension: parseIntSafe(ev.target.value, 0) }))}
          disabled={readOnly}
        >
          {DIGIT_OPTS.map((d) => (
            <option key={d} value={String(d)}>
              {d}
            </option>
          ))}
        </Form.Select>
      </Form.Group>

      <Form.Group className="mb-3" controlId="spAnioAprobacion">
        <Form.Label>Año de aprobación</Form.Label>
        <Form.Control
          type="number"
          name="anioAprobacion"
          min={1000}
          max={9999}
          value={data.anioAprobacion}
          onChange={handleDigitChange('anioAprobacion')}
          required
          readOnly={readOnly}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="spAnioCaducidad">
        <Form.Label>Año de caducidad del plan</Form.Label>
        <Form.Control type="text" value={anioCaducidad} readOnly disabled />
        <Form.Text className="text-muted">Año de aprobación + validez + extensión (calculado).</Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="spTiempoRestante">
        <Form.Label>Tiempo restante hasta la caducidad</Form.Label>
        <Form.Control
          type="text"
          readOnly
          disabled
          value={tiempoHastaCaducidad.label}
          className={tiempoHastaCaducidad.expired ? 'text-danger' : undefined}
        />
        <Form.Text className="text-muted">
          Calculado hasta el 31 de diciembre del año de caducidad según cohorte y extensión.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="spPeriodoExtension">
        <Form.Label className="mb-2 d-block">¿En periodo de extensión?</Form.Label>
        <span className="visually-hidden" aria-live="polite">
          {enPeriodoExtension
            ? 'En periodo de extensión: sí.'
            : 'En periodo de extensión: no.'}
        </span>
        {/* Sin inputs nativos: Bootstrap fuerza azul en radios; aquí el color es 100 % nuestro */}
        <div
          className={styles.extensionRadios}
          role="radiogroup"
          aria-readonly="true"
          aria-label="Periodo de extensión, solo lectura"
        >
          <div className={styles.extensionOption}>
            <span
              className={`${styles.radioDot} ${
                !enPeriodoExtension ? styles.radioDotNoOn : styles.radioDotIdle
              }`}
              aria-hidden
            />
            <span
              className={
                !enPeriodoExtension ? styles.extensionLabelNo : styles.extensionLabelMuted
              }
            >
              NO
            </span>
          </div>
          <div className={styles.extensionOption}>
            <span
              className={`${styles.radioDot} ${
                enPeriodoExtension ? styles.radioDotSiOn : styles.radioDotIdle
              }`}
              aria-hidden
            />
            <span
              className={
                enPeriodoExtension ? styles.extensionLabelSi : styles.extensionLabelMuted
              }
            >
              SI
            </span>
          </div>
        </div>
        <Form.Text className="text-muted">
          Solo lectura: se actualiza según la fecha actual y los años de validez y extensión.
        </Form.Text>
      </Form.Group>

      <Form.Group className="mb-3" controlId="spCargaHoraria">
        <Form.Label>Carga horaria (horas)</Form.Label>
        <Form.Control
          type="number"
          name="cargaHoraria"
          min={1}
          max={99999}
          value={data.cargaHoraria}
          onChange={handleDigitChange('cargaHoraria')}
          required
          readOnly={readOnly}
        />
      </Form.Group>

      {hasId && unlocked ? (
        <Form.Group className="mb-3" controlId="spEstado">
          <Form.Label>Estado</Form.Label>
          <div>
            <Form.Check
              inline
              type="radio"
              id="sp-active-si"
              name="spEstadoActivo"
              label="Activo"
              checked={data.active !== false}
              onChange={() => openActiveConfirm(true)}
            />
            <Form.Check
              inline
              type="radio"
              id="sp-active-no"
              name="spEstadoActivo"
              label="Inactivo"
              checked={data.active === false}
              onChange={() => openActiveConfirm(false)}
            />
          </div>
          <Form.Text className="text-muted">El cambio requiere confirmación.</Form.Text>
        </Form.Group>
      ) : null}

      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {hasId ? 'Guardar cambios' : 'Guardar'}
      </Button>

      <ConfirmChangeEstadoModal
        show={!!activeFormConfirm}
        onHide={() => !activeFormSaving && setActiveFormConfirm(null)}
        kind="planestudio"
        itemName={data.nombre || data.numeroResolucion || ''}
        fromActive={activeFormConfirm ? activeFormConfirm.fromActive : true}
        toActive={activeFormConfirm ? activeFormConfirm.toActive : true}
        onConfirm={confirmActiveForm}
        confirming={activeFormSaving}
      />
    </Form>
  );
};

export default FormStudyPlan;
