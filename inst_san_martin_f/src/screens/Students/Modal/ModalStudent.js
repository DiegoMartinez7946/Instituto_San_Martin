import React from 'react';
import { Modal } from 'react-bootstrap';

import FormStudent from '../Form/FormStudent';

const ModalStudent = ({ show, handleClose, saveEvent, data, degrees }) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  const isEdit = data && typeof data === 'object' && data.id;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg">
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar alumno' : 'Agregar alumno'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormStudent
          key={isEdit ? data.id : 'new'}
          dataEntry={data}
          degrees={degrees}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalStudent;
