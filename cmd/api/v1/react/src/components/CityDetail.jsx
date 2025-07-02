import { useEffect, useState } from "react";
import "./CityDetail.css"

export default function CityDetail({ city }) {
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  useEffect(() => {
    if (!city.imagesURL || city.imagesURL.length <= 1) return;

    const interval = setInterval(() => {
      setCurrentImageIndex(prev =>
        prev === city.imagesURL.length - 1 ? 0 : prev + 1
      );
    }, 5000);

    return () => clearInterval(interval);
  }, [city]);

  return (
    <>
      <div className="city-detail">
        <div className="city-info">
          <h1>{city.name}</h1>
          <pre className="city-description">{city.caption}</pre>
        </div>

        <div className="city-images">
          <div className="image-slider">
            <img
              src={"http://localhost:8080" + city.imagesURL[currentImageIndex]}
              alt={city.name}
            />
          </div>
          <div className="image-indicators">
            {city.imagesURL.map((_, index) => (
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
