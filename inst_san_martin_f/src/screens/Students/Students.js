import React from 'react';
import { Container, Row } from 'react-bootstrap';

import styles from './Students.module.css';

const Students = () => {
  return (
    <Container className={styles.container}>
      <Row>
        <h1>Alumnos</h1>
      </Row>
    </Container>
  );
};

export default Students;
