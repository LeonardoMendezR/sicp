fetch('http://tu-backend.com/api/endpoint')
  .then(response => response.json())
  .then(data => {
    console.log(data);
    // Actualiza tu estado o UI con los datos recibidos
  })
  .catch(error => {
    console.error('Error:', error);
  });