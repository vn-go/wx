import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 5,
  duration: '5s',
};

const BASE = 'http://localhost:8081';

export default function () {
  const single = {
    file: http.file('hello from k6', 'hello.txt', 'text/plain'),
  };
  const res1 = http.post(`${BASE}/upload`, single);
  check(res1, {
    'single status 200': (r) => r.status === 200,
    'single filename ok': (r) => r.json().file === 'hello.txt',
  });

  const multi = {
    files: [
      http.file('alpha', 'a.txt', 'text/plain'),
      http.file('beta', 'b.txt', 'text/plain'),
    ],
  };
  const res2 = http.post(`${BASE}/upload`, multi);
  check(res2, {
    'multi status 200': (r) => r.status === 200,
    'multi length 2': (r) => Array.isArray(r.json().files) && r.json().files.length === 2,
  });

  sleep(1);
}
