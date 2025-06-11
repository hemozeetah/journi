import axios from "axios";
import { useEffect, useState } from "react";
import FloatingModal from "../components/FloatingModal";
import PostForm from "../components/PostForm";
import PostList from "../components/PostList";

export default function HomePage({ token, claims }) {
  const [showModal, setShowModal] = useState(false);

  const [posts, setPosts] = useState([]);
  const [cities, setCities] = useState([]);
  const [places, setPlaces] = useState([]);

  useEffect(() => {
    axios.get("http://localhost:8080/v1/posts")
      .then(res => {
        const postsData = res.data;
        const enrichedPostPromises = postsData.map(async post => {
          const userPromise = axios.get(`http://localhost:8080/v1/users/${post.userID}`)
            .then(userRes => ({
              userName: userRes.data.name,
              userProfile: userRes.data.profile
            }))
            .catch(err => {
              console.log(`Error fetching user for post ${post.id}:`, err.response?.data);
              return {
                userName: 'Unknown',
                userProfile: ''
              };
            });
          const placePromise = post.placeID
            ? axios.get(`http://localhost:8080/v1/places/${post.placeID}`)
              .then(placeRes => ({
                placeName: placeRes.data.name
              }))
              .catch(err => {
                console.log(`Error fetching place for post ${post.id}:`, err.response?.data);
                return {
                  placeName: 'Unknown Place'
                };
              })
            : Promise.resolve({ placeName: '' });
          const [userData, placeData] = await Promise.all([userPromise, placePromise]);
          return ({
            ...post,
            ...userData,
            ...placeData
          });
        });
        Promise.all(enrichedPostPromises)
          .then(postsWithDetails => {
            setPosts(postsWithDetails);
            console.log(postsWithDetails);
          });
      })
      .catch(err => {
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
  }, []);

  return (
    <>
      <div className="post-list-container">
        <h1>Post List Page</h1>
        <div className="post-list">
          <div className="post-card add-post-card" onClick={() => setShowModal(true)}>
            <div className="add-post-content">
              <div className="plus-sign">+</div>
            </div>
          </div>
          <PostList posts={posts}/>
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
