fetch('http://tu-backend.com/api/endpoint')
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json();
  })
  .then(data => {
    console.log(data);
    // Aquí puedes actualizar tu estado o UI con los datos recibidos
  })
  .catch(error => {
    console.error('Hubo un problema con la solicitud:', error);
  });