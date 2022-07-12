import React from 'react';
import { Container, Row } from 'react-bootstrap';

import styles from './Main.module.css';

const Main = () => {
  return (
    <Container className={styles.container}>
      <Row>
        <h1>Main</h1>
      </Row>
    </Container>
  );
};

export default Main;
