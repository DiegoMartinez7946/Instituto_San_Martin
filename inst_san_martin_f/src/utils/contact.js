const EMAIL_RE = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

/** Correo opcional: vacío OK; si hay texto, formato con @ y dominio con extensión. */
export const validateCorreoElectronicoOpcional = (email) => {
  const e = email != null ? String(email).trim() : '';
  if (!e) {
    return { ok: true, message: '' };
  }
  if (e.length > 254) {
    return { ok: false, message: 'El correo electrónico es demasiado largo' };
  }
  if (!EMAIL_RE.test(e)) {
    return {
      ok: false,
      message: 'El correo electrónico debe ser válido (incluir @ y dominio con extensión)'
    };
  }
  return { ok: true, message: '' };
};

/** Correo obligatorio y formato válido. */
export const validateCorreoElectronicoRequerido = (email) => {
  const e = email != null ? String(email).trim() : '';
  if (!e) {
    return { ok: false, message: 'El correo electrónico es obligatorio' };
  }
  return validateCorreoElectronicoOpcional(e);
};

/** Teléfono opcional: vacío OK; si hay dígitos, solo números y entre 7 y 15 (nacional/internacional). */
export const validateTelefonoOpcional = (digits) => {
  const p = digits != null ? String(digits).trim() : '';
  if (!p) {
    return { ok: true, message: '' };
  }
  if (!/^\d+$/.test(p)) {
    return {
      ok: false,
      message: 'El teléfono solo puede contener números (sin letras ni caracteres especiales)'
    };
  }
  if (p.length < 7 || p.length > 15) {
    return {
      ok: false,
      message: 'El teléfono debe tener entre 7 y 15 dígitos'
    };
  }
  return { ok: true, message: '' };
};

/** Deja solo dígitos (máx. 15) para el valor del campo teléfono. */
export const sanitizeTelefonoInput = (value) =>
  String(value != null ? value : '').replace(/\D/g, '').slice(0, 15);
