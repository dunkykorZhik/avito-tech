import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
export let options = {
  vus: 1000,               // number of virtual users
  duration: '1m',          // test length (adjust to stress/soak test)
  thresholds: {
    http_req_duration: ['p(95)<50'],     // 95% of requests < 50ms
    http_req_failed: ['rate<0.0001'],    // <0.01% failures (99.99% success)
  },
};

// Replace with your API base
const BASE_URL = 'http://localhost:8080';

// Test credentials
const USERNAME = 'testuser';
const PASSWORD = 'testpass';

// Shared auth token (gets refreshed once for all VUs)
let token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlciIsImV4cCI6MTc1NTc5MDk2NywiaWF0IjoxNzU1NzkwNjY3fQ.z8wl01OrK9pIp5MbzIVCEPqgZ8dA2Jvjr6Qymvow5xc;

// Authenticate once before test
export function setup() {
    let loginRes = http.post(`${BASE_URL}/api/auth`, JSON.stringify({
      username: USERNAME,
      password: PASSWORD,
    }), { headers: { 'Content-Type': 'application/json' } });
  
    console.log("Login status:", loginRes.status);
    console.log("Login body:", loginRes.body);
  
    check(loginRes, { 'login succeeded': (r) => r.status === 200 });
  
    let body = JSON.parse(loginRes.body);
  
    // Adjust this depending on your API response format
    return { token: body.jwt };
  }
  

// Each virtual user runs this
export default function (data) {
  let authHeaders = {
    headers: {
      'Authorization': `Bearer ${data.token}`,
      'Content-Type': 'application/json',
    },
  };

  // 1. Get info
  let infoRes = http.get(`${BASE_URL}/api/info`, authHeaders);
  check(infoRes, { 'info ok': (r) => r.status === 200 });

  // 2. Send coins
  let sendRes = http.post(`${BASE_URL}/api/sendCoin`, JSON.stringify({
    toUser: 'korkem',
    amount: Math.floor(Math.random() * 10) + 1,
  }), authHeaders);
  check(sendRes, { 'send ok': (r) => r.status === 200 });

  // 3. Buy item
  let buyRes = http.get(`${BASE_URL}/api/buy/item1`, authHeaders);
  check(buyRes, { 'buy ok': (r) => r.status === 200 });

  // Short pause to simulate user think-time
  sleep(1);
}
