const { defineConfig } = require("cypress");

module.exports = defineConfig({
  e2e: {
    supportFile: false,
    allowCypressEnv: false,
    baseUrl: "http://localhost:36000",
  },
});
