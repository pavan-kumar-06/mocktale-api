import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    realistic_load: {
      executor: "ramping-arrival-rate",
      startRate: 1000,
      timeUnit: "1s",
      preAllocatedVUs: 50,
      maxVUs: 500,
      stages: [
        { target: 5000, duration: "30s" },
        { target: 5000, duration: "5m" },
        { target: 1000, duration: "30s" },
      ],
    },
  },
  thresholds: {
    "http_req_duration{api:vote}": ["p(95)<1000"],
    "http_req_duration{api:movie}": ["p(95)<500"],
    "http_req_duration{api:rating}": ["p(95)<500"],
    http_req_failed: ["rate<0.05"],
  },
};

const baseURL = "http://localhost:8080";

// Your 5 movie slugs with 0 initial votes
const movieSlugs = ["100-kaadhal-2019", "hustlers-2019", "this-way-up-2019", "the-cursed-2020", "beyonce-bowl-2024"];

export default function () {
  const randomSlug = movieSlugs[Math.floor(Math.random() * movieSlugs.length)];
  const apiType = Math.random();

  let response;
  let apiName;

  if (apiType < 0.4) {
    apiName = "movie";
    response = http.get(`${baseURL}/api/movies/${randomSlug}`, {
      tags: { api: apiName },
    });
  } else if (apiType < 0.7) {
    apiName = "rating";
    response = http.get(`${baseURL}/api/ratings/${randomSlug}`, {
      tags: { api: apiName },
    });
  } else {
    apiName = "vote";
    response = http.post(
      `${baseURL}/api/user/vote`,
      JSON.stringify({
        movie_slug: randomSlug,
        option_chosen: Math.floor(Math.random() * 4), // 0 or 1
      }),
      {
        headers: { "Content-Type": "application/json" },
        tags: { api: apiName },
      }
    );
  }

  check(response, {
    [`${apiName} status 2xx`]: (r) => r.status >= 200 && r.status < 300,
    [`${apiName} response time`]: (r) => r.timings.duration < 2000,
  });
}
