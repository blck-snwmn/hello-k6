import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
    // 注文登録
    const order = { id: Math.floor(Math.random() * 100000).toString(), item: 'example', quantity: 1 };
    const headers = { 'Content-Type': 'application/json' };
    const res = http.post('http://localhost:8080/orders', JSON.stringify(order), { headers });
    check(res, { 'status is 201': (r) => r.status === 201 });

    // 注文状態更新
    const update = { id: order.id, shipped: true };
    const res2 = http.put('http://localhost:8080/orders', JSON.stringify(update), { headers });
    check(res2, { 'status is 200': (r) => r.status === 200 });

    // 注文一覧取得
    const res3 = http.get('http://localhost:8080/orders');
    check(res3, { 'status is 200': (r) => r.status === 200 });

    sleep(1);
}

export const options = {
    vus: 100, // 並列ユーザー数
    duration: '30s', // テスト実行時間
};