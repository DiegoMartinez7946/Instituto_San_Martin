import React  from 'react';
import { Container, Row, Col } from 'react-bootstrap';

import styles from './Home.module.css';

const Home = () => {

  return (
    <Container className={styles.container}>
      <Row className="justify-content-center">
        <Col xs lg="4" className="justify-content-center">
          <img src="/ista.png" alt="fondo"/>
        </Col>
      </Row>
    </Container>
  )

};

export default Home;
