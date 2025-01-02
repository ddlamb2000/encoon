import { p as postMessage } from "../../../../chunks/kafka.js";
const POST = async ({ params, request, url }) => {
  const auth = request.headers.get("Authorization");
  if (auth === "" || auth === void 0 || auth.length < 10) {
    console.error("Not authorized");
    return json({ error: "Not authorized" }, { status: 401 });
  }
  const tokenString = auth?.substring(7);
  try {
    const arrayToken = tokenString.split(".");
    const tokenPayload = JSON.parse(atob(arrayToken[1]));
    const now = (/* @__PURE__ */ new Date()).toISOString();
    const nowDate = Date.parse(now);
    const tokenExpirationDate = Date.parse(tokenPayload.expires);
    if (nowDate > tokenExpirationDate) {
      return json({ error: "Authorization expided" }, { status: 401 });
    }
  } catch (error) {
    return json({ error: "Not authorized" }, { status: 401 });
  }
  return postMessage(params, request, url);
};
export {
  POST
};
