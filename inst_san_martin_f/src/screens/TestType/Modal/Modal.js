import React from 'react';
import { Modal } from 'react-bootstrap';

import FormTestType from '../Form/Form';

const ModalTestType = ({ show, handleClose, saveEvent, data }) => {
  const eventHandler = async (e) => saveEvent(e);

  const isEdit = data && typeof data === 'object' && data.ID;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar Tipo de Examen' : 'Agregar Tipo de Examen'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormTestType 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalTestType;
