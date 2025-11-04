// grpc.js (TỐI ƯU – 1000+ req/s)
import { check } from 'k6';
import grpc from 'k6/net/grpc';

const client = new grpc.Client();

// INIT CONTEXT: CHỈ LOAD Ở ĐÂY
client.load(['.'], 'dynamic.proto');

export const options = {
    vus: 100,
    duration: '30s',
};

export default function () {
    // MỖI VU: CONNECT + INVOKE + CLOSE
    client.connect('localhost:50051', { plaintext: true });

    const request = {
        grpc_endpoint: 'main.user.Login',
        endcoder: 'gob',
        args: 'JX8DAQL/gAABAgEIVXNlcm5hbWUBDAABCFBhc3N3b3JkAQwAAAAS/4ABBWFkbWluAQYxMjM0NTYA'
    };

    const response = client.invoke('dynamic.DynamicService/Invoke', request);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
        'no error': (r) => r && r.message.error === '',
        'result received': (r) => r && r.message.result.length > 0,
    });

    client.close();
    // sleep(0.01); // nếu cần giới hạn
}