import { fail } from "k6";
import { adminApi } from "./adminApi.js";
import { getAdApi } from "./getAdApi.js";

const baseUrl = "http://api";

export const options = {
  discardResponseBodies: true,
  scenarios: {
    addAd: {
      executor: "per-vu-iterations",
      exec: "addAd",
      vus: 10,
      iterations: 100,
      maxDuration: "1m",
      gracefulStop: "1s",
      tags: { my_custom_tag: "addAd" },
      env: { MYVAR: "addAd" },
    },
    getAd: {
      executor: "constant-vus",
      exec: "getAd",
      vus: 500,
      duration: "30s",
      startTime: "0s",
      gracefulStop: "1s",
      tags: { my_custom_tag: "getAd" },
      env: { MYVAR: "getAd" },
    },
  },
};

export function addAd() {
  if (__ENV.MYVAR != "addAd") fail();
  adminApi(baseUrl);
}

export function getAd() {
  if (__ENV.MYVAR != "getAd") fail();
  getAdApi(baseUrl);
}
