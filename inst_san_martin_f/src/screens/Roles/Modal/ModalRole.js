import React from 'react';
import { Modal } from 'react-bootstrap';

import FormRole from '../Form/FormRole';

const ModalRole = ({ show, handleClose, saveEvent, data }) => {
  const eventHandler = async (e) => saveEvent(e);

  const isEdit = data && typeof data === 'object' && data.ID;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar Rol' : 'Agregar Rol'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormRole 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalRole;
