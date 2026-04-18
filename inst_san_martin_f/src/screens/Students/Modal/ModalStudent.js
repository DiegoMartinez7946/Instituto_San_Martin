import React from 'react';
import { Modal } from 'react-bootstrap';

import FormStudent from '../Form/FormStudent';

const ModalStudent = ({ show, handleClose, saveEvent, data, degrees, shifts, changeActive }) => {

  const eventHandler = async (e) => saveEvent(e);

  const isEdit =
    data &&
    typeof data === 'object' &&
    !!(data.id || data.idAlumno);

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar alumno' : 'Agregar alumno'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormStudent
          key={isEdit ? data.id || data.idAlumno : 'new'}
          dataEntry={data}
          degrees={degrees}
          shifts={shifts}
          saveData={(e) => eventHandler(e)}
          changeActive={changeActive}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalStudent;
