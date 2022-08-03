import React from 'react';
import { Modal } from 'react-bootstrap';

import FormDegree from '../Form/FormDegree';

const ModalDegree = ({show, handleClose, saveEvent, data}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>{ data ? 'Editar Carrera' : 'Agregar Carrera'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormDegree 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalDegree;
