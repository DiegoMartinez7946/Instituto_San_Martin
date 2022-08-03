import React from 'react';
import { Col } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPencil, faTrash, faKey, faMagnifyingGlass, faCircleCheck } from '@fortawesome/free-solid-svg-icons';

import styles from './Table.module.css';

const BodyCol = ({data, colActions, actions, colNumber}) => {

  const eventHandler = (e) => {
    colActions(e);
  };

  return (
    <Col className={colNumber > 0 ? '' : styles.hiddenDiv}>
      {
        data ?
        <div className={styles.BodyCol}>{data}</div> :
        <div className={styles.BodyCol}>
          { actions.indexOf('e') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('edit')} icon={faPencil}/> : null } 
          { actions.indexOf('d') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('delete')} icon={faTrash}/> : null }        
          { actions.indexOf('p') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('password')} icon={faKey}/> : null }        
          { actions.indexOf('s') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('view')} icon={faMagnifyingGlass}/> : null }        
          { actions.indexOf('c') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('check')} icon={faCircleCheck}/> : null }        
        </div>
      }
    </Col>
  )
};

export default BodyCol;
