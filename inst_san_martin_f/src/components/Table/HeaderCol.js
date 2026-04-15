import React from 'react';
import { Col } from 'react-bootstrap';
import { transformTitles } from '../../utils/transformTitles';

import styles from './Table.module.css';

const HeaderCol = ({ colNumber, data, wide, actionsHeader }) => {
  return (
    <Col
      xs={actionsHeader ? 'auto' : undefined}
      className={[
        colNumber > 0 ? '' : styles.hiddenDiv,
        wide ? styles.colWide : '',
        actionsHeader ? styles.BodyColActions : ''
      ]
        .filter(Boolean)
        .join(' ')}
    >
      <div className={styles.HeaderCol}>{transformTitles(data)}</div>
    </Col>
  );
};

export default HeaderCol;
