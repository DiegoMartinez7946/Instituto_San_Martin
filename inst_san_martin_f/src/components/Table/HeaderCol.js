import React from 'react';
import { Col } from 'react-bootstrap';
import { transformTitles } from '../../utils/transformTitles';

import styles from './Table.module.css';

const HeaderCol = ({colNumber, data}) => {
  return (
    <Col className={colNumber > 0 ? '' : styles.hiddenDiv}>
      <div className={styles.HeaderCol}>{transformTitles(data)}</div>
    </Col>
  ); 
}

export default HeaderCol;
