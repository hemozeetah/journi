import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router';
import "./UserSearch.css";

export default function UserSearch({ users }) {
  const navigate = useNavigate();

  const [searchTerm, setSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const searchRef = useRef(null);

  const handleSearch = (e) => {
    const term = e.target.value;
    setSearchTerm(term);

    if (term.trim() === '') {
      setSearchResults([]);
      setIsDropdownOpen(false);
      return;
    }

    const results = users.filter(user =>
      user.name.toLowerCase().includes(term.toLowerCase())
    );

    setSearchResults(results);
    setIsDropdownOpen(results.length > 0);
  };

  const handleClick = (userID) => {
    setIsDropdownOpen(false);
    navigate(`/users/${userID}`);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (searchRef.current && !searchRef.current.contains(event.target)) {
        setIsDropdownOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="navbar-center" ref={searchRef}>
      <input
        type="text"
        placeholder="Search users..."
        className="search-bar"
        value={searchTerm}
        onChange={handleSearch}
      />

      {isDropdownOpen && (
        <div className="search-dropdown">
          <div className="search-results">
            {searchResults.map(user => (
              <div
                key={user.id}
                className="user-result"
                onClick={() => handleClick(user.id)}
              >
                <div className={`user-result-role-${user.role}`}>
                  {user.name} {user.role !== 'user' && '*'}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
