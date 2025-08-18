import axios from 'axios';

const fetchData = async () => {
  try {
    const response = await axios.get('http://tu-backend.com/api/endpoint');
    console.log(response.data);
  } catch (error) {
    console.error('Error al hacer la solicitud:', error);
  }
};

fetchData();