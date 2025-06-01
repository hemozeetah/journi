import { useEffect, useState } from "react";
import { useParams } from "react-router";
import PlaceDetail from "../components/PlaceDetail";
import axios from "axios";

export default function PlaceDetailPage() {
  const { id } = useParams();

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

  return (
    <>
      <PlaceDetail place={place} />
    </>
  );
}
