import './PlaceList.css';

export default function PlaceList({ places }) {

  return (
    <>
      <div className="place-list-scroll">
        {places.map((place, index) => (
          <div key={index} className="place-card">
            <div className="place-image">
              {/* TODO put image */}
              <div className="image-placeholder"></div>
            </div>
            <div className="place-info">
              <h3 className="place-name">{place.name}</h3>
              <p className="place-caption" title={place.caption}>
                {place.caption.length > 100
                  ? `${place.caption.substring(0, 100)}...`
                  : place.caption}
              </p>
            </div>
          </div>
        ))}
      </div>
    </>
  );
}
