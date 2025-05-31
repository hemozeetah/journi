import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import CityDetail from "../components/CityDetail";
import FloatingModal from "../components/FloatingModal";
import PlaceForm from "../components/PlaceForm";
import PlaceList from "../components/PlaceList";

export default function CityDetailPage({ claims, token }) {
  const isAdmin = claims && claims.role === "admin";
  const [showModal, setShowModal] = useState(false);

  const { id } = useParams();

  const [city, setCity] = useState(null);
  const [places, setPlaces] = useState([]);

  useEffect(() => {
    axios.get("http://localhost:8080/v1/cities/" + id)
      .then(res => {
        setCity(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
    axios.get("http://localhost:8080/v1/places?city_id=" + encodeURIComponent(id))
      .then(res => {
        setPlaces(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
  }, []);

  if (!city) {
    return <div>City not found</div>;
  }

  return (
    <>
      <CityDetail city={city} />
      <div className="place-list-container">
        <h2 className="place-list-title">{city.name}'s Places</h2>
        <PlaceList places={places} claims={claims} token={token} />
        {isAdmin && (
          <>
            <div className="place-list-container place-card add-place-card" onClick={() => setShowModal(true)}>
              <div className="add-place-content">
                <div className="plus-sign">+</div>
              </div>
            </div>
            {showModal && (
              <FloatingModal setShowModal={setShowModal}>
                <PlaceForm
                  token={token}
                  city={city}
                  setPlaces={setPlaces}
                  setShowModal={setShowModal}
                />
              </FloatingModal>
            )}
          </>
        )}
      </div>
    </>
  );
}
