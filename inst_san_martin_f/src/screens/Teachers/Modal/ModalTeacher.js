import React from 'react';
import { Modal } from 'react-bootstrap';

import FormTeacher from '../Form/FormTeacher';

const ModalTeacher = ({
  show,
  handleClose,
  saveEvent,
  data,
  degrees,
  shifts,
  titulosHabilitantes,
  modalidades,
  changeActive
}) => {
  const eventHandler = async (e) => saveEvent(e);

  const isEdit = data && typeof data === 'object' && data.id;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar docente' : 'Agregar docente'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormTeacher
          key={isEdit ? data.id : 'new'}
          dataEntry={data}
          degrees={degrees}
          shifts={shifts}
          titulosHabilitantes={titulosHabilitantes}
          modalidades={modalidades}
          saveData={(e) => eventHandler(e)}
          changeActive={changeActive}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalTeacher;
