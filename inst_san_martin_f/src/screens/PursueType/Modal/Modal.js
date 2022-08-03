import React from 'react';
import { Modal } from 'react-bootstrap';

import FormPursueType from '../Form/Form';

const ModalPursueType = ({show, handleClose, saveEvent, data}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>{ data ? 'Editar Modalidad' : 'Agregar Modalidad'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormPursueType 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalPursueType;
