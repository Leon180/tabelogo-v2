import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// 自定義指標
const errorRate = new Rate('errors');
const restaurantSearchDuration = new Trend('restaurant_search_duration');
const successfulRequests = new Counter('successful_requests');

// Mock 餐廳 Place IDs
const MOCK_PLACES = [
    'mock_tokyo_ramen_1',
    'mock_osaka_sushi_1',
    'mock_kyoto_tempura_1',
    'mock_fukuoka_ramen_1',
    'mock_tokyo_sushi_1',
];

// 測試配置
export const options = {
    stages: [
        { duration: '30s', target: 10 },   // 熱身：30秒內增加到10個用戶
        { duration: '1m', target: 50 },    // 增加負載：1分鐘內增加到50個用戶
        { duration: '3m', target: 50 },    // 穩定負載：維持50個用戶3分鐘
        { duration: '1m', target: 100 },   // 峰值測試：1分鐘內增加到100個用戶
        { duration: '2m', target: 100 },   // 峰值維持：維持100個用戶2分鐘
        { duration: '30s', target: 0 },    // 冷卻：30秒內降到0
    ],
    thresholds: {
        'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95%請求<500ms, 99%<1s
        'http_req_failed': ['rate<0.01'],                  // 錯誤率<1%
        'errors': ['rate<0.01'],                           // 自定義錯誤率<1%
    },
};

// 設置階段 - 獲取認證 token
export function setup() {
    // 嘗試註冊新用戶
    const timestamp = Date.now();
    const registerRes = http.post(
        'http://localhost:8080/api/v1/auth/register',
        JSON.stringify({
            email: `k6test${timestamp}@example.com`,
            username: `k6test${timestamp}`,
            password: 'Test1234!',
        }),
        {
            headers: { 'Content-Type': 'application/json' },
        }
    );

    let token;
    if (registerRes.status === 201 || registerRes.status === 200) {
        const body = registerRes.json();
        token = body.access_token;
    } else {
        console.error(`Registration failed with status ${registerRes.status}`);
        console.error(`Response: ${registerRes.body}`);
        throw new Error('Failed to register user');
    }

    if (!token) {
        throw new Error('Token is null or undefined');
    }

    console.log(`Setup complete. Token obtained: ${token.substring(0, 20)}...`);
    return { token };
}

// 主測試函數
export default function (data) {
    // 隨機選擇一個 Mock 餐廳
    const placeId = MOCK_PLACES[Math.floor(Math.random() * MOCK_PLACES.length)];

    // 發送請求
    const startTime = new Date();
    const res = http.get(
        `http://localhost:18082/api/v1/restaurants/quick-search/${placeId}`,
        {
            headers: {
                'Authorization': `Bearer ${data.token}`,
            },
            tags: { name: 'QuickSearch' },
        }
    );
    const duration = new Date() - startTime;

    // 記錄自定義指標
    restaurantSearchDuration.add(duration);

    // 驗證響應
    const checkResult = check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
        'has restaurant data': (r) => {
            try {
                const body = r.json();
                return body.restaurant && body.restaurant.name;
            } catch (e) {
                return false;
            }
        },
        'restaurant name is mock': (r) => {
            try {
                const body = r.json();
                return body.restaurant.name && body.restaurant.name.includes('Mock');
            } catch (e) {
                return false;
            }
        },
    });

    // 記錄錯誤
    errorRate.add(!checkResult);

    // 記錄成功請求
    if (checkResult) {
        successfulRequests.add(1);
    }

    // 模擬真實用戶行為 - 隨機等待
    sleep(Math.random() * 2 + 1); // 1-3秒隨機等待
}

// 清理階段
export function teardown(data) {
    console.log('Test completed successfully');
}
