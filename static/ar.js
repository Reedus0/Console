const api = `/api`
data = []
String.prototype.replaceAll = function(search, replacement) {
    var target = this;
    return target.split(search).join(replacement);
};
async function onsel(){
    console.log(document.getElementById("combo").value);
    textel = document.getElementById("text")
    text = data.filter(x => x.Name == document.getElementById("combo").value)[0].Text
    console.log(text)
    text = text.replaceAll("==", "<br><br>", -1)
    text = text.replaceAll("_", "")
    textel.innerHTML = text
}

async function http_get() {
    const response = await fetch(api + "/getFolders")
    const folders = await response.json()
    data = folders
    combbx = document.getElementById("combo")
    for (const folder of folders) {
        console.log(folder)
        combo.innerHTML += `<option value="${folder.Name}">${folder.Name}</option>`
    }
    onsel()
}

async function autoreg(){
    if (document.getElementById("prefix").value == "" || document.getElementById("count").value == "" || document.getElementById("combo").value == "") {
        alert("All fields are required")
        return
    }
    const response = await fetch(api + "/autoReg", { method: "POST", body: document.getElementById("prefix").value +"\n"+ document.getElementById("count").value +"\n"+ document.getElementById("combo").value+"\n"+ document.getElementById("tables").value })
    document.location.href = "index.html"
}
http_get()