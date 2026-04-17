import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';

import ModalStudent from './Modal/ModalStudent';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';
import ConfirmChangeEstadoModal from '../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree } from '../../context/Global/actions/DegreeActions';
import {
  getStudents,
  addStudent,
  updateStudent,
  changeActiveStudent
} from '../../context/Global/actions/StudentActions';

import styles from './Students.module.css';
import listToolbar from '../common/ListToolbar.module.css';
import { etiquetaNivel } from '../../constant/nivelesAcademicos';

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
    nivelAprobado: s.nivelAprobado ? etiquetaNivel(s.nivelAprobado) : '—',
    carreras: (s.degreeIds || []).map(did => degreeNameById(degrees, did)).join(', '),
    estado: s.active !== false
  }));

const Students = () => {

  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [activeConfirm, setActiveConfirm] = useState(null);
  const [activeConfirmLoading, setActiveConfirmLoading] = useState(false);

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
    try {
      const isNew = !e.id;
      const result = isNew
        ? await addStudent(globalDispatch, e)
        : await updateStudent(globalDispatch, e);

      buildNotification(result);
      const codeNum =
        result && result.code !== undefined && result.code !== null
          ? Number(result.code)
          : NaN;
      if (codeNum !== 200 && codeNum !== 201) {
        return result;
      }
      await getStudents(globalDispatch);
      setShow(false);
      setDataRow(null);
      return result;
    } catch (err) {
      const data = err.response && err.response.data;
      const fail =
        data && data.message
          ? { ...data, code: data.code || err.response.status }
          : { message: err.message || 'Error de red', code: 400 };
      buildNotification(fail);
      return fail;
    }
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
      const list = globalState.students || [];
      const full =
        list.find((s) => String(s.id) === String(row.idAlumno)) ||
        (row.dni ? list.find((s) => String(s.dni) === String(row.dni)) : null);
      if (!full) {
        showError(
          'No se encontró el alumno en memoria. Recargue la página o use Agregar.',
          'warning'
        );
        return;
      }
      setDataRow(full);
      setShow(true);
    }
    if (ev === 'check') {
      const full = (globalState.students || []).find(
        (s) => String(s.id) === String(row.idAlumno)
      );
      const cur = full ? full.active !== false : row.estado === true;
      setActiveConfirm({
        id: row.idAlumno,
        name: row.nombre || '',
        fromActive: cur,
        toActive: !cur
      });
    }
  };

  const confirmStudentActive = async () => {
    if (!activeConfirm) return;
    setActiveConfirmLoading(true);
    const result = await changeActiveStudent(globalDispatch, {
      id: activeConfirm.id,
      active: activeConfirm.toActive
    });
    buildNotification(result);
    setActiveConfirmLoading(false);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum === 200) {
      setDataRow((prev) => {
        if (!prev || prev.id !== activeConfirm.id) return prev;
        return { ...prev, active: activeConfirm.toActive };
      });
      setActiveConfirm(null);
    } else if (codeNum === 199) {
      setActiveConfirm(null);
    }
  };

  const tableData = useMemo(
    () => buildTableRows(globalState.students, globalState.degrees),
    [globalState.students, globalState.degrees]
  );
  const filteredTableData = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return tableData;
    return tableData.filter((r) => String(r.nombre || '').toLowerCase().includes(q));
  }, [tableData, appliedSearch]);

  return (
    <React.Fragment>
      <Container fluid className={styles.container}>
        {error}
        <Row>
          <h1>Alumnos</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
            <Button
              className="w-100"
              variant="primary"
              size="sm"
              onClick={() => addStudentEvent()}
            >
              Agregar
            </Button>
          </div>
          <div className={listToolbar.toolbarHalf}>
            <InputGroup className="w-100">
              <Form.Control
                type="text"
                value={searchText}
                onChange={(e) => {
                  const v = e.target.value;
                  setSearchText(v);
                  if (!String(v).trim()) setAppliedSearch('');
                }}
                placeholder="Buscar por nombre"
              />
              <Button type="button" className={listToolbar.buscarBtn} onClick={() => setAppliedSearch(searchText)}>
                Buscar
              </Button>
            </InputGroup>
          </div>
        </div>
        <br />
        <Row className={styles.tableRowWrap}>
          <Col xs={12} className="px-2 px-md-3 min-w-0">
            {filteredTableData.length > 0 ? (
              <Table
                key="students-table"
                tableEvents={(e, d) => tableEvents(e, d)}
                data={filteredTableData}
                actions="ec"
              />
            ) : (
              <p className="text-muted">No hay alumnos registrados.</p>
            )}
          </Col>
        </Row>
      </Container>
      <ModalStudent
        show={show}
        handleClose={closeStudentEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow}
        degrees={globalState.degrees || []}
        changeActive={async (payload) => {
          const result = await changeActiveStudent(globalDispatch, payload);
          buildNotification(result);
          return result;
        }}
      />
      <ConfirmChangeEstadoModal
        show={!!activeConfirm}
        onHide={() => !activeConfirmLoading && setActiveConfirm(null)}
        kind="alumno"
        itemName={activeConfirm ? activeConfirm.name : ''}
        fromActive={activeConfirm ? activeConfirm.fromActive : true}
        toActive={activeConfirm ? activeConfirm.toActive : true}
        onConfirm={confirmStudentActive}
        confirming={activeConfirmLoading}
      />
    </React.Fragment>
  );
};

export default Students;
