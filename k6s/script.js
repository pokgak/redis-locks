import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 10,
  duration: '30s',
};

export default function() {
  const url = 'http://localhost:8080/order';
  const payload = JSON.stringify({ key1: 'value1', key2: 'value2' });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}
