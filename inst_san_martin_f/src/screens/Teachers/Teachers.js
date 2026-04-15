import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col } from 'react-bootstrap';

import ModalTeacher from './Modal/ModalTeacher';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree } from '../../context/Global/actions/DegreeActions';
import {
  getTeachers,
  addTeacher,
  updateTeacher,
  getTitulosHabilitantes,
  getModalidadesDocente
} from '../../context/Global/actions/TeacherActions';

import styles from './Teachers.module.css';
import { etiquetaNivel } from '../../constant/nivelesAcademicos';

const degreeNameById = (degrees, hexId) => {
  const d = (degrees || []).find((x) => (x.id || x.ID) === hexId);
  return d ? d.name : hexId;
};

const tituloLabel = (titulos, id) => {
  const t = (titulos || []).find((x) => (x.id || x.ID) === id);
  if (!t) return id || '—';
  return String(t.codigo).toUpperCase() === 'SI' ? 'Sí' : 'No';
};

const modalidadLabel = (mods, id) => {
  const m = (mods || []).find((x) => (x.id || x.ID) === id);
  if (!m) return id || '—';
  const c = String(m.codigo || '').toLowerCase();
  if (c === 'titular') return 'Titular';
  if (c === 'provisional') return 'Provisional';
  if (c === 'suplente') return 'Suplente';
  return m.codigo;
};

const buildTableRows = (teachers, degrees, titulos, modalidades) =>
  (teachers || []).map((t) => ({
    idDocente: t.id || '',
    nombre: t.name,
    email: t.email || '',
    telefono: t.phone || '',
    dni: t.dni,
    direccion: t.address || '',
    enseniaEn: (t.enseniaEn || []).map((n) => etiquetaNivel(n)).join(', '),
    carreras: (t.degreeIds || []).map((did) => degreeNameById(degrees, did)).join(', '),
    tituloHabilitante: tituloLabel(titulos, t.tituloHabilitanteId),
    modalidad: modalidadLabel(modalidades, t.modalidadId)
  }));

function Teachers() {
  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);
  const [titulosHabilitantes, setTitulosHabilitantes] = useState([]);
  const [modalidades, setModalidades] = useState([]);

  const showError = (message, type) => {
    message
      ? setError(
          <Notification message={message} type={type} show={showError} />
        )
      : setError(null);
  };

  const buildNotification = (result) => {
    if (!result) return;
    const code = Number(result.code);
    switch (code) {
      case 199:
        showError(result.message, 'warning');
        break;
      case 200:
      case 201:
        showError(result.message, 'success');
        break;
      case 400:
      case 403:
        showError(result.message || 'No autorizado', 'danger');
        break;
      default:
        break;
    }
  };

  useEffect(() => {
    const load = async () => {
      try {
        await getDegree(globalDispatch);
        const [tit, mod] = await Promise.all([
          getTitulosHabilitantes(),
          getModalidadesDocente()
        ]);
        setTitulosHabilitantes(Array.isArray(tit) ? tit : []);
        setModalidades(Array.isArray(mod) ? mod : []);
        await getTeachers(globalDispatch);
      } catch (e) {
        const msg =
          (e.response && e.response.data && e.response.data.message) ||
          e.message ||
          'Error al cargar docentes';
        showError(msg, 'danger');
      }
    };
    load();
  }, [globalDispatch]);

  const saveEventHandler = async (e) => {
    const isNew = !e.id;
    try {
      const result = isNew
        ? await addTeacher(globalDispatch, e)
        : await updateTeacher(globalDispatch, e);
      buildNotification(result);
      const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
      if (codeNum !== 200 && codeNum !== 201) {
        return;
      }
      await getTeachers(globalDispatch);
      setShow(false);
      setDataRow(null);
    } catch (err) {
      const data = err.response && err.response.data;
      buildNotification(
        data && data.message
          ? { ...data, code: data.code || err.response.status }
          : { message: err.message || 'Error de red', code: 400 }
      );
    }
  };

  const addTeacherEvent = () => {
    setDataRow(null);
    setShow(true);
  };

  const closeTeacherEvent = () => {
    setDataRow(null);
    setShow(false);
  };

  const tableEvents = (ev, row) => {
    if (ev === 'edit') {
      const full = (globalState.teachers || []).find((t) => t.id === row.idDocente);
      setDataRow(full || row);
      setShow(true);
    }
  };

  const tableData = useMemo(
    () =>
      buildTableRows(
        globalState.teachers,
        globalState.degrees,
        titulosHabilitantes,
        modalidades
      ),
    [globalState.teachers, globalState.degrees, titulosHabilitantes, modalidades]
  );

  return (
    <React.Fragment>
      <Container className={styles.container}>
        {error}
        <Row>
          <h1>Docentes</h1>
        </Row>
        <hr />
        <br />
        <Row className="justify-content-center">
          <Col xs lg="4">
            <Button className="w-100" variant="primary" size="sm" onClick={() => addTeacherEvent()}>
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
                key="teachers-table"
                tableEvents={(e, d) => tableEvents(e, d)}
                data={tableData}
                actions="e"
              />
            ) : (
              <p className="text-muted">No hay docentes registrados.</p>
            )}
          </Col>
          <Col xs lg="1" />
        </Row>
      </Container>
      <ModalTeacher
        show={show}
        handleClose={closeTeacherEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow}
        degrees={globalState.degrees || []}
        titulosHabilitantes={titulosHabilitantes}
        modalidades={modalidades}
      />
    </React.Fragment>
  );
}

Teachers.displayName = 'Teachers';

export default Teachers;
