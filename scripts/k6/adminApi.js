import { check } from "k6";
import http from "k6/http";
import {
  randomAdTitle,
  randomCountry,
  randomGender,
  randomPlatform,
  randomTimestamp,
} from "./utils.js";

export function adminApi(baseUrl) {
  let url = `${baseUrl}/api/v1/ad`;
  let data = {
    title: randomAdTitle(),
    startAt: randomTimestamp(false),
    endAt: randomTimestamp(true),
    conditions: {
      ageStart: Math.floor(Math.random() * 50),
      ageEnd: Math.floor(Math.random() * 50) + Math.floor(Math.random() * 50),
      gender: randomGender(Math.floor(Math.random() * 3)),
      country: randomCountry(Math.floor(Math.random() * 5)),
      platform: randomPlatform(Math.floor(Math.random() * 3)),
    },
  };

  let res = http.post(url, JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
  });

  check(res, { "status is 200": (r) => r.status === 200 });
}
