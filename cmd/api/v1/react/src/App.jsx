import { useState } from 'react';
import { Route, BrowserRouter as Router, Routes } from 'react-router';
import './App.css';
import Navbar from './components/Navbar';
import Cities from './pages/Cities';
import Home from './pages/Home';
import Places from './pages/Places';
import Programs from './pages/Programs';

function App() {
  const [claims, setClaims] = useState(() => {
    const saved = localStorage.getItem('claims');
    return saved ? JSON.parse(saved) : null;
  });

  const [token, setToken] = useState(() => {
    return localStorage.getItem('token') || null;
  });


  return (
    <Router>
      <Navbar
        claims={claims}
        setClaims={setClaims}
        setToken={setToken}
      />

      <Routes>
        <Route
          path="/"
          element={
            <Home
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/cities"
          element={
            <Cities
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/places"
          element={
            <Places />
          } />
        <Route
          path="/programs"
          element={
            <Programs />
          } />
        <Route
          path="*"
          element={
            <Home />
          } />
      </Routes>
    </Router>
  );
}

export default App
