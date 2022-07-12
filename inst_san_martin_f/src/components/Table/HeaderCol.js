import React from 'react';
import { Col } from 'react-bootstrap';

import styles from './Table.module.css';

const HeaderCol = ({data}) => {
  return (
    <Col>
      <div className={styles.HeaderCol}>{data}</div>
    </Col>
  ); 
}

export default HeaderCol;
