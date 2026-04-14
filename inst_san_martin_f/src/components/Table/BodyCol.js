import React from 'react';
import { Col } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPencil, faTrash, faKey, faMagnifyingGlass, faCircleCheck } from '@fortawesome/free-solid-svg-icons';

import styles from './Table.module.css';

const BodyCol = ({data, colActions, actions, colNumber}) => {

  const eventHandler = (e) => {
    if (colActions) {
      colActions(e);
    }
  };

  const isActionsColumn = typeof colActions === 'function';

  return (
    <Col className={colNumber > 0 ? '' : styles.hiddenDiv}>
      {
        isActionsColumn ?
        <div className={styles.BodyCol}>
          { actions && actions.indexOf('e') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('edit')} icon={faPencil}/> : null } 
          { actions && actions.indexOf('d') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('delete')} icon={faTrash}/> : null }        
          { actions && actions.indexOf('p') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('password')} icon={faKey}/> : null }        
          { actions && actions.indexOf('s') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('view')} icon={faMagnifyingGlass}/> : null }        
          { actions && actions.indexOf('c') > -1 ? <FontAwesomeIcon onClick={() => eventHandler('check')} icon={faCircleCheck}/> : null }        
        </div> :
        <div className={styles.BodyCol}>
          {data === undefined || data === null ? '' : String(data)}
        </div>
      }
    </Col>
  )
};

export default BodyCol;
