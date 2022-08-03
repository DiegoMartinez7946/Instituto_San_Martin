import React from 'react';
import { Modal } from 'react-bootstrap';

import FormRole from '../Form/FormRole';

const ModalRole = ({show, handleClose, saveEvent, data}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>{ data ? 'Editar Rol' : 'Agregar Rol'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormRole 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalRole;
