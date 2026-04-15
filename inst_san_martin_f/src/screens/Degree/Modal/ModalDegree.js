import React from 'react';
import { Modal } from 'react-bootstrap';

import FormDegree from '../Form/FormDegree';

const ModalDegree = ({ show, handleClose, saveEvent, data, changeActive }) => {
  const eventHandler = async (e) => saveEvent(e);

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{data && (data.id || data.ID) ? 'Editar Carrera' : 'Agregar Carrera'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormDegree dataEntry={data} saveData={eventHandler} changeActive={changeActive} />
      </Modal.Body>
    </Modal>
  );
};

export default ModalDegree;
