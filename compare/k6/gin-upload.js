import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,
  duration: '5s',
};

const BASE = 'http://localhost:8081';
const newLocal = `./gin-upload.js`;
const file1 = open(newLocal, 'b');

export default function () {
  const single = {
    file: http.file(file1, 'HttpServer.go', 'text/plain'),
  };
  const res1 = http.post(`${BASE}/upload`, single);
  check(res1, {
    'single status 200': (r) => r.status === 200,
    'single filename ok': (r) => r.json().file === 'HttpServer.go',
  });

  sleep(1);
}
