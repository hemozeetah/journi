import { useState } from "react";
import { useNavigate } from "react-router";
import "./PostList.css";

function Post({ post }) {
  const navigate = useNavigate();

  const handleUserClick = (userID) => {
    navigate(`/users/${userID}`);
  };

  const handlePlaceClick = (placeID) => {
    navigate(`/places/${placeID}`);
  };

  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  return (
    <>
      <div
        key={post.id}
        className="post-card"
      >
        <div className="user-container">
          <img
            src={"http://localhost:8080" + post.userProfile}
            alt="User profile"
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
        <h3
          className="place-name"
          onClick={() => handlePlaceClick(post.placeID)}
        >
          <i>at</i> {post.placeName}
        </h3>
        <p className="post-caption">{post.caption}</p>
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

    </>
  );
}

export default function PostList({ posts }) {

  return (
    <>
      {posts.map((post) => (
        <Post post={post} />
      ))}
    </>
  );
}
