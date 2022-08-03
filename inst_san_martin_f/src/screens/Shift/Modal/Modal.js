import React from 'react';
import { Modal } from 'react-bootstrap';

import FormShift from '../Form/Form';

const ModalShift = ({show, handleClose, saveEvent, data}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>{ data ? 'Editar Turno' : 'Agregar Turno'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormShift 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalShift;
