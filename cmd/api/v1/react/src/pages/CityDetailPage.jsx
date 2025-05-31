import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import CityDetail from "../components/CityDetail";

export default function CityDetailPage({ claims, token }) {
  const { id } = useParams();

  const [city, setCity] = useState(null);
  useEffect(() => {
    axios.get("http://localhost:8080/v1/cities/" + id)
      .then(res => {
        setCity(res.data);
        console.log(res.data);
      }).catch(err => {
        console.log(err.response.data)
      })
  }, []);

  if (!city) {
    return <div>City not found</div>;
  }

  return (
    <>
      <CityDetail city={city} />
    </>
  );
}
