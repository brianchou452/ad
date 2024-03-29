import { check } from "k6";
import http from "k6/http";
import {
  randomAge,
  randomCountry,
  randomGender,
  randomLimit,
  randomOffset,
  randomPlatform,
} from "./utils.js";

export function getAdApi(baseUrl) {
  let url = `${baseUrl}/api/v1/ad?offset=${randomOffset()}&limit=${randomLimit()}`;

  if (Math.random() > 0.5) url += `&age=${randomAge()}`;
  if (Math.random() > 0.5) url += `&gender=${randomGender()[0]}`;
  if (Math.random() > 0.5) url += `&platform=${randomPlatform()[0]}`;
  if (Math.random() > 0.5) url += `&country=${randomCountry()[0]}`;

  let res = http.get(url, {
    headers: { "Content-Type": "application/json" },
    tags: {
      name: "Get AD API",
    },
  });

  check(
    res,
    { "Get AD API status is 200": (r) => r.status === 200 },
    { name: "Get AD API" }
  );
}
