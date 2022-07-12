import React from 'react';
import { Container, Row } from 'react-bootstrap';

import styles from './Users.module.css';

const Users = () => {
  return (
    <Container className={styles.container}>
      <Row>
        <h1>Users</h1>
      </Row>
    </Container>
  );
};

export default Users;
