import React, { useState, useEffect, useMemo } from 'react';
import { Button, Container, Row, Col, Form, InputGroup } from 'react-bootstrap';
import ModalStudyPlan from './Modal/ModalStudyPlan';
import Table from '../../components/Table/Table';
import Notification from '../../components/Notification/Notification';
import ConfirmChangeEstadoModal from '../../components/ConfirmChangeEstadoModal/ConfirmChangeEstadoModal';

import { useGlobal } from '../../context/Global/GlobalProvider';
import {
  getStudyPlans,
  addStudyPlan,
  updateStudyPlan,
  changeActiveStudyPlan
} from '../../context/Global/actions/StudyPlanActions';

import styles from './StudyPlan.module.css';
import listToolbar from '../common/ListToolbar.module.css';

const StudyPlan = () => {
  const [globalState, globalDispatch] = useGlobal();
  const [show, setShow] = useState(false);
  const [dataRow, setDataRow] = useState('');
  const [rows, setRows] = useState([]);
  const [searchText, setSearchText] = useState('');
  const [appliedSearch, setAppliedSearch] = useState('');
  const [error, setError] = useState(null);
  const [estadoConfirm, setEstadoConfirm] = useState(null);
  const [estadoConfirmLoading, setEstadoConfirmLoading] = useState(false);

  const showError = (message, type) => {
    message
      ? setError(
          <Notification message={message} type={type} show={showError} />
        )
      : setError(null);
  };

  const refreshPlans = () => getStudyPlans(globalDispatch);

  const buildNotification = (result) => {
    const c = result && result.code !== undefined && result.code !== null ? Number(result.code) : NaN;
    switch (c) {
      case 199:
        showError(result.message, 'warning');
        break;
      case 200:
      case 201:
        showError(result.message, 'success');
        break;
      case 400:
      case 422:
        showError(result.message, 'danger');
        break;
      default:
        if (result && result.message) showError(result.message, 'danger');
        break;
    }
  };

  useEffect(() => {
    refreshPlans().catch(() => {});
  }, []);

  useEffect(() => {
    const list = Array.isArray(globalState.studyPlans) ? globalState.studyPlans : [];
    setRows(
      list.map((p) => ({
        id: p.id || p.ID,
        nombre: p.nombre,
        numeroResolucion: p.numeroResolucion,
        cohorte: p.cohorte,
        extension: p.extension,
        anioAprobacion: p.anioAprobacion,
        anioCaducidad: p.anioCaducidad,
        cargaHoraria: p.cargaHoraria,
        active: p.active
      }))
    );
  }, [globalState.studyPlans]);

  const saveEventHandler = async (e) => {
    const result = !(e.id || e.ID)
      ? await addStudyPlan(globalDispatch, e)
      : await updateStudyPlan(globalDispatch, e);

    buildNotification(result);
    const codeNum = result && result.code !== undefined ? Number(result.code) : NaN;
    if (codeNum !== 200 && codeNum !== 201) {
      return result;
    }
    await refreshPlans();
    setShow(false);
    setDataRow('');
    return result;
  };

  const tableEvents = async (ev, d) => {
    if (ev === 'edit') {
      const full = (globalState.studyPlans || []).find((p) => (p.id || p.ID) === (d.id || d.ID)) || d;
      setDataRow(full);
      setShow(true);
    }
    if (ev === 'check') {
      const cur = d.active === true;
      setEstadoConfirm({ row: d, fromActive: cur, toActive: !cur });
    }
  };

  const confirmEstadoToggle = async () => {
    if (!estadoConfirm) return;
    setEstadoConfirmLoading(true);
    const item = { ...estadoConfirm.row, active: estadoConfirm.toActive };
    const result = await changeActiveStudyPlan(globalDispatch, item);
    buildNotification(result);
    setEstadoConfirmLoading(false);
    if (Number(result.code) === 200) {
      setEstadoConfirm(null);
      await refreshPlans();
    }
  };

  const filtered = useMemo(() => {
    const q = String(appliedSearch || '').trim().toLowerCase();
    if (!q) return rows;
    return rows.filter((r) => String(r.numeroResolucion || '').toLowerCase().includes(q));
  }, [rows, appliedSearch]);

  return (
    <React.Fragment>
      <Container fluid className={styles.container}>
        {error}
        <Row>
          <h1>Planes de estudio</h1>
        </Row>
        <hr />
        <br />
        <div className={listToolbar.toolbarRow}>
          <div className={listToolbar.toolbarHalf}>
            <Button
              className="w-100"
              variant="primary"
              size="xs"
              onClick={() => {
                setDataRow('');
                setShow(true);
              }}
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
                placeholder="Buscar por número de resolución"
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
            <Table
              key="studyplan"
              tableEvents={(e, d) => tableEvents(e, d)}
              data={filtered}
              actions="ec"
              wideKeys={[]}
            />
          </Col>
        </Row>
      </Container>
      <ModalStudyPlan
        show={show}
        handleClose={() => {
          setDataRow('');
          setShow(false);
        }}
        saveEvent={(e) => saveEventHandler(e)}
        data={dataRow || ''}
        changeActive={async (payload) => {
          const result = await changeActiveStudyPlan(globalDispatch, payload);
          buildNotification(result);
          return result;
        }}
      />
      <ConfirmChangeEstadoModal
        show={!!estadoConfirm}
        onHide={() => !estadoConfirmLoading && setEstadoConfirm(null)}
        kind="planestudio"
        itemName={
          estadoConfirm && estadoConfirm.row
            ? estadoConfirm.row.nombre || estadoConfirm.row.numeroResolucion || ''
            : ''
        }
        fromActive={estadoConfirm ? estadoConfirm.fromActive : true}
        toActive={estadoConfirm ? estadoConfirm.toActive : true}
        onConfirm={confirmEstadoToggle}
        confirming={estadoConfirmLoading}
      />
    </React.Fragment>
  );
};

export default StudyPlan;
