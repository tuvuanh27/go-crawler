import axios from 'axios';

export const apiClient = axios.create({
  baseURL: 'http://localhost:5002/api/v1/',
});

// config interceptors for apiClient when response status >= 400
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    // if error.response.data blob convert to json
    if (error.response.data instanceof Blob) {
      const errorText = await error.response.data.text();

      return Promise.reject(JSON.parse(errorText));
    }

    return Promise.reject(error.response.data);
  },
);

export const API_METHODS = {
  GET: 'GET',
  POST: 'POST',
  PUT: 'PUT',
  PATCH: 'PATCH',
  DELETE: 'DELETE',
  OPTIONS: 'OPTIONS',
} as const;
