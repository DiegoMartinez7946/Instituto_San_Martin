import React from 'react';
import { Modal } from 'react-bootstrap';

import FormUser from '../Form/FormUser';

const ModalRole = ({show, handleClose, saveEvent}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>Agregar Usuario</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormUser 
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
      <Modal.Footer>
      </Modal.Footer>
    </Modal>
  );
};

export default ModalRole;
