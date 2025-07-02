import axios from "axios";
import { useState } from "react";

export default function ProgramForm({ cities, places, claims, token, setPrograms, setShowModal }) {
  const [data, setData] = useState({
    caption: '',
    startDate: '',
    endDate: ''
  });
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setData({
      ...data,
      [name]: value
    });
  };

  const [selectedPlaces, setSelectedPlaces] = useState([]);
  const [cityPlaces, setCityPlaces] = useState([]);

  const [selectedCityID, setSelectedCityID] = useState('');
  const [selectedCity, setSelectedCity] = useState('');
  const [selectedPlaceID, setSelectedPlaceID] = useState('');
  const [selectedPlace, setSelectedPlace] = useState('');
  const [placeStartDate, setPlaceStartDate] = useState('');
  const [placeEndDate, setPlaceEndDate] = useState('');

  const handleCityChange = (e) => {
    setSelectedCityID(e.target.value);
    setSelectedCity(cities.find(city => {
      return city.id === e.target.value;
    }));
    setCityPlaces(places.filter((place) => {
      return place.cityID === e.target.value;
    }));
  };

  const handlePlaceChange = (e) => {
    setSelectedPlaceID(e.target.value);
    setSelectedPlace(places.find(place => {
      return place.id === e.target.value;
    }));
  };

  const handleAddPlace = () => {
    if (!selectedPlaceID || !placeStartDate || !placeEndDate) return;

    const newPlace = {
      placeID: selectedPlaceID,
      placeName: selectedPlace.name,
      cityName: selectedCity.name,
      startDatetime: placeStartDate,
      endDatetime: placeEndDate
    };

    setSelectedPlaces([...selectedPlaces, newPlace]);

    // Reset the place selection fields
    setSelectedPlaceID('');
    setPlaceStartDate('');
    setPlaceEndDate('');
  };

  const handleRemovePlace = (index) => {
    setSelectedPlaces(selectedPlaces.filter((_, idx) => {
      return idx !== index;
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = {
      ...data,
      startDate: new Date(data.startDate).toISOString(),
      endDate: new Date(data.startDate).toISOString()
    };

    axios.post("http://localhost:8080/v1/programs", formData, {
      headers: {
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
        setPrograms(programs => [...programs, {
          ...res.data,
          companyName: claims.name,
        }]);
        const journeyPromises = selectedPlaces.map(place => {
          const data = {
            programID: res.data.id,
            placeID: place.placeID,
            startDatetime: new Date(place.startDatetime).toISOString(),
            endDatetime: new Date(place.endDatetime).toISOString(),
          };
          console.log(data);

          return axios.post("http://localhost:8080/v1/journeys", data, {
            headers: {
              'Authorization': "Bearer " + token
            }
          });
        });

        // Wait for all journeys to be created
        return Promise.all(journeyPromises);
      })
      .then(journeyResponses => {
        journeyResponses.forEach(res => console.log(res.data));
        setShowModal(false);
      })
      .catch(err => {
        if (err.response) {
          console.error("Server responded with error:", err.response.data);
          console.error("Status code:", err.response.status);
        } else if (err.request) {
          console.error("No response received:", err.request);
        } else {
          console.error("Request error:", err.message);
        }
      });
  };

  return (
    <>
      <h2>Add Program</h2>
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
          <label>Start Date:</label>
          <input
            type="date"
            name="startDate"
            value={data.startDate}
            onChange={handleInputChange}
            required
          />
        </div>
        <div className="form-group">
          <label>End Date:</label>
          <input
            type="date"
            name="endDate"
            value={data.endDate}
            onChange={handleInputChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Add Places:</label>
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
              value={selectedPlaceID}
              onChange={handlePlaceChange}
            >
              <option disabled defaultValue="" value=""> -- select place -- </option>
              {cityPlaces.map(place => (
                <option key={place.id} value={place.id}>{place.name}</option>
              ))}
            </select>

            <input
              type="datetime-local"
              value={placeStartDate}
              onChange={(e) => setPlaceStartDate(e.target.value)}
              placeholder="Start date"
            />

            <input
              type="datetime-local"
              value={placeEndDate}
              onChange={(e) => setPlaceEndDate(e.target.value)}
              placeholder="End date"
            />

            <button
              type="button"
              onClick={handleAddPlace}
              className="add-place-button"
            >
              Add Place
            </button>
          </div>

          {/* Display selected places */}
          <div className="selected-places">
            <h4>Selected Places:</h4>
            {selectedPlaces.length === 0 ? (
              <p>No places added yet</p>
            ) : (
              <ul>
                {selectedPlaces.map((place, index) => (
                  <li key={index}>
                    {place.cityName} - {place.placeName} ({place.startDatetime} to {place.endDatetime})
                    <button
                      type="button"
                      onClick={() => handleRemovePlace(index)}
                      className="remove-place-button"
                    >
                      Remove
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        <button type="submit" className="submit-button">
          Add Program
        </button>
      </form>
    </>
  );
}
