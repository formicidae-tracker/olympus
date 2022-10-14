const PROXY_CONFIG = [
	  {
		  context: [
			  "/olympus",
			  "/api",
		  ],
		  "target": "http://localhost",
		  "secure": false
	  }
]

module.exports = PROXY_CONFIG;
