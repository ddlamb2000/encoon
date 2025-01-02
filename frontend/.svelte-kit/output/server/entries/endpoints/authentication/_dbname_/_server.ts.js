import { p as postMessage } from "../../../../chunks/kafka.js";
const POST = async ({ params, request, url }) => {
  return postMessage(params, request, url);
};
export {
  POST
};
