import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';

import ModalTeacher from './Modal/ModalTeacher';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';
import ConfirmChangeEstadoModal from '../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getDegree } from '../../context/Global/actions/DegreeActions';
import { getShift } from '../../context/Global/actions/ShiftActions';
import {
  getTeachers,
  addTeacher,
  updateTeacher,
  changeActiveTeacher,
  getTitulosHabilitantes,
  getModalidadesDocente
} from '../../context/Global/actions/TeacherActions';

import styles from './Teachers.module.css';
import listToolbar from '../common/ListToolbar.module.css';
import { etiquetaNivel } from '../../constant/nivelesAcademicos';

const degreeNameById = (degrees, hexId) => {
  const d = (degrees || []).find((x) => (x.id || x.ID) === hexId);
  return d ? d.name : hexId;
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

const carreraYModalidadResumen = (t, degrees, modalidades) => {
  const parts = (t.careers || []).map((c) => {
    const nombre = degreeNameById(degrees, c.degreeId);
    const mod = modalidadLabel(modalidades, c.modalidadId);
    return `${nombre}: ${mod}`;
  });
  return parts.length ? parts.join('; ') : '—';
};

const buildTableRows = (teachers, degrees, modalidades) =>
  (teachers || []).map((t) => ({
    idDocente: t.id || '',
    nombre: t.name,
    email: t.email || '',
    telefono: t.phone || '',
    dni: t.dni,
    direccion: t.address || '',
    enseniaEn: (t.enseniaEn || []).map((n) => etiquetaNivel(n)).join(', '),
    resumen: carreraYModalidadResumen(t, degrees, modalidades),
    estado: t.active !== false
  }));

function Teachers() {
  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState(null);
  const [error, setError] = useState(null);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [titulosHabilitantes, setTitulosHabilitantes] = useState([]);
  const [modalidades, setModalidades] = useState([]);
  const [shiftsData, setShiftsData] = useState([]);
  const [activeConfirm, setActiveConfirm] = useState(null);
  const [activeConfirmLoading, setActiveConfirmLoading] = useState(false);

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
        const shiftsLoaded = await getShift(globalDispatch);
        setShiftsData(Array.isArray(shiftsLoaded) ? shiftsLoaded : []);
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
        return result;
      }
      await getTeachers(globalDispatch);
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
    if (ev === 'check') {
      const full = (globalState.teachers || []).find((t) => t.id === row.idDocente);
      const cur = full ? full.active !== false : row.estado === true;
      setActiveConfirm({
        id: row.idDocente,
        name: row.nombre || '',
        fromActive: cur,
        toActive: !cur
      });
    }
  };

  const confirmTeacherActive = async () => {
    if (!activeConfirm) return;
    setActiveConfirmLoading(true);
    const result = await changeActiveTeacher(globalDispatch, {
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
    () =>
      buildTableRows(globalState.teachers, globalState.degrees, modalidades),
    [globalState.teachers, globalState.degrees, modalidades]
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
          <h1>Docentes</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
            <Button className="w-100" variant="primary" size="sm" onClick={() => addTeacherEvent()}>
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
                key="teachers-table"
                tableEvents={(e, d) => tableEvents(e, d)}
                data={filteredTableData}
                actions="ec"
              />
            ) : (
              <p className="text-muted">No hay docentes registrados.</p>
            )}
          </Col>
        </Row>
      </Container>
      <ModalTeacher
        show={show}
        handleClose={closeTeacherEvent}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow}
        degrees={globalState.degrees || []}
        shifts={(Array.isArray(shiftsData) && shiftsData.length ? shiftsData : globalState.shifts) || []}
        titulosHabilitantes={titulosHabilitantes}
        modalidades={modalidades}
        changeActive={async (payload) => {
          const result = await changeActiveTeacher(globalDispatch, payload);
          buildNotification(result);
          return result;
        }}
      />
      <ConfirmChangeEstadoModal
        show={!!activeConfirm}
        onHide={() => !activeConfirmLoading && setActiveConfirm(null)}
        kind="docente"
        itemName={activeConfirm ? activeConfirm.name : ''}
        fromActive={activeConfirm ? activeConfirm.fromActive : true}
        toActive={activeConfirm ? activeConfirm.toActive : true}
        onConfirm={confirmTeacherActive}
        confirming={activeConfirmLoading}
      />
    </React.Fragment>
  );
}

Teachers.displayName = 'Teachers';

export default Teachers;
