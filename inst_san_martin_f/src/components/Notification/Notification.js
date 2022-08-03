import React from 'react';
import { Toast, ToastContainer } from 'react-bootstrap';

const Notification = ({message, type, show}) => {
  return (
    <ToastContainer position={'top-end'} >
      <Toast
          className="d-inline-block m-1"
          bg={type}
          onClose={() => show(null)} delay={3000} autohide
        >
          <Toast.Header>
            <strong className="me-auto">{type === 'success' ? 'Ok' : 'Error'}</strong>
          </Toast.Header>
          <Toast.Body className={type === 'warning' ? 'text-black' : 'text-white' }>
            { message }
          </Toast.Body>
        </Toast>
    </ToastContainer>
  );
};

export default Notification;
