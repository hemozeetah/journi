import { useState } from 'react';
import { Route, BrowserRouter as Router, Routes } from 'react-router';
import './App.css';
import Navbar from './components/Navbar';
import CityListPage from './pages/CityListPage';
import CityDetailPage from './pages/CityDetailPage';
import PlaceDetailPage from './pages/PlaceDetailPage';
import ProgramListPage from './pages/ProgramListPage';
import ProgramDetailPage from './pages/ProgramDetailPage';
import PostListPage from './pages/PostListPage';
import ProfilePage from './pages/ProfilePage';
import SubscriberDetailPage from './pages/SubscriberDetailPage';
import HomePage from './pages/HomePage';

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
            <HomePage
            />
          } />
        <Route
          path="/posts"
          element={
            <PostListPage
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
          path="/programs/:id"
          element={
            <ProgramDetailPage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/users/:id"
          element={
            <ProfilePage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="/subscribers/:id"
          element={
            <SubscriberDetailPage
              claims={claims}
              token={token}
            />
          } />
        <Route
          path="*"
          element={
            <PostListPage
              claims={claims}
              token={token}
            />
          } />
      </Routes>
    </Router>
  );
}

export default App
