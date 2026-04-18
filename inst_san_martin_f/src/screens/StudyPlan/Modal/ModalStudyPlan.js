import React from 'react';
import { Modal } from 'react-bootstrap';

import FormStudyPlan from '../Form/FormStudyPlan';

const ModalStudyPlan = ({ show, handleClose, saveEvent, data, changeActive }) => {
  const eventHandler = async (e) => saveEvent(e);

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{data && (data.id || data.ID) ? 'Editar plan de estudio' : 'Agregar plan de estudio'}</Modal.Title>
      </Modal.Header>
      <Modal.Body style={{ overflowX: 'hidden', maxWidth: '100%' }}>
        <FormStudyPlan dataEntry={data} saveData={eventHandler} changeActive={changeActive} />
      </Modal.Body>
    </Modal>
  );
};

export default ModalStudyPlan;
