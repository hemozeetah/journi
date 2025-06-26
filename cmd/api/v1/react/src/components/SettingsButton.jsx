import React, { useEffect, useRef, useState } from 'react';
import './SettingsButton.css';

export default function SettingsButton({ onEdit, onDelete }) {
  const [isOpen, setIsOpen] = useState(false);
  const menuRef = useRef(null);

  const toggleMenu = () => {
    setIsOpen(!isOpen);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="settings-container" ref={menuRef}>
      <button className="settings-button" onClick={toggleMenu}>
        <span className="settings-dot"></span>
        <span className="settings-dot"></span>
        <span className="settings-dot"></span>
      </button>

      {isOpen && (
        <div className="settings-menu">
          <button
            className="settings-menu-item"
            onClick={() => {
              onEdit();
              setIsOpen(false);
            }}
          >
            Edit
          </button>
          <button
            className="settings-menu-item delete"
            onClick={() => {
              onDelete();
              setIsOpen(false);
            }}
          >
            Delete
          </button>
        </div>
      )}
    </div>
  );
};
