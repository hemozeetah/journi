import { useState } from 'react';
import { Route, BrowserRouter as Router, Routes } from 'react-router';
import './App.css';
import Navbar from './components/Navbar';
import CityListPage from './pages/CityListPage';
import Home from './pages/Home';
import CityDetailPage from './pages/CityDetailPage';
import PlaceDetailPage from './pages/PlaceDetailPage';
import ProgramListPage from './pages/ProgramListPage';

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
            <CityListPage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/cities/:id"
          element={
            <CityDetailPage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/places/:id"
          element={
            <PlaceDetailPage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/programs"
          element={
            <ProgramListPage
              claims={claims}
              token={token}
            />
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
