import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    const payload = JSON.stringify({ x: 12, y: 33, z: 2222 });
    const headers = { 'Content-Type': 'application/json' };
    const response = http.post('http://localhost:8080', payload, { headers });
    sleep(1);
}