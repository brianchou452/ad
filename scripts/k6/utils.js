export const randomOffset = () => Math.floor(Math.random() * 100);
export const randomLimit = () => Math.floor(Math.random() * 100) + 1;
export const randomAge = () => Math.floor(Math.random() * 101);
export const randomGender = (count = 1) => getRandomElements(genders, count);
export const randomPlatform = (count = 1) =>
  getRandomElements(platforms, count);
export const randomCountry = (count = 1) => getRandomElements(countries, count);
export const randomAdTitle = () => getRandomElements(titles, 1);
export const randomTimestamp = (offset = false) => {
  const currentDate = new Date();
  const randomOffset = Math.floor(Math.random() * 12 * 60 * 60 * 1000);
  if (offset) {
    currentDate.setDate(currentDate.getDate() + 1);
  }
  const timestamp = new Date(
    currentDate.getTime() - randomOffset
  ).toISOString();
  return timestamp;
};

const getRandomElements = (arr, count) => {
  const selectedElements = new Set();
  while (selectedElements.size < count) {
    selectedElements.add(arr[Math.floor(Math.random() * arr.length)]);
  }
  return Array.from(selectedElements);
};

const titles = [
  "AD 1",
  "AD 2",
  "AD 3",
  "AD 4",
  "AD 5",
  "我超級",
  "超級",
  "想去Dcard",
  "Dcard",
  "的拉拉",
];
const genders = ["M", "F"];
const platforms = ["android", "ios", "web"];
const countries = ["TW", "JP", "HK", "VN"];

// const countries = [
//   "AF",
//   "AX",
//   "AL",
//   "DZ",
//   "AS",
//   "AD",
//   "AO",
//   "AI",
//   "AQ",
//   "AG",
//   "AR",
//   "AM",
//   "AW",
//   "AU",
//   "AT",
//   "AZ",
//   "BS",
//   "BH",
//   "BD",
//   "BB",
//   "BY",
//   "BE",
//   "BZ",
//   "BJ",
//   "BM",
//   "BT",
//   "BO",
//   "BQ",
//   "BA",
//   "BW",
//   "BV",
//   "BR",
//   "IO",
//   "BN",
//   "BG",
//   "BF",
//   "BI",
//   "CV",
//   "KH",
//   "CM",
//   "CA",
//   "KY",
//   "CF",
//   "TD",
//   "CL",
//   "CN",
//   "CX",
//   "CC",
//   "CO",
//   "KM",
//   "CG",
//   "CD",
//   "CK",
//   "CR",
//   "CI",
//   "HR",
//   "CU",
//   "CW",
//   "CY",
//   "CZ",
//   "DK",
//   "DJ",
//   "DM",
//   "DO",
//   "EC",
//   "EG",
//   "SV",
//   "GQ",
//   "ER",
//   "EE",
//   "SZ",
//   "ET",
//   "FK",
//   "FO",
//   "FJ",
//   "FI",
//   "FR",
//   "GF",
//   "PF",
//   "TF",
//   "GA",
//   "GM",
//   "GE",
//   "DE",
//   "GH",
//   "GI",
//   "GR",
//   "GL",
//   "GD",
//   "GP",
//   "GU",
//   "GT",
//   "GG",
//   "GN",
//   "GW",
//   "GY",
//   "HT",
//   "HM",
//   "VA",
//   "HN",
//   "HK",
//   "HU",
//   "IS",
//   "IN",
//   "ID",
//   "IR",
//   "IQ",
//   "IE",
//   "IM",
//   "IL",
//   "IT",
//   "JM",
//   "JP",
//   "JE",
//   "JO",
//   "KZ",
//   "KE",
//   "KI",
//   "KP",
//   "KR",
//   "KW",
//   "KG",
//   "LA",
//   "LV",
//   "LB",
//   "LS",
//   "LR",
//   "LY",
//   "LI",
//   "LT",
//   "LU",
//   "MO",
//   "MK",
//   "MG",
//   "MW",
//   "MY",
//   "MV",
//   "ML",
//   "MT",
//   "MH",
//   "MQ",
//   "MR",
//   "MU",
//   "YT",
//   "MX",
//   "FM",
//   "MD",
//   "MC",
//   "MN",
//   "ME",
//   "MS",
//   "MA",
//   "MZ",
//   "MM",
//   "NA",
//   "NR",
//   "NP",
//   "NL",
//   "NC",
//   "NZ",
//   "NI",
//   "NE",
//   "NG",
//   "NU",
//   "NF",
//   "MP",
//   "NO",
//   "OM",
//   "PK",
//   "PW",
//   "PS",
//   "PA",
//   "PG",
//   "PY",
//   "PE",
//   "PH",
//   "PN",
//   "PL",
//   "PT",
//   "PR",
//   "QA",
//   "RE",
//   "RO",
//   "RU",
//   "RW",
//   "BL",
//   "SH",
//   "KN",
//   "LC",
//   "MF",
//   "PM",
//   "VC",
//   "WS",
//   "SM",
//   "ST",
//   "SA",
//   "SN",
//   "RS",
//   "SC",
//   "SL",
//   "SG",
//   "SX",
//   "SK",
//   "SI",
//   "SB",
//   "SO",
//   "ZA",
//   "GS",
//   "SS",
//   "ES",
//   "LK",
//   "SD",
//   "SR",
//   "SJ",
//   "SE",
//   "CH",
//   "SY",
//   "TW",
//   "TJ",
//   "TZ",
//   "TH",
//   "TL",
//   "TG",
//   "TK",
//   "TO",
//   "TT",
//   "TN",
//   "TR",
//   "TM",
//   "TC",
//   "TV",
//   "UG",
//   "UA",
//   "AE",
//   "GB",
//   "US",
//   "UM",
//   "UY",
//   "UZ",
//   "VU",
//   "VE",
//   "VN",
//   "VG",
//   "VI",
//   "WF",
//   "EH",
//   "YE",
//   "ZM",
//   "ZW",
// ];
