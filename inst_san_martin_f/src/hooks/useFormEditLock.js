import { useState, useEffect, useCallback } from 'react';

/**
 * Bloqueo de edición para registros ya persistidos.
 * @param {string} entityKey - Cadena no vacía = registro existente (bloquear hasta desbloquear). Vacío = alta nueva sin bloqueo.
 * @param {*} syncKey - Al cambiar (p. ej. dataEntry), se vuelve a bloquear.
 */
export function useFormEditLock(entityKey, syncKey) {
  const isLockedEntity = Boolean(entityKey);
  const [unlocked, setUnlocked] = useState(false);

  useEffect(() => {
    setUnlocked(false);
  }, [entityKey, syncKey]);

  const readOnly = isLockedEntity && !unlocked;

  const armLockAfterSave = useCallback(() => {
    setUnlocked(false);
  }, []);

  return {
    isLockedEntity,
    unlocked,
    setUnlocked,
    readOnly,
    armLockAfterSave
  };
}
