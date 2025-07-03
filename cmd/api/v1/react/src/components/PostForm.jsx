import axios from "axios";
import { useEffect, useState } from "react";

export default function PostForm({ post = null, cities, places, token, claims, setPosts, setShowModal }) {
  const [data, setData] = useState({
    caption: '',
    placeID: ''
  });
  const [images, setImages] = useState([]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setData({
      ...data,
      [name]: value
    });
  };

  const handleImageChange = (e) => {
    if (e.target.files) {
      setImages(Array.from(e.target.files));
    }
  };

  const removeImage = (index) => {
    const updatedImages = [...images];
    updatedImages.splice(index, 1);

    setImages(updatedImages);
  };

  const [cityPlaces, setCityPlaces] = useState([]);

  const [selectedCityID, setSelectedCityID] = useState('');

  useEffect(() => {
    if (post) {
      setData({
        caption: post.caption,
        placeID: post.placeID,
      });
      const cityID = places.find(place => place.id === post.placeID).cityID;
      setSelectedCityID(cityID);
      setCityPlaces(places.filter((place) => {
        return place.cityID === cityID;
      }));
    }
  }, [post]);

  const handleCityChange = (e) => {
    setSelectedCityID(e.target.value);
    setCityPlaces(places.filter((place) => {
      return place.cityID === e.target.value;
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append('data', JSON.stringify(data));
    images.forEach((image, _) => {
      formData.append(`images`, image);
    });

    const url = post
      ? `http://localhost:8080/v1/posts/${post.id}`
      : "http://localhost:8080/v1/posts";

    const method = post ? 'put' : 'post';

    axios[method](url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        if (post) {
          setPosts(posts => posts.map(p => {
            if (p.id !== post.id) {
              return p;
            }
            return {
              ...res.data,
              userName: claims.name,
              userProfile: claims.profileURL,
              placeName: cityPlaces.find(place => place.id === data.placeID).name
            };
          }));
        } else {
          setPosts(posts => [...posts, {
            ...res.data,
            userName: claims.name,
            userProfile: claims.profileURL,
            placeName: cityPlaces.find(place => place.id === data.placeID).name
          }]);
        }
        setShowModal(false);
      })
      .catch(err => {
        console.log(err.response.data)
      });
  };

  return (
    <>
      <h2>{post ? 'Edit Post' : 'Add Post'}</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Caption:</label>
          <textarea
            name="caption"
            value={data.caption}
            onChange={handleInputChange}
            rows={5}
            required
          />
        </div>
        <div className="form-group">
          <label>Place:</label>
          <div className="place-selection">
            <select
              value={selectedCityID}
              onChange={handleCityChange}
            >
              <option disabled defaultValue="" value=""> -- select city -- </option>
              {cities.map(city => (
                <option key={city.id} value={city.id}>{city.name}</option>
              ))}
            </select>
            <select
              name="placeID"
              value={data.placeID}
              onChange={handleInputChange}
            >
              <option disabled defaultValue="" value=""> -- select place -- </option>
              {cityPlaces.map(place => (
                <option key={place.id} value={place.id}>{place.name}</option>
              ))}
            </select>
          </div>
        </div>
        <div className="form-group">
          <label>Select Images:</label>
          <input
            type="file"
            name="images"
            onChange={handleImageChange}
            multiple
            accept="image/*"
            required={!post}
          />

          {/* All image previews */}
          {images.length > 0 && (
            <div className="image-grid">
              {images.map((image, index) => (
                <div key={index} className="image-preview">
                  <img
                    src={URL.createObjectURL(image)}
                    alt={`image ${index}`}
                  />
                  <button
                    type="button"
                    onClick={() => removeImage(index)}
                    className="remove-btn"
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
        <button type="submit" className="submit-button">
          {post ? 'Update Post' : 'Add Post'}
        </button>
      </form>
    </>
  );
}
