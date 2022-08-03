import React from 'react';
import { Modal } from 'react-bootstrap';

import FormTestType from '../Form/Form';

const ModalTestType = ({show, handleClose, saveEvent, data}) => {

  const eventHandler = (e) => {
    saveEvent(e);
  };

  return (
    <Modal show={show} onHide={handleClose} centered>
      <Modal.Header closeButton>
        <Modal.Title>{ data ? 'Editar Tipo de Examen' : 'Agregar Tipo de Examen'}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <FormTestType 
          dataEntry={data}
          saveData={(e) => eventHandler(e)}
        />
      </Modal.Body>
    </Modal>
  );
};

export default ModalTestType;
