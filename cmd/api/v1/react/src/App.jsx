import { useState } from 'react';
import './App.css';
import CreatePost from './components/CreatePost';
import Navbar from './components/Navbar';

function App() {
  const [claims, setClaims] = useState(() => {
    const saved = localStorage.getItem('claims');
    return saved ? JSON.parse(saved) : null;
  });

  const [token, setToken] = useState(() => {
    return localStorage.getItem('token') || null;
  });


  return (
    <>
      <Navbar
        claims={claims}
        setClaims={setClaims}
        setToken={setToken}
      />
      {token && <p>{token}</p>}
      {claims && (
        <div>
          <p>{claims.name}</p>
          <p>{claims.role}</p>
        </div>
      )}

      <CreatePost />
    </>
  )
}

export default App
