import React, { useState, useEffect } from 'react';
import { Form, Button, InputGroup } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faEye, faEyeSlash } from '@fortawesome/free-solid-svg-icons';

import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';

const FormChange = ({ email, saveData }) => {
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock('pwd-change', email);

  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);

  const [data, setData] = useState({
    email,
    currentPassword: '',
    newPassword: ''
  });

  useEffect(() => {
    setData((prev) => ({ ...prev, email }));
  }, [email]);

  const handleInputChange = (event) => {
    setData({
      ...data,
      [event.target.name]: event.target.value
    });
  };

  const sendData = async (e) => {
    e.preventDefault();
    const res = await saveData(data);
    if (isSaveSuccess(res)) {
      setData({
        email,
        currentPassword: '',
        newPassword: ''
      });
      setShowCurrentPassword(false);
      setShowNewPassword(false);
      armLockAfterSave();
    }
  };

  return (
    <Form onSubmit={sendData}>
      <FormEditLockBanner
        entityKey="pwd-change"
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
        onCancelUnlock={() => setUnlocked(false)}
        unlockVariant="warning"
      />

      <Form.Group className="mb-3" controlId="formBasicEmail">
        <Form.Control type="hidden" name="email" value={data.email} />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPasswordCurrent">
        <Form.Label>Password Actual</Form.Label>
        <InputGroup>
          <Form.Control
            autoFocus
            type={showCurrentPassword ? 'text' : 'password'}
            placeholder="Ingrese password actual"
            name="currentPassword"
            onChange={handleInputChange}
            value={data.currentPassword}
            readOnly={readOnly}
            autoComplete="current-password"
          />
          <Button
            variant="outline-secondary"
            type="button"
            title={showCurrentPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            aria-label={showCurrentPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            onClick={() => setShowCurrentPassword((v) => !v)}
            tabIndex={-1}
            disabled={readOnly}
          >
            <FontAwesomeIcon icon={showCurrentPassword ? faEyeSlash : faEye} />
          </Button>
        </InputGroup>
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPasswordNew">
        <Form.Label>Password Nueva</Form.Label>
        <InputGroup>
          <Form.Control
            type={showNewPassword ? 'text' : 'password'}
            placeholder="Ingrese nueva password"
            name="newPassword"
            onChange={handleInputChange}
            value={data.newPassword}
            readOnly={readOnly}
            autoComplete="new-password"
          />
          <Button
            variant="outline-secondary"
            type="button"
            title={showNewPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            aria-label={showNewPassword ? 'Ocultar contraseña' : 'Mostrar contraseña'}
            onClick={() => setShowNewPassword((v) => !v)}
            tabIndex={-1}
            disabled={readOnly}
          >
            <FontAwesomeIcon icon={showNewPassword ? faEyeSlash : faEye} />
          </Button>
        </InputGroup>
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        Actualizar
      </Button>
    </Form>
  );
};

export default FormChange;
