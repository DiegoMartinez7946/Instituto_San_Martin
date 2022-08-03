import React, { useState } from 'react';

import styles from './Layout.module.css';
import NavigationBar from './NavigationBar/NavigationBar';
import SideBar from './SideBar/SideBar';
import Backdrop from './Backdrop/Backdrop';

const Layout = (props) => {
    const [showSideBar, setShowSideBar] = useState(false);

    const toggleClickHandler = () => {
        setShowSideBar(current => !current);
    };

    const backdropClickHandler = () => {
        setShowSideBar(false);
    }

    return (
        <div className={styles.layout}>
            <NavigationBar toggleClick={toggleClickHandler}/>
            <SideBar show={showSideBar} click={backdropClickHandler}/>
            { showSideBar ? <Backdrop backdropClick={backdropClickHandler}/> : null }
            <main className={styles.content}>
                {props.children}
            </main>
        </div>
    );
};

export default Layout;
