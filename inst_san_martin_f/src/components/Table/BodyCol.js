import React from 'react';
import { Col } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPencil, faTrash } from '@fortawesome/free-solid-svg-icons';

import styles from './Table.module.css';

const BodyCol = ({data, actions}) => {

  return (
    <Col>
      {
        data ?
        <div className={styles.BodyCol}>{data}</div> :
        <div className={styles.BodyCol}>
          <FontAwesomeIcon onClick={() => actions('edit')} icon={faPencil}/> 
          <FontAwesomeIcon onClick={() => actions('delete')} icon={faTrash}/> 
        </div>
      }
    </Col>
  )
};

export default BodyCol;
