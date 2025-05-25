export default function Home({ token, claims }) {
  return (
    <>
      <h1>Home Page</h1>
      {token && <p>{token}</p>}
      {claims && (
        <div>
          <p>{claims.name}</p>
          <p>{claims.role}</p>
        </div>
      )}
    </>
  );
}
