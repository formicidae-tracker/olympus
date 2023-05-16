var defaultTarget = "";
BACKEND_PORT = process.env["BACKEND_PORT"];

module.exports = [
  {
    context: ["/api/**"],
    target:
    BACKEND_PORT ?  "http://localhost:" + BACKEND_PORT: defaultTarget,
    secure: true,
    changeOrigin: true,
  },
];
