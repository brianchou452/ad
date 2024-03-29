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
    title: randomAdTitle()[0],
    startAt: randomTimestamp(false),
    endAt: randomTimestamp(true),
    conditions: {
      // ageStart: Math.floor(Math.random() * 50),
      // ageEnd: Math.floor(Math.random() * 50) + Math.floor(Math.random() * 50),
      // gender: randomGender(Math.floor(Math.random() * 3)),
      // country: randomCountry(Math.floor(Math.random() * 5)),
      // platform: randomPlatform(Math.floor(Math.random() * 3)),
    },
  };

  if (Math.random() > 0.5) {
    data.conditions.ageStart = Math.floor(Math.random() * 50);
    data.conditions.ageEnd =
      data.conditions.ageStart + Math.floor(Math.random() * 50);
  }
  if (Math.random() > 0.5) {
    data.conditions.gender = randomGender(1 + Math.floor(Math.random() * 1));
  }
  if (Math.random() > 0.5) {
    data.conditions.country = randomCountry(1 + Math.floor(Math.random() * 3));
  }
  if (Math.random() > 0.5) {
    data.conditions.platform = randomPlatform(
      1 + Math.floor(Math.random() * 2)
    );
  }

  let res = http.post(url, JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
    tags: {
      name: "Admin API",
    },
  });

  if (res.status !== 200) {
    console.log(data);
  }

  check(
    res,
    { "Admin API status is 200": (r) => r.status === 200 },
    { name: "Admin API" }
  );
}
