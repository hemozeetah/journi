import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router";
import "./PostList.css";

function Post({ post, claims, token }) {
  const navigate = useNavigate();
  const [currentImageIndex, setCurrentImageIndex] = useState(0);
  const [showMenu, setShowMenu] = useState(false);
  const menuRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setShowMenu(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleUserClick = (userID) => {
    navigate(`/users/${userID}`);
  };

  const handlePlaceClick = (placeID) => {
    navigate(`/places/${placeID}`);
  };

  const toggleMenu = () => {
    setShowMenu(!showMenu);
  };

  const handleEdit = () => {
    console.log("Edit post", post.id);
    setShowMenu(false);
  };

  const handleDelete = () => {
    console.log("Delete post", post.id);
    setShowMenu(false);
  };

  return (
    <div key={post.id} className="post-card">
      <div className="post-header">
        <div className="user-container">
          <img
            src={post.userProfile ? "http://localhost:8080" + post.userProfile : '/profile.png'}
            alt="profile"
            className="user-profile"
            onClick={() => handleUserClick(post.userID)}
          />
          <h2
            className="user-name"
            onClick={() => handleUserClick(post.userID)}
          >
            {post.userName}
          </h2>
        </div>
        {claims && claims.id === post.userID && (
          <div className="post-actions" ref={menuRef}>
            <button className="post-menu-button" onClick={toggleMenu}>
              <span>⋮</span>
            </button>
            {showMenu && (
              <div className="post-menu">
                <button onClick={handleEdit}>Edit</button>
                <button onClick={handleDelete}>Delete</button>
              </div>
            )}
          </div>
        )}
      </div>
      <h3
        className="place-name"
        onClick={() => handlePlaceClick(post.placeID)}
      >
        <i>at {post.placeName}</i>
      </h3>
      <pre className="post-caption">{post.caption}</pre>
      <div className="post-images">
        <div className="image-slider">
          <img
            src={"http://localhost:8080" + post.imagesURL[currentImageIndex]}
            alt={post.userName}
          />
        </div>
        <div className="image-indicators">
          {post.imagesURL.map((_, index) => (
            <div
              key={index}
              className={`indicator ${index === currentImageIndex ? 'active' : ''}`}
              onClick={() => setCurrentImageIndex(index)}
            />
          ))}
        </div>
      </div>
    </div>
  );
}

export default function PostList({ posts, claims, token }) {
  return (
    <>
      {posts.map((post) => (
        <Post
          key={post.id}
          post={post}
          claims={claims}
          token={token}
        />
      ))}
    </>
  );
}
