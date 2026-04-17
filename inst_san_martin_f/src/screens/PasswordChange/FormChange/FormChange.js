import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';

import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';

const FormChange = ({ email, saveData }) => {
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock('pwd-change', email);

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
        <Form.Control
          autoFocus
          type="text"
          placeholder="Ingrese password actual"
          name="currentPassword"
          onChange={handleInputChange}
          value={data.currentPassword}
          readOnly={readOnly}
        />
      </Form.Group>

      <Form.Group className="mb-3" controlId="formBasicPasswordNew">
        <Form.Label>Password Nueva</Form.Label>
        <Form.Control
          type="text"
          placeholder="Ingrese nueva password"
          name="newPassword"
          onChange={handleInputChange}
          value={data.newPassword}
          readOnly={readOnly}
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        Actualizar
      </Button>
    </Form>
  );
};

export default FormChange;
