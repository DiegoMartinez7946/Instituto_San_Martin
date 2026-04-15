import React, { useState, useEffect, useMemo } from 'react';
import { Row } from 'react-bootstrap';
import lodash from 'lodash';

import HeaderCol from './HeaderCol';
import BodyCol from './BodyCol';

import styles from './Table.module.css';


const DEFAULT_WIDE_KEYS = ['resumen', 'carreras', 'enseniaen'];

const Table = ({ data, tableEvents, actions, wideKeys = DEFAULT_WIDE_KEYS }) => {
  const wideSet = useMemo(
    () => new Set((wideKeys || []).map((k) => String(k).toLowerCase())),
    [wideKeys]
  );

  const [header, setHeader] = useState(null);
  const [body, setBody] = useState(null);
  const [dataTable, setDataTable] = useState([]);

  useEffect(() => {
    if (data) {
      setDataTable(data);
    }
  }, [data]);

  useEffect(() => {
    if (!lodash.isEmpty(dataTable)) {
      loadHeader();
      loaderBody();
    }
  }, [dataTable, wideSet]);

  const loadHeader = () => {
    const keys = Object.keys(dataTable[0]);
    const values = keys.map((colKey, i) => (
      <HeaderCol
        key={colKey}
        data={lodash.capitalize(colKey)}
        colNumber={i}
        wide={wideSet.has(colKey.toLowerCase())}
      />
    ));
    setHeader([
      values,
      <HeaderCol
        key="__acciones__"
        data="Acciones"
        colNumber={actions ? keys.length : 0}
        wide={false}
        actionsHeader
      />
    ]);
  };

  const loaderBody = () => {
    const keys = Object.keys(dataTable[0]);
    const values = dataTable.map((row, i) => (
      <Row key={i} className={styles.tableRow}>
        {keys.map((colKey, j) => (
          <BodyCol
            key={colKey}
            wide={wideSet.has(colKey.toLowerCase())}
            data={typeof row[colKey] === 'boolean' ? (row[colKey] === true ? 'Activo' : 'Inactivo') : row[colKey]}
            colNumber={j}
            actions={actions}
          />
        ))}
        <BodyCol
          key="__actions__"
          data=""
          colActions={(e) => eventsHandler(e, row)}
          actions={actions}
          colNumber={actions ? keys.length + 1 : 0}
        />
      </Row>
    ));

    setBody(values);
  };

  const eventsHandler = (e, d) => {
    tableEvents(e, d);
  };

  return (
    <div className={styles.Table}>
      <Row className={styles.tableRow}>
        {header}
      </Row>
      {body}
    </div>
  )
};

export default Table;
