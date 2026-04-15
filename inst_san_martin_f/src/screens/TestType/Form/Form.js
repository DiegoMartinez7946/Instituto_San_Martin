import React, { useState, useEffect } from 'react';
import { Form, Button } from 'react-bootstrap';
import lodash from 'lodash';

import FormEditLockBanner from '../../../components/FormEditLockBanner/FormEditLockBanner';
import { useFormEditLock } from '../../../hooks/useFormEditLock';
import { isSaveSuccess } from '../../../utils/saveResult';

const FormTestType = ({ dataEntry, saveData }) => {
  const [data, setData] = useState({
    ID: '',
    type: ''
  });

  const entry = dataEntry && typeof dataEntry === 'object' ? dataEntry : {};
  const lockEntityKey = entry.ID || '';
  const { readOnly, unlocked, setUnlocked, armLockAfterSave } = useFormEditLock(lockEntityKey, dataEntry);

  useEffect(() => {
    const e = dataEntry && typeof dataEntry === 'object' ? dataEntry : {};
    setData({
      ID: e.ID || '',
      type: e.type || ''
    });
  }, [dataEntry]);

  const handleInputChange = (event) => {
    setData({
      ...data,
      [event.target.name]: lodash.toUpper(event.target.value)
    });
  };

  const sendData = async (e) => {
    e.preventDefault();
    const res = await saveData(data);
    if (isSaveSuccess(res)) {
      armLockAfterSave();
    }
  };

  const isEdit = !!(entry && entry.ID);

  return (
    <Form onSubmit={sendData}>
      <Form.Group className="mb-3" controlId="formBasicID">
        <Form.Control type="hidden" disabled name="ID" value={data ? data.ID : ''} />
      </Form.Group>

      <FormEditLockBanner
        entityKey={lockEntityKey}
        unlocked={unlocked}
        onUnlock={() => setUnlocked(true)}
      />

      <Form.Group className="mb-3" controlId="formBasicType">
        <Form.Label>Tipo de Examen</Form.Label>
        <Form.Control
          type="text"
          placeholder="Ingrese tipo de examen"
          name="type"
          onChange={handleInputChange}
          value={data ? data.type : ''}
          autoFocus
          readOnly={readOnly}
        />
      </Form.Group>
      <br />
      <Button variant="primary" type="submit" className="w-100" disabled={readOnly}>
        {isEdit ? 'Actualizar' : 'Guardar'}
      </Button>
    </Form>
  );
};

export default FormTestType;
