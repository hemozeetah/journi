import { Link } from "react-router";

export default function HomePage() {
  return (
    <>
      <div className="home-page">
        {/* Features Section */}
        <section className="features">
          <div className="container">
            <h2>Your Gateway to Syrian Adventures</h2>
            <div className="features-grid">
              <div className="feature-card">
                <div className="feature-icon">🏙️</div>
                <h3>Explore Cities</h3>
                <p>Discover the rich history and culture of Syria's most beautiful cities</p>
                <Link to="/cities" className="feature-link">View Cities →</Link>
              </div>

              <div className="feature-card">
                <div className="feature-icon">✈️</div>
                <h3>Join Journeys</h3>
                <p>Embark on curated travel experiences with expert guides</p>
                <Link to="/programs" className="feature-link">Browse Journeys →</Link>
              </div>

              <div className="feature-card">
                <div className="feature-icon">📝</div>
                <h3>Share Stories</h3>
                <p>Document your travels and inspire others with your experiences</p>
                <Link to="/posts" className="feature-link">See Posts →</Link>
              </div>
            </div>
          </div>
        </section>

        {/* For Companies Section */}
        <section className="for-companies">
          <div className="container">
            <div className="company-content">
              <h2>For Travel Companies</h2>
              <p>Manage and promote your Syrian travel experiences with our professional tools</p>
              <ul className="company-features">
                <li>Create and manage travel packages</li>
                <li>Connect with adventure seekers</li>
                <li>Track bookings and payments</li>
                <li>Get detailed analytics</li>
              </ul>
            </div>
          </div>
        </section>
      </div>
    </>
  );
}
