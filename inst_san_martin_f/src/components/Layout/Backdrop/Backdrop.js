import React from 'react';

import styles from './Backdrop.module.css';

const Backdrop = ({ backdropClick }) => (
    <div className={styles.backdrop} onClick={backdropClick}/>
);

export default Backdrop;
