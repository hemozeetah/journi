import axios from 'axios';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router';
import FloatingModal from './FloatingModal';
import './Navbar.css';
import SignInForm from './SignInForm';
import SignUpForm from './SignUpForm';
import UserSearch from './UserSearch';

export default function Navbar({ claims, setClaims, setToken }) {
  const [showModal, setShowModal] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);

  const [users, setUsers] = useState([]);

  useEffect(() => {
    axios.get("http://localhost:8080/v1/users")
      .then(res => {
        setUsers(res.data);
      })
      .catch(err => {
        console.log(err.response.data);
      })
  }, []);

  const handleLogout = () => {
    setClaims(null);
    setToken(null);
    localStorage.removeItem('token');
    localStorage.removeItem('claims');
  }

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <Link to="/" className="nav-link">Home</Link>
        <Link to="/posts" className="nav-link">Posts</Link>
        <Link to="/cities" className="nav-link">Cities</Link>
        <Link to="/programs" className="nav-link">Programs</Link>
      </div>

      <UserSearch users={users} />

      <div className="navbar-right">
        {claims && (
          <>
            <Link to={`/users/${claims.id}`} className="nav-link">{claims.name}</Link>
            <button className="auth-button" onClick={handleLogout}>
              logout
            </button>
          </>
        )}
        {!claims && (
          <button className="auth-button" onClick={() => setShowModal(true)}>
            Sign In/Sign Up
          </button>
        )}
      </div>

      {showModal && (
        <FloatingModal setShowModal={setShowModal}>
          {isSignUp ? (
            <SignUpForm
              setIsSignUp={setIsSignUp}
            />
          ) : (
            <SignInForm
              setClaims={setClaims}
              setToken={setToken}
              setShowModal={setShowModal}
              setIsSignUp={setIsSignUp}
            />
          )
          }
        </FloatingModal>
      )}
    </nav>
  );
};
