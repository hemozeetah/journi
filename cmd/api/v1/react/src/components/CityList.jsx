import { useNavigate } from 'react-router';
import './CityList.css';

export default function CityList({ cities }) {
  const navigate = useNavigate();

  const handleCityClick = (cityID) => {
    navigate(`/cities/${cityID}`);
  };

  return (
    <>
      {cities.map((city) => (
        <div
          key={city.id}
          className="city-card"
          onClick={() => handleCityClick(city.id)}
        >
          <h3>{city.name}</h3>
          <pre className="city-caption">
            {city.caption.length > 100
              ? `${city.caption.substring(0, 100)}...`
              : city.caption}
          </pre>
        </div>
      ))}
    </>
  );
}
