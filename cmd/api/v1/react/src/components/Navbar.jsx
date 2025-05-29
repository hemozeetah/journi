import axios from 'axios';
import React, { useState } from 'react';
import { Link } from 'react-router';
import FloatingModal from './FloatingModal';
import './Navbar.css';

export default function Navbar({ claims, setClaims, setToken }) {
  const [showModal, setShowModal] = useState(false);
  const [isSignUp, setIsSignUp] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const toggleSignUp = () => {
    setIsSignUp(!isSignUp);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    if (isSignUp) {
      if (formData.password !== formData.confirmPassword) {
        alert("Passwords don't match!");
        return;
      }
      axios.post("http://localhost:8080/v1/users", {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        passwordConfirm: formData.confirmPassword
      }).then(_ => {
        // TODO flash message
        toggleSignUp()
      }).catch(err => {
        console.log(err)
      })
    } else {
      axios.get("http://localhost:8080/v1/auth/token", {
        auth: {
          username: formData.email,
          password: formData.password,
        }
      }).then(res => {
        setClaims(res.data.claims);
        setToken(res.data.token);
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('claims', JSON.stringify(res.data.claims));
        setShowModal(false);
      }).catch(err => {
        console.log(err)
      })
    }
  };

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
        <Link to="/cities" className="nav-link">Cities</Link>
        <Link to="/places" className="nav-link">Places</Link>
        <Link to="/programs" className="nav-link">Programs</Link>
      </div>

      <div className="navbar-center">
        <input
          type="text"
          placeholder="Search..."
          className="search-bar"
        />
      </div>

      <div className="navbar-right">
        {claims && (
          <button className="auth-button" onClick={handleLogout}>
            logout
          </button>
        )}
        {!claims && (
          <button className="auth-button" onClick={() => setShowModal(true)}>
            Sign In/Sign Up
          </button>
        )}
      </div>

      {showModal && (
        <FloatingModal setShowModal={setShowModal}>
          <h2>{isSignUp ? 'Sign Up' : 'Sign In'}</h2>
          <form onSubmit={handleSubmit}>
            {isSignUp && (
              <div className="form-group">
                <label>Name</label>
                <input
                  type="text"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  required
                />
              </div>
            )}
            <div className="form-group">
              <label>Email</label>
              <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleInputChange}
                required
              />
            </div>
            <div className="form-group">
              <label>Password</label>
              <input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                required
              />
            </div>
            {isSignUp && (
              <div className="form-group">
                <label>Confirm Password</label>
                <input
                  type="password"
                  name="confirmPassword"
                  value={formData.confirmPassword}
                  onChange={handleInputChange}
                  required
                />
              </div>
            )}
            <button type="submit" className="submit-button">
              {isSignUp ? 'Sign Up' : 'Sign In'}
            </button>
          </form>
          <p className="toggle-auth">
            {isSignUp ? 'Already have an account? ' : "Don't have an account? "}
            <button onClick={toggleSignUp} className="toggle-button">
              {isSignUp ? 'Sign In' : 'Sign Up'}
            </button>
          </p>
        </FloatingModal>
      )}
    </nav>
  );
};
