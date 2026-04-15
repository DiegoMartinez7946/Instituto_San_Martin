/** DNI: exactamente 7 u 8 caracteres, solo dígitos (sin puntos ni letras). */
export const validateDNI = (dni) => {
  const d = dni != null ? String(dni).trim() : '';
  if (!d) {
    return { ok: false, message: 'El DNI es obligatorio' };
  }
  if (!/^\d{7,8}$/.test(d)) {
    return {
      ok: false,
      message: 'El DNI debe tener 7 u 8 dígitos numéricos (sin puntos ni letras)'
    };
  }
  return { ok: true, message: '' };
};
