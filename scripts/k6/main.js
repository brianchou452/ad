// import { adminApi } from "./adminApi.js";
import { getAdApi } from "./getAdApi.js";
const baseUrl = "http://api";

export let options = {
  stages: [
    // Ramp-up from 1 to 5 VUs in 5s
    { duration: "5s", target: 500 },
    // Stay at rest on 5 VUs for 10s
    { duration: "20s", target: 500 },
    // Ramp-down from 5 to 0 VUs for 5s
    // { duration: "5s", target: 0 },
  ],
};

export default function () {
  // adminApi(baseUrl);
  getAdApi(baseUrl);
}
