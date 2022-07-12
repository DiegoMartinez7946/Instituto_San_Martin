import React from 'react';
import { Routes, Route } from 'react-router-dom';

import Protected from './utils/Protected';
import Logout from './utils/Logout';

import Login from './screens/Login/Login';
import Main from './screens/Main/Main';
import Roles from './screens/Roles/Roles';
import Students from './screens/Students/Students';
import Users from './screens/Users/Users';
import Home from './screens/Home/Home';
import Layout from './components/Layout/Layout';

const App = () => {


  const routes = (
    <Routes>
      <Route exact path="/" element={<Home />} />
      <Route exact path="/login" element={<Login/>} />
      <Route exact path="/logout" element={<Logout />} />
      <Route exact path="/main" element={<Protected><Main/></Protected>} />
      <Route exact path="/roles" element={<Protected><Roles/></Protected>} />
      <Route exact path="/users" element={<Protected><Users/></Protected>} />
      <Route exact path="/students" element={<Protected><Students/></Protected>} />
    </Routes>
  );
  return (
      <Layout>
        {routes} 
      </Layout>
  );
}

export default App;
