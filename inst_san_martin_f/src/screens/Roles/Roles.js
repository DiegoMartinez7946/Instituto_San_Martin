import React, { useState, useEffect } from 'react';
import { Button, Container, Row, Col } from 'react-bootstrap';
import ModalRole from './Modal/ModalRole';
import Table from '../../components/Table/Table';

import { useGlobal } from '../../context/Global/GlobalProvider';
import { getRole } from '../../context/Global/actions/RoleActions';

import styles from './Roles.module.css';

const Roles = () => {

  const [show, setShow] = useState(false);
  const [globalState, globalDispatch] = useGlobal();

  useEffect(() => {
    const getAllRoles = async () => {
      await getRole(globalDispatch);
      return; 
    };
    getAllRoles();
  }, []);

  const addRoleEvent = () => {
    setShow(current => !current);
  }

  return (
    <React.Fragment>
      <Container className={styles.container}>
        <Row>
          <h1>Roles</h1>
        </Row>
        <Row className="justify-content-center">
          <Col xs lg="4">
              <Button 
                className="w-100"
                variant="primary"
                size="xs"
                onClick={() => addRoleEvent()}
              >
                Agregar
              </Button>
          </Col>
        </Row>
        <br />  
        <Row>
          <Col xs lg="2"></Col>
          <Col>
            { globalState.roles.length ? 
              <Table
                key={'role'}
                data={globalState.roles}
              /> :
              <p>No hay datos disponibles</p>
            }
          </Col>
          <Col xs lg="2"></Col>
        </Row>
      </Container>
      <ModalRole show={show}/>
    </React.Fragment>
  );
};

export default Roles;
