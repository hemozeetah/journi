import { useNavigate } from 'react-router';
import './PlaceList.css';

export default function PlaceList({ places }) {
  const navigate = useNavigate();

  const handlePlaceClick = (placeID) => {
    navigate(`/places/${placeID}`);
  };

  return (
    <>
      <div className="place-list-scroll">
        {places.map((place, index) => (
          <div
            key={index}
            className="place-card"
            onClick={() => handlePlaceClick(place.id)}
          >
            <div className="place-image">
              {/* TODO put image */}
              <div className="image-placeholder"></div>
            </div>
            <div className="place-info">
              <h3 className="place-name">{place.name}</h3>
              <pre className="place-caption" title={place.caption}>
                {place.caption.length > 100
                  ? `${place.caption.substring(0, 100)}...`
                  : place.caption}
              </pre>
            </div>
          </div>
        ))}
      </div>
    </>
  );
}
