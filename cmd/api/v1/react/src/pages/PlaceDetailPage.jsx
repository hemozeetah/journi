import axios from "axios";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import PlaceDetail from "../components/PlaceDetail";
import PlaceForm from "../components/PlaceForm";
import FloatingModal from "../components/FloatingModal";
import SettingsButton from "../components/SettingsButton";

export default function PlaceDetailPage({ claims, token }) {
  const { id } = useParams();
  const navigate = useNavigate();

  const isAdmin = claims && claims.role === "admin";
  const [showModalEdit, setShowModalEdit] = useState(false);

  const [place, setPlace] = useState(null);

  useEffect(() => {
    axios.get("http://localhost:8080/v1/places/" + id)
      .then(res => {
        setPlace(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      });
  }, []);

  if (!place) {
    return <div>Place not found</div>
  }

  const handleEdit = () => {
    setShowModalEdit(true);
  };

  const handleDelete = () => {
    axios.delete(`http://localhost:8080/v1/places/${id}`, {
      headers: {
        'Authorization': "Bearer " + token
      }
    })
      .then(res => {
        console.log(res.data);
        navigate(`/cities/${place.cityID}`);
      })
      .catch(err => {
        console.log(err.response.data);
      });
  };

  return (
    <>
      <PlaceDetail place={place} />
      {isAdmin && (
        <>
          {showModalEdit && (
            <FloatingModal setShowModal={setShowModalEdit}>
              <PlaceForm
                token={token}
                setShowModal={setShowModalEdit}
                place={place}
                setPlace={setPlace}
              />
            </FloatingModal>
          )}
          <SettingsButton
            onEdit={handleEdit}
            onDelete={handleDelete}
          />
        </>
      )}
    </>
  );
}
