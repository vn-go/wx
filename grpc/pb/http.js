// http-login.js
import { check } from 'k6';
import http from 'k6/http';

export const options = {
    vus: 100,
    duration: '30s',
    thresholds: {
        'http_req_duration': ['p(95)<100'],  // 95% request < 100ms
        'http_req_failed': ['rate<0.01'],    // <1% lỗi
    },
};

const BASE_URL = 'http://localhost:8080';
const PAYLOAD = JSON.stringify({
    Username: 'admin',
    Password: '123456',
});

const PARAMS = {
    headers: {
        'Content-Type': 'application/json',
    },
};

export default function () {
    const res = http.post(`${BASE_URL}/api/user-controller/login`, PAYLOAD, PARAMS);

    // Kiểm tra response
    check(res, {
        'status is 200': (r) => r.status === 200,

    });

    // Log response nếu cần debug
    // console.log('Response:', res.body);

    //sleep(1); // Nghỉ 1 giây giữa các request
}