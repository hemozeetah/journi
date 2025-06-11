import axios from "axios";
import { useState } from "react";

export default function PostForm({ cities, places, token, claims, setPosts, setShowModal }) {
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

  const [cityPlaces, setCityPlaces] = useState([]);

  const [selectedCityID, setSelectedCityID] = useState('');

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
    axios.post("http://localhost:8080/v1/posts", formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        // TODO flash message
        setPosts(posts => [...posts, {
          ...res.data,
          userName: claims.name,
          userProfile: claims.profile,
          placeName: cityPlaces.find(place => place.id === data.placeID).name
        }]);
        setShowModal(false);
      })
      .catch(err => {
        console.log(err.response.data)
      });
  };

  return (
    <>
      <h2>Add Post</h2>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Caption:</label>
          <input
            type="text"
            name="caption"
            value={data.caption}
            onChange={handleInputChange}
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
            required
          />
        </div>
        <button type="submit" className="submit-button">
          Add Post
        </button>
      </form>
    </>
  );
}
