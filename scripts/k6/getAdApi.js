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

  const age = randomAge();
  if (age !== 0) url += `&age=${age}`;
  const gender = randomGender();
  if (gender !== "") url += `&gender=${gender}`;
  const platform = randomPlatform();
  if (platform !== "") url += `&platform=${platform}`;
  const country = randomCountry();
  if (country !== "") url += `&country=${country}`;

  let res = http.get(url, {
    headers: { "Content-Type": "application/json" },
  });

  check(res, { "status is 200": (r) => r.status === 200 });
}
