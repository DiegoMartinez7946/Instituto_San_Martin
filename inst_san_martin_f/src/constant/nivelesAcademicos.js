/** Debe coincidir con la colección `jerarquia` del backend */
export const JERARQUIA_POR_NIVEL = {
  inicial: 1,
  primario: 2,
  secundario: 3,
  terciario: 4,
  universitario: 4
};

export const NIVELES_ORDENADOS = Object.keys(JERARQUIA_POR_NIVEL);

export const jerarquiaDeNivel = (nivel) => {
  const k = nivel != null ? String(nivel).toLowerCase().trim() : '';
  if (!k) return null;
  return Object.prototype.hasOwnProperty.call(JERARQUIA_POR_NIVEL, k) ? JERARQUIA_POR_NIVEL[k] : null;
};

export const etiquetaNivel = (nivel) => {
  if (!nivel) return '';
  return nivel.charAt(0).toUpperCase() + nivel.slice(1).toLowerCase();
};
