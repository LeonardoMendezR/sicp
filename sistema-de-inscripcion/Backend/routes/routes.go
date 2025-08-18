fetch('http://localhost:5000/api/endpoint') // Cambia la URL según tu configuración
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
  .then(data => {
    console.log(data); // Maneja los datos recibidos
  })
  .catch(error => {
    console.error('There was a problem with the fetch operation:', error);
  });