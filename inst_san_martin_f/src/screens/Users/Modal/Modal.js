import React from 'react';
import { Modal } from 'react-bootstrap';

import FormUser from '../Form/FormUser';

const UsersModal = ({ show, handleClose, saveEvent, data, roles, shifts, changeActive }) => {
  const eventHandler = async (e) => saveEvent(e);

  const isEdit = data && typeof data === 'object' && data.id;

  return (
    <Modal show={show} onHide={handleClose} centered size="lg" scrollable>
      <Modal.Header closeButton>
        <Modal.Title>{isEdit ? 'Editar usuario' : 'Agregar usuario'}</Modal.Title>
      </Modal.Header>
      <Modal.Body
        style={{
          maxHeight: 'min(78vh, 720px)',
          overflowY: 'auto',
          overflowX: 'hidden',
          maxWidth: '100%'
        }}
      >
        <FormUser
          key={isEdit ? data.id : 'new'}
          dataEntry={data || null}
          roles={roles}
          shifts={shifts}
          saveData={eventHandler}
          changeActive={changeActive}
        />
      </Modal.Body>
      <Modal.Footer />
    </Modal>
  );
};

export default UsersModal;
