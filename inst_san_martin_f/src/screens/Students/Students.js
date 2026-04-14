import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col } from 'react-bootstrap';

import ModalStudent from './Modal/ModalStudent';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree } from '../../context/Global/actions/DegreeActions';
import { getStudents, addStudent, updateStudent } from '../../context/Global/actions/StudentActions';

import styles from './Students.module.css';

const degreeNameById = (degrees, hexId) => {
  const d = (degrees || []).find(x => (x.id || x.ID) === hexId);
  return d ? d.name : hexId;
};

const buildTableRows = (students, degrees) =>
  (students || []).map(s => ({
    idAlumno: s.id || '',
    nombre: s.name,
    email: s.email || '',
    telefono: s.phone || '',
    dni: s.dni,
    direccion: s.address || '',
    carreras: (s.degreeIds || []).map(did => degreeNameById(degrees, did)).join(', ')
  }));

const Students = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);

  const showError = (message, type) => {
    message
      ? setError(
        <Notification
          message={message}
          type={type}
          show={showError}
        />
      )
      : setError(null);
  };

  const buildNotification = (result) => {
    switch (result.code) {
      case 199:
        showError(result.message, 'warning');
        break;
      case 200:
      case 201:
        showError(result.message, 'success');
        break;
      case 400:
        showError(result.message, 'danger');
        break;
      default:
        break;
    }
  };

  useEffect(() => {
    const load = async () => {
      await getDegree(globalDispatch);
      await getStudents(globalDispatch);
    };
    load();
  }, [globalDispatch]);

  const saveEventHandler = async (e) => {
    const isNew = !e.id;
    const result = isNew
      ? await addStudent(globalDispatch, e)
      : await updateStudent(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined && result.code !== null
      ? Number(result.code)
      : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return;
    }
    await getStudents(globalDispatch);
    setShow(false);
    setDataRow(null);
  };

  const addStudentEvent = () => {
    setDataRow(null);
    setShow(true);
  };

  const closeStudentEvent = () => {
    setDataRow(null);
    setShow(false);
  };

  const tableEvents = (ev, row) => {
    if (ev === 'edit') {
      const full = (globalState.students || []).find(s => s.id === row.idAlumno);
      setDataRow(full || row);
      setShow(true);
    }
  };

  const tableData = useMemo(
    () => buildTableRows(globalState.students, globalState.degrees),
    [globalState.students, globalState.degrees]
  );

  return (
    <React.Fragment>
      <Container className={styles.container}>
        {error}
        <Row>
          <h1>Alumnos</h1>
        </Row>
        <hr />
        <br />
        <Row className="justify-content-center">
          <Col xs lg="4">
            <Button
              className="w-100"
              variant="primary"
              size="sm"
              onClick={() => addStudentEvent()}
            >
              Agregar
            </Button>
          </Col>
        </Row>
        <br />
        <Row>
          <Col xs lg="1" />
          <Col>
            {tableData.length > 0 ? (
              <Table
                key="students-table"
                tableEvents={(e, d) => tableEvents(e, d)}
                data={tableData}
                actions="e"
              />
            ) : (
              <p className="text-muted">No hay alumnos registrados.</p>
            )}
          </Col>
          <Col xs lg="1" />
        </Row>
      </Container>
      <ModalStudent
        show={show}
        handleClose={closeStudentEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow}
        degrees={globalState.degrees || []}
      />
    </React.Fragment>
  );
};

export default Students;
