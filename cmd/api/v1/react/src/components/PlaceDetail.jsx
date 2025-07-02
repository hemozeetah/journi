import { useEffect, useState } from "react";
import "./PlaceDetail.css"

export default function PlaceDetail({ place }) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  useEffect(() => {
    if (!place.imagesURL || place.imagesURL.length <= 1) return;

    const interval = setInterval(() => {
      setCurrentImageIndex(prev =>
        prev === place.imagesURL.length - 1 ? 0 : prev + 1
      );
    }, 5000);

    return () => clearInterval(interval);
  }, [place]);

  return (
    <>
      <div className="place-detail">
        <div className="place-info">
          <h1>{place.name}</h1>
          <h3 className="place-type">{place.type}</h3>
          <pre className="place-description">{place.caption}</pre>
        </div>

        <div className="place-images">
          <div className="image-slider">
            <img
              src={"http://localhost:8080" + place.imagesURL[currentImageIndex]}
              alt={place.name}
            />
          </div>
          <div className="image-indicators">
            {place.imagesURL.map((_, index) => (
              <div
                key={index}
                className={`indicator ${index === currentImageIndex ? 'active' : ''}`}
                onClick={() => setCurrentImageIndex(index)}
              />
            ))}
          </div>
        </div>
      </div>
    </>
  );
}
