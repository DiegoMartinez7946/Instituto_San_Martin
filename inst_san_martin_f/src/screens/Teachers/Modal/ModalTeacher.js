import React from 'react';
import { Modal } from 'react-bootstrap';

import FormTeacher from '../Form/FormTeacher';

const ModalTeacher = ({
  show,
  handleClose,
  saveEvent,
  data,
  degrees,
  titulosHabilitantes,
  modalidades
}) => {
  const eventHandler = (e) => {
    saveEvent(e);
  };

  const isEdit = data && typeof data === 'object' && data.id;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg">
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar docente' : 'Agregar docente'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormTeacher
          key={isEdit ? data.id : 'new'}
          dataEntry={data}
          degrees={degrees}
          titulosHabilitantes={titulosHabilitantes}
          modalidades={modalidades}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalTeacher;
