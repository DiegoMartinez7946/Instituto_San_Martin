import React from 'react';
import { Button } from 'react-bootstrap';

/**
 * @param {string} entityKey - Si no hay valor, no se muestra el banner.
 * @param {boolean|undefined} estadoActivo - Si es boolean, se muestra el estado activo/inactivo.
 */
const FormEditLockBanner = ({
  entityKey,
  estadoActivo,
  unlocked,
  onUnlock,
  onCancelUnlock,
  unlockVariant = 'outline-primary'
}) => {
  if (!entityKey) return null;

  const hasEstado = typeof estadoActivo === 'boolean';

  return (
    <div className="border rounded p-3 mb-3 bg-light">
      <div className="d-flex flex-wrap align-items-start justify-content-between gap-3">
        <div className="flex-grow-1" style={{ minWidth: '200px' }}>
          {hasEstado ? (
            <>
              <div className="text-muted small">Estado del registro</div>
              <div className={`fw-semibold ${estadoActivo ? 'text-success' : 'text-secondary'}`}>
                {estadoActivo ? 'Activo' : 'Inactivo'}
              </div>
              {!unlocked ? (
                <p className="text-muted small mb-0 mt-2">
                  Los datos están protegidos. Pulse Desbloquear para modificar el formulario
                  {hasEstado ? ' o el estado.' : '.'}
                </p>
              ) : null}
            </>
          ) : (
            <p className="text-muted small mb-0">
              {unlocked
                ? 'Puede editar los campos.'
                : 'Los datos están protegidos. Pulse Desbloquear para poder modificarlos.'}
            </p>
          )}
        </div>
        <div className="d-flex align-items-center gap-2">
          {unlocked && typeof onCancelUnlock === 'function' ? (
            <Button type="button" variant="danger" onClick={onCancelUnlock}>
              Cancelar
            </Button>
          ) : null}
          <Button type="button" variant={unlockVariant} disabled={unlocked} onClick={onUnlock}>
            Desbloquear
          </Button>
        </div>
      </div>
    </div>
  );
};

export default FormEditLockBanner;
