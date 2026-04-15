import React from 'react';
import { Container, Row } from 'react-bootstrap';

import { useGlobal } from '../../context/Global/GlobalProvider';

import styles from './Main.module.css';

const Main = () => {
  const [globalState] = useGlobal();
  const t = globalState.portalTime;

  return (
    <Container className={styles.container}>
      <Row>
        <h1 className={styles.title}>home</h1>
      </Row>
      {t ? (
        <Row className={styles.clockRow}>
          <div className={styles.clockCard}>
            <p className={styles.clockCard__label}>Hora en Argentina</p>
            <p className={styles.clockCard__hora}>{t.hora}</p>
            <dl className={styles.clockCard__grid}>
              <dt>Fecha</dt>
              <dd>{t.fecha}</dd>
              <dt>Día</dt>
              <dd>{t.dia}</dd>
              <dt>Mes</dt>
              <dd>{t.mes}</dd>
              <dt>Año</dt>
              <dd>{t.año}</dd>
            </dl>
          </div>
        </Row>
      ) : null}
    </Container>
  );
};

export default Main;
