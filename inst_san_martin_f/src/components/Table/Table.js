import React, { useState, useEffect } from 'react';
import { Row } from 'react-bootstrap';
import lodash from 'lodash';

import HeaderCol from './HeaderCol';
import BodyCol from './BodyCol';

import styles from './Table.module.css';


const Table = ({data}) => {

  const [header, setHeader] = useState(null);
  const [body, setBody] = useState(null);
 
  const loadHeader = () => {
    const values = Object.keys(
      data[0]).map((d, i) => {
        return <HeaderCol 
          key={i} 
          data={lodash.capitalize(d)}
        />
      });
    setHeader([values, <HeaderCol 
          key={data.length + 1} 
          data={'Acciones'}
        />]);  
  };

  const actionsEvent = (e) => { 
    console.log('acciones', e);
  }

  const loaderBody = () => {
    const values = data.map((d, i) => {
      return <Row key={i}>
        {Object.values(d).map((x, j) => {
          return <BodyCol key={j} data={x}/>
        })}
        <BodyCol 
          key={Object.values(data[0]).length + 1} 
          actions={(e) => actionsEvent(e)}
          data={''} 
        />
      </Row>
    });

    setBody(values);
  };

  useEffect(() => {
    loadHeader();
    loaderBody();
  }, []);


  return (
    <div className={styles.Table}>
      <Row>
        {header}
      </Row>
      {body}
    </div>
  )
};

export default Table;
