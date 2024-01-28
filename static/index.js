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
            console.log(arr);
            if (modified[bot_name] == undefined) {
                generateConsole(bot_name)
                    //rouind to integer
                modified[bot_name] = arr.Timestamp
                setInterval(() => {
                    document.getElementById(bot_name).querySelector(".console__date-text").innerHTML = `${Math.round(Date.now() / 1000 - modified[bot_name])} seconds ago`
                }, 1000)
            }
            modified[bot_name] = arr.Timestamp
            const element = document.getElementById(bot_name)
            const consoleText = element.querySelector(".console__text")
                //reverse arr
            arr = arr.Logs.reverse()
            arr = arr.map((ell) => {
                msg = ell.split(" ")
                type = msg.shift()
                console.log(type);
                filteredMsg = msg.join(" ").replace("<", "&lt;").replace(">", "&gt;").replace("/", "&#x2F;").replace("\\", "&#39;")
                if (type === "INFO") return `<p class="console__info">${filteredMsg}</p>`
                if (type === "ERROR") return `<p class="console__error">${filteredMsg}</p>`
                if (type === "WARN") return `<p class="console__warn">${filteredMsg}</p>`
                return `<p class="console__msg">${filteredMsg}</p>`
            })

            element.classList.add("_anim")
            consoleText.innerHTML = arr.join("<br>")
            consoleText.scrollTop = consoleText.scrollHeight
        }
    }

    elements = document.getElementsByClassName("console")
    elementsArray = Array.from(elements)
    setTimeout(() => elementsArray.forEach(element => element.classList.remove("_anim")), 300)
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
    const console =
        `
<div class="console" id="${name}">
    <div class="console__inner">
        <div class="console__buttons">
            <button class="console__button" onclick="reset(${name})">R</button>
        </div>
        <div class="console__console">
            <div class="console__name">
                <h3 class="console__title">${name}</h3>
            </div>
            <div class="console__date">
                <h3 class="console__title console__date-text"></h3>
            </div>
            <div class="console__text"></div>
        </div>
    </div>
</div>
`
    document.getElementById("main").innerHTML += console
    elements = document.getElementsByClassName("console")
    elementsArray = Array.from(elements)
    elementsArray.forEach(element => element.querySelector(".console__text").scrollTop = element.querySelector(".console__text").scrollHeight)
}

async function reset(name) {
    const response = await fetch(api + "/reset", { method: "POST", body: name.id })
}

setInterval(poll, 2000)