// Ejemplo en React
import React, { useEffect, useState } from 'react';

const App = () => {
  const [data, setData] = useState([]);

  useEffect(() => {
    fetch('http://localhost:5000/api/endpoint') // Cambia la URL según tu configuración
      .then(response => response.json())
      .then(data => setData(data))
      .catch(error => console.error('Error:', error));
  }, []);

  return (
    <div>
      <h1>Datos del Backend</h1>
      <ul>
        {data.map(item => (
          <li key={item.id}>{item.name}</li> // Ajusta según la estructura de tus datos
        ))}
      </ul>
    </div>
  );
};

export default App;