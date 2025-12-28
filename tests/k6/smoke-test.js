import http from 'k6/http';
import { check, sleep } from 'k6';

// 簡單的煙霧測試 - 快速驗證系統基本功能
export const options = {
    vus: 5,              // 5個虛擬用戶
    duration: '30s',     // 運行30秒
    thresholds: {
        'http_req_duration': ['p(95)<200'],  // 95%請求<200ms
        'http_req_failed': ['rate<0.05'],     // 錯誤率<5%
    },
};

const MOCK_PLACES = [
    'mock_tokyo_ramen_1',
    'mock_osaka_sushi_1',
    'mock_kyoto_tempura_1',
];

export function setup() {
    const res = http.post(
        'http://localhost:8080/api/v1/auth/register',
        JSON.stringify({
            email: `smoke${Date.now()}@test.com`,
            username: 'smoketest',
            password: 'Test1234!',
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    return { token: res.json('access_token') };
}

export default function (data) {
    const placeId = MOCK_PLACES[Math.floor(Math.random() * MOCK_PLACES.length)];

    const res = http.get(
        `http://localhost:18082/api/v1/restaurants/quick-search/${placeId}`,
        { headers: { 'Authorization': `Bearer ${data.token}` } }
    );

    check(res, {
        'status is 200': (r) => r.status === 200,
        'has restaurant': (r) => r.json('restaurant') !== null,
    });

    sleep(1);
}
