import React from 'react';
import { Modal, Button } from 'react-bootstrap';

const entityPhrase = (kind) => {
  if (kind === 'docente') return 'del docente';
  if (kind === 'alumno') return 'del alumno';
  if (kind === 'usuario') return 'del usuario';
  return 'de la carrera';
};

const ConfirmChangeEstadoModal = ({
  show,
  onHide,
  kind,
  itemName,
  fromActive,
  toActive,
  onConfirm,
  confirming
}) => {
  const fromTxt = fromActive ? 'activo' : 'inactivo';
  const toTxt = toActive ? 'activo' : 'inactivo';
  const phrase = entityPhrase(kind || 'carrera');

  return (
    <Modal show={show} onHide={onHide} centered backdrop={confirming ? 'static' : true}>
      <Modal.Header closeButton={!confirming}>
        <Modal.Title>Confirmar cambio de estado</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <p className="mb-0">
          Está por cambiar el estado {phrase}
          {itemName ? ` «${itemName}»` : ''} de <strong>{fromTxt}</strong> a <strong>{toTxt}</strong>.
        </p>
        <p className="text-muted small mt-2 mb-0">Confirme para guardar el cambio.</p>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onHide} disabled={confirming}>
          Cancelar
        </Button>
        <Button variant="primary" onClick={onConfirm} disabled={confirming}>
          {confirming ? 'Guardando…' : 'OK'}
        </Button>
      </Modal.Footer>
    </Modal>
  );
};

export default ConfirmChangeEstadoModal;
