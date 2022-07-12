import React from 'react';

import styles from './ToggleButton.module.css';

const ToggleButton = ({ click }) => (
    <button className={styles.toggle_button} onClick={click}>
        <div className={styles.toggle_button__line}/>
        <div className={styles.toggle_button__line}/>
        <div className={styles.toggle_button__line}/>
    </button>
);

export default ToggleButton;
