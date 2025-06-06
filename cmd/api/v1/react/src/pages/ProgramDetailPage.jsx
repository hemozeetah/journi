import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router";
import ProgramDetail from "../components/ProgramDetail";

export default function ProgramDetailPage() {
  const { id } = useParams();

  const [program, setProgram] = useState(null);
  const [company, setCompany] = useState(null);
  const [places, setPlaces] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const programRes = await axios.get(`http://localhost:8080/v1/programs/${id}`);
        setProgram(programRes.data);

        const companyRes = await axios.get(`http://localhost:8080/v1/users/${programRes.data.companyID}`);
        setCompany(companyRes.data);

        const journeysRes = await axios.get(`http://localhost:8080/v1/journeys?program_id=${id}`);
        const placesData = await Promise.all(
          journeysRes.data.map(journey =>
            axios.get(`http://localhost:8080/v1/places/${journey.placeID}`)
              .then(res => {
                return {
                  ...res.data,
                  startDatetime: journey.startDatetime,
                  endDatetime: journey.endDatetime
                }
              })
              .catch(err => {
                console.log(err.response.data);
                return null; // or handle error differently
              })
          )
        );

        setPlaces(placesData.filter(place => place !== null));

      } catch (err) {
        console.log(err.response?.data || err.message);
      }
    };

    fetchData();
  }, []);

  if (!program || !company) {
    return <div>Program not found</div>;
  }

  return (
    <>
      <ProgramDetail
        program={program}
        company={company}
        places={places}
      />
    </>
  );
}
