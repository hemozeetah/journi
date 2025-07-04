import axios from "axios";
import { useEffect, useState } from "react";
import "./ProgramForm.css";

export default function ProgramForm({ cities, places, claims, token, program = null, setProgram, currentPlaces, setPrograms, setShowModal }) {
  const [data, setData] = useState({
    caption: '',
    startDate: '',
    endDate: ''
  });
  const [errors, setErrors] = useState({
    dateRange: '',
    placeDateRange: '',
    existingPlaces: ''
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setData({
      ...data,
      [name]: value
    });
    // Clear errors when changing dates
    if (name === 'startDate' || name === 'endDate') {
      setErrors({...errors, dateRange: '', existingPlaces: ''});
      if (selectedPlaces.length > 0 && newData.startDate && newData.endDate) {
        validateExistingPlaces(newData.startDate, newData.endDate);
      }
    }
  };

  const [selectedPlaces, setSelectedPlaces] = useState([]);
  const [cityPlaces, setCityPlaces] = useState([]);

  const [selectedCityID, setSelectedCityID] = useState('');
  const [selectedCity, setSelectedCity] = useState('');
  const [selectedPlaceID, setSelectedPlaceID] = useState('');
  const [selectedPlace, setSelectedPlace] = useState('');
  const [placeStartDate, setPlaceStartDate] = useState('');
  const [placeEndDate, setPlaceEndDate] = useState('');

  useEffect(() => {
    if (program) {
      setData({
        caption: program.caption,
        startDate: new Date(program.startDate).toISOString().split('T')[0],
        endDate: new Date(program.endDate).toISOString().split('T')[0]
      });
      setSelectedPlaces(currentPlaces.map(place => ({
        ...place,
        cityName: cities.find(city => city.id === place.cityID).name
      })));
    }
  }, [program]);

  const validateProgramDates = () => {
    if (data.startDate && data.endDate) {
      const start = new Date(data.startDate);
      const end = new Date(data.endDate);

      if (start > end) {
        setErrors({ ...errors, dateRange: 'End date must be after start date' });
        return false;
      }
    }
    return true;
  };

  const validateExistingPlaces = (newStartDate, newEndDate) => {
    const programStart = new Date(newStartDate);
    const programEnd = new Date(newEndDate);
    const programEndPlusOne = new Date(programEnd);
    programEndPlusOne.setDate(programEndPlusOne.getDate() + 1);


    const invalidPlaces = selectedPlaces.filter(place => {
      const placeStart = new Date(place.startDatetime);
      const placeEnd = new Date(place.endDatetime);
      return placeStart < programStart || placeEnd > programEndPlusOne;
    });

    if (invalidPlaces.length > 0) {
      setErrors({
        ...errors, existingPlaces:
          `${invalidPlaces.length} place(s) are outside the new program date range. Please adjust or remove them.`
      });
      return false;
    }

    return true;
  };


  const validatePlaceDates = () => {
    if (!data.startDate || !data.endDate) {
      setErrors({ ...errors, placeDateRange: 'Please set program dates first' });
      return false;
    }

    if (!placeStartDate || !placeEndDate) {
      setErrors({ ...errors, placeDateRange: 'Both place dates are required' });
      return false;
    }

    const programStart = new Date(data.startDate);
    const programEnd = new Date(data.endDate);
    const programEndPlusOne = new Date(programEnd);
    programEndPlusOne.setDate(programEndPlusOne.getDate() + 1);

    const placeStart = new Date(placeStartDate);
    const placeEnd = new Date(placeEndDate);

    if (placeStart >= placeEnd) {
      setErrors({ ...errors, placeDateRange: 'Place end date must be after start date' });
      return false;
    }

    if (placeStart < programStart || placeEnd > programEndPlusOne) {
      setErrors({ ...errors, placeDateRange: 'Place dates must be within program date range' });
      return false;
    }

    return true;
  };

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
    if (!selectedPlaceID) {
      setErrors({ ...errors, placeDateRange: 'Please select a place' });
      return;
    }
    if (!validatePlaceDates()) return;

    const newPlace = {
      id: selectedPlaceID,
      name: selectedPlace.name,
      startDatetime: placeStartDate,
      endDatetime: placeEndDate,
      cityName: selectedCity.name
    };

    setSelectedPlaces([...selectedPlaces, newPlace]);
    setErrors({ ...errors, placeDateRange: '' });

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
    if (!validateProgramDates()) return;
    if (!validateExistingPlaces(data.startDate, data.endDate)) {
      return;
    }


    const formData = {
      ...data,
      startDate: new Date(data.startDate).toISOString(),
      endDate: new Date(data.endDate).toISOString()
    };

    const url = program
      ? `http://localhost:8080/v1/programs/${program.id}`
      : "http://localhost:8080/v1/programs";

    const method = program ? 'put' : 'post';

    console.log(method, url);

    axios[method](url, formData, {
      headers: {
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
        if (program) {
          setProgram(res.data);

          // First, get all existing journeys for this program
          return axios.get(`http://localhost:8080/v1/journeys?program_id=${program.id}`, {
            headers: {
              'Authorization': "Bearer " + token
            }
          }).then(journeysResponse => {
            // Then delete each journey
            const deletePromises = journeysResponse.data.map(journey => {
              return axios.delete(`http://localhost:8080/v1/journeys/${journey.id}`, {
                headers: {
                  'Authorization': "Bearer " + token
                }
              });
            });
            return Promise.all(deletePromises);
          }).then(() => res); // Return the original program response
        } else {
          setPrograms(programs => [...programs, {
            ...res.data,
            companyName: claims.name,
          }]);
          return res; // Return the response for the new program
        }
      })
      .then(res => {
        // Now create new journeys with the selected places
        const journeyPromises = selectedPlaces.map(place => {
          const data = {
            programID: res.data.id,
            placeID: place.id,
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

        return Promise.all(journeyPromises);
      })
      .then(journeyResponses => {
        journeyResponses.forEach(res => console.log(res.data));
        setShowModal(false);
        if (program) {
          // TODO update instead of refresh page
          window.location.reload();
        }
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
      <h2>{program ? 'Edit Program' : 'Add Program'}</h2>
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
          {errors.dateRange && <p className="error">{errors.dateRange}</p>}
          {errors.existingPlaces && <p className="error">{errors.existingPlaces}</p>}
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
              onChange={(e) => {
                setPlaceStartDate(e.target.value);
                setErrors({ ...errors, placeDateRange: '' });
              }}
              placeholder="Start date"
              style={{ marginTop: 10 }}
            />

            <input
              type="datetime-local"
              value={placeEndDate}
              onChange={(e) => {
                setPlaceEndDate(e.target.value);
                setErrors({ ...errors, placeDateRange: '' });
              }}
              placeholder="End date"
              style={{ marginTop: 10, marginBottom: 10 }}
            />

            <button
              type="button"
              onClick={handleAddPlace}
              className="add-place-button"
            >
              Add Place
            </button>
            {errors.placeDateRange && <p className="error">{errors.placeDateRange}</p>}
          </div>

          {/* Display selected places */}
          <div className="selected-places">
            <h4>Selected Places:</h4>
            {selectedPlaces.length === 0 ? (
              <p>No places added yet</p>
            ) : (
              <ul>
                {selectedPlaces.map((place, index) => (
                  <li key={index} style={{ padding: 4 }}>
                    {place.cityName} - {place.name} ({new Date(place.startDatetime).toLocaleString()} to {new Date(place.endDatetime).toLocaleString()})
                    <button
                      type="button"
                      onClick={() => handleRemovePlace(index)}
                      className="remove-place-button"
                      style={{ marginLeft: 8 }}
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
          {program ? 'Update Program' : 'Add Program'}
        </button>
      </form>
    </>
  );
}
