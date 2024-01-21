modified = {}

async function poll() {
    const response = await fetch("http://127.0.0.1:9999")
    const logs = await response.json()
    update(logs)
}

function update(logs) {
    for (let arr of logs) {
        console.log(modified)
        console.log(logs)
        const bot_name = arr[0]
        if (modified[bot_name] == undefined) {
            generateConsole(bot_name)
            modified[bot_name] = 0
            setInterval(() => {
                modified[bot_name] += 1
                document.getElementById(bot_name).querySelector(".console__date-text").innerHTML = `${modified[bot_name]} seconds ago`
            }, 1000)
        }
        modified[bot_name] = 0
        const element = document.getElementById(bot_name)
        const consoleText = element.querySelector(".console__text")
        arr = arr.filter((ell, index) => index != 0)
        arr = arr.map((ell) => {
            msg = ell.split(" ")
            type = msg.shift()
            filteredMsg = msg.join(" ").replace("<", "&lt;").replace(">", "&gt;").replace("/", "&#x2F;").replace("\\", "&#39;")
            if (type === "INFO") return `<p class="console__info">${filteredMsg}</p>`
            if (type === "ERROR") return `<p class="console__error">${filteredMsg}</p>`
            return `<p class="console__msg>${filteredMsg}</p3>`
        })

        element.classList.add("_anim")
        consoleText.innerHTML = arr.join("<br>")
        consoleText.scrollTop = consoleText.scrollHeight
    }

    elements = document.getElementsByClassName("console")
    elementsArray = Array.from(elements)
    setTimeout(() => elementsArray.forEach(element => element.classList.remove("_anim")), 1000)
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
}

async function reset(name) {
    const response = await fetch("http://127.0.0.1:9999", { method: "POST", body: name.id })
    const logs = await response.json()
}

setInterval(poll, 2000)
