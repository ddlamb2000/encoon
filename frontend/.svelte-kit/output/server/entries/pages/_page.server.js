function load({ cookies }) {
  const visited = cookies.get("visited");
  console.log("cookies.set");
  cookies.set("visited", "true", { path: "/" });
  return {
    visited: visited === "true"
  };
}
export {
  load
};
