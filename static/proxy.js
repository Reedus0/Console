const api = `/api`
async function http_get() {
    const response = await fetch(api + "/getProxyes")
    const proxy = await response.text()
    document.getElementById("proxy_list").innerHTML = proxy
}
String.prototype.replaceAll = function(search, replacement) {
    var target = this;
    return target.split(search).join(replacement);
};
async function check() {
    data = document.getElementById("proxy_list").value
    document.getElementById("proxy_check").innerHTML = "Загрузка..."

    console.log(data);
    const response = await fetch(api + "/checkProxy", { method: "POST", body: data })
    proxy = await response.text()
    proxy = proxy.replaceAll("false", "❌").replaceAll("true", "✅")
    document.getElementById("proxy_check").innerHTML = proxy

}
async function update() {
    data = document.getElementById("proxy_list").value
    const response = await fetch(api + "/updateProxy", { method: "POST", body: data })
    res = await response.text()
    alert(res)

}