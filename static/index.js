/* trunk-ignore-all(prettier) */
const modified = {}
const api = `/api`

async function poll() {
    const response = await fetch(api + "/getLogs")
    const logs = await response.json()
    update(logs)
}

function update(logs) {
    for (let k in logs) {
        arr = logs[k]
        const bot_name = arr.Name
        if (modified[bot_name] != arr.Timestamp) {
            if (modified[bot_name] == undefined) {
                generateConsole(bot_name)
                //rouind to integer
                modified[bot_name] = arr.Timestamp
                setInterval(() => {
                    document.getElementById(bot_name).querySelector(".console__date-text").innerHTML = `${Math.round(Date.now() / 1000 - modified[bot_name])} seconds ago`
                    rd = Math.round(Date.now() / 1000 - modified[bot_name]) * 10
                    if (rd > 255)
                        rd = 255
                    grn = 255 - rd / 2
                    document.getElementById(bot_name + "_title").style.color = `rgba(${rd},${grn},0,1)`
                }, 1000)
            }
            modified[bot_name] = arr.Timestamp
            const element = document.getElementById(bot_name)
            const consoleText = element.querySelector(".console__text")
            arr = arr.Logs
            // if (bot_name == "server_log") {
            //     arr = arr[0].split("\n")
            //     arr = arr.reverse()

            // }
            arr = arr.reverse()
            //reverse arr
            arr = arr.map((ell) => {
                msg = ell.split(" ")
                type = msg.shift()
                filteredMsg = msg.join(" ").replace("<", "&lt;").replace(">", "&gt;").replace("/", "&#x2F;").replace("\\", "&#39;")
                if (type === "INFO") return `<p class="console__info">${filteredMsg}</p>`
                if (type === "ERROR") return `<p class="console__error">${filteredMsg}</p>`
                if (type === "WARN") return `<p class="console__warn">${filteredMsg}</p>`
                return `<p class="console__msg">${filteredMsg}</p>`
            })

            element.classList.add("_anim")
            consoleText.innerHTML = arr.join("")
            consoleText.scrollTop = consoleText.scrollHeight
        }
    }

    elements = document.getElementsByClassName("console")
    elementsArray = Array.from(elements)
    setTimeout(() => elementsArray.forEach(element => element.classList.remove("_anim")), 300)
}

function deleteConsole(name) {
    document.getElementById(name).remove()
    delete modified[name]
}
function format_date(date) {
    hour = date.getHours()
    minute = date.getMinutes()
    day = date.getDate()
    month = date.getMonth() + 1
    year = date.getFullYear()

    return `${zeroPad(day, 2)}.${zeroPad(month, 2)}.${year} ${zeroPad(hour, 2)}:${zeroPad(minute, 2)}`
}

const zeroPad = (num, places) => String(num).padStart(places, '0')

function generateConsole(name) {
    console = ""
    if (name == "server_log") {
        console =
            `
        <div class="console" id="${name}">
        <div class="console__inner">
            <div class="console__console" style="width:100%">
                <div class="console__name" style="width:100%">
                    <h3 id="${name}_title" class="console__title" style="color:rgba(255,255,255,1)">${name}</h3>
                </div>
                <div class="console__date" style="width:100%">
                    <h3 class="console__title console__date-text"></h3>
                </div>
                <div class="console__text" style="width:100%"></div>
            </div>
        </div>
    </div>`
    } else {
        console =
            `
        <div class="console" id="${name}">
        <div class="console__inner">
            <div class="console__buttons">
                <button class="console__button" onclick="sendsig('${name}','r')">R</button>
                <button class="console__button" onclick="sendsig('${name}','d'); deleteConsole('${name}')">D</button>
                <button class="console__button" onclick="sendsig('${name}','a')">A</button>
            </div>
            <div class="console__console">
                <div class="console__name">
                    <h3 id="${name}_title" class="console__title" style="color:rgba(255,255,255,1)">${name}</h3>
                </div>
                <div class="console__date">
                    <h3 class="console__title console__date-text"></h3>
                </div>
                <div class="console__text"></div>
            </div>
        </div>
    </div>
`}
    document.getElementById("main").innerHTML += console
    elements = document.getElementsByClassName("console")
    elementsArray = Array.from(elements)
    elementsArray.forEach(element => element.querySelector(".console__text").scrollTop = element.querySelector(".console__text").scrollHeight)
}

async function sendsig(name, sig) {
    const response = await fetch(api + "/sig" + sig, { method: "POST", body: name })
}

setInterval(poll, 2000)