import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import FloatingModal from "../components/FloatingModal";
import PostForm from "../components/PostForm";
import PostList from "../components/PostList";
import UserDetail from "../components/UserDetail";

export default function ProfilePage({ token, claims }) {
  const [showModal, setShowModal] = useState(false);

  const { id } = useParams();

  const [user, setUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [cities, setCities] = useState([]);
  const [places, setPlaces] = useState([]);

  useEffect(() => {
    axios.get(`http://localhost:8080/v1/users/${id}`)
      .then(res => {
        const userData = res.data;
        setUser(userData);
        console.log(res.data);
        axios.get(`http://localhost:8080/v1/posts?user_id=${id}`)
          .then(res => {
            const postsData = res.data;
            const placePromise = postsData.map(post =>
              axios.get(`http://localhost:8080/v1/places/${post.placeID}`)
                .then(res => ({
                  ...post,
                  userName: userData.name,
                  userProfile: userData.profileImageURL,
                  placeName: res.data.name
                }))
                .catch(err => {
                  console.log(`Error fetching place for post ${post.id}:`, err.response?.data);
                  return {
                    ...post,
                    userName: userData.name,
                    userProfile: userData.profileImageURL,
                    placeName: 'Unknown Place'
                  }
                }));
            Promise.all(placePromise)
              .then(postsWithPlaces => {
                setPosts(postsWithPlaces);
                console.log(postsWithPlaces);
              });
          })
          .catch(err => {
            console.log(err.response.data)
          });
      }).catch(err => {
        console.log(err.response.data)
      });
    axios.get("http://localhost:8080/v1/cities")
      .then(res => {
        setCities(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
    axios.get("http://localhost:8080/v1/places")
      .then(res => {
        setPlaces(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
  }, [id]);

  return (
    <>
      <UserDetail
        user={user}
        claims={claims}
        token={token}
        />
      <div className="post-list-container">
        <div className="post-list">
          {claims && claims.id === id && (
            <div className="post-card add-post-card" onClick={() => setShowModal(true)}>
              <div className="add-post-content">
                <div className="plus-sign">+</div>
              </div>
            </div>
          )}
          <PostList posts={posts} />
        </div>
      </div>
      {showModal && (
        <FloatingModal setShowModal={setShowModal}>
          <PostForm
            cities={cities}
            places={places}
            token={token}
            claims={claims}
            setPosts={setPosts}
            setShowModal={setShowModal}
          />
        </FloatingModal>
      )}
    </>
  );
}
