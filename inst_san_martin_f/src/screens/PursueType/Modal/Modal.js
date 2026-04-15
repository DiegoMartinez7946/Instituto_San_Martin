import React from 'react';
import { Modal } from 'react-bootstrap';

import FormPursueType from '../Form/Form';

const ModalPursueType = ({ show, handleClose, saveEvent, data }) => {
  const eventHandler = async (e) => saveEvent(e);

  const isEdit = data && typeof data === 'object' && data.ID;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar Modalidad' : 'Agregar Modalidad'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormPursueType 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalPursueType;
