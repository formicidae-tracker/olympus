var defaultTarget = "";
BACKEND_PORT = process.env["BACKEND_PORT"];

module.exports = [
  {
    context: ["/api/**"],
    target:
      BACKEND_PORT.length == 0
        ? defaultTarget
        : "http://localhost:" + BACKEND_PORT,
    secure: true,
    changeOrigin: true,
  },
];
