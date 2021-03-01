let conn
let uuid = ""
let useUuidField = false
let joinUrl = ""

window.onload = function(){
    var pathArray = window.location.pathname.split('/');
    console.log(pathArray)
    if(pathArray[1] === "joinMe") {
        uuid = pathArray[2]
        document.getElementById("start").style.display = "none"
        document.getElementById("join").style.display = "block"
    }
}

function createGame() {
    var req = new XMLHttpRequest()
    req.onreadystatechange = function () {
        if (req.readyState === 4 && req.status === 200) {
            let response = JSON.parse(req.responseText);
            uuid = response.Payload.UUID
            joinUrl = new URL("/joinMe/" + response.Payload.Key, window.location.href);
            document.getElementById("joinUrl").innerText = joinUrl
        }
        document.getElementById("start").style.display = "none"
        document.getElementById("join").style.display = "block"
    }
    req.open("GET", "/game/create", true); // true for asynchronous
    req.send(null);
}

function joinLobby() {
    document.getElementById("start").style.display = "none"
    document.getElementById("join").style.display = "block"
    document.getElementById("joinUuid").style.display = "block"
    useUuidField = true
}

function joinGame() {
    let name = document.getElementById("joinName").value
    if (useUuidField) {
        uuid = document.getElementById("joinUuid").value
    }
    let playerUuid = uuidv4()
    var url = new URL("/game/" + uuid + "/join?name=" + name + "&uuid=e603cadb-6dde-4072-a468-9deaf66e8fc6", window.location.href);
    url.protocol = url.protocol.replace('http', 'ws');
    url.protocol = url.protocol.replace('https', 'wss');
    conn = new WebSocket(url.href)
    document.getElementById("join").style.display = "none"
    document.getElementById("lobby").style.display = "block"

    conn.onmessage = function (event) {
        var msg = JSON.parse(event.data);
        switch (msg.Type) {
            case "round":
                answered = false
                document.getElementById("lobby").style.display = "none"
                document.getElementById("round").style.display = "block"
                document.getElementById("buttons").innerHTML = ""
                document.getElementById("postImage").src = ""

                document.getElementById("postTitle").innerText = msg.Payload.Post.Title
                document.getElementById("postContent").innerHTML = htmlDecode(msg.Payload.Post.Content)
                if(msg.Payload.Post.Type === "Image"){
                    document.getElementById("postImage").src = msg.Payload.Post.Url
                }
                let subreddits = msg.Payload.Subreddits
                subreddits.forEach(function (entry) {
                    document.getElementById("buttons").innerHTML += '<div class="col-3 mb-1 me-1">' +
                        '<button id="' + entry + '" onclick="answerPost(this)" data-answer="' + entry + '" type="button" class="btn btn-light w-100 ">' + entry + '</button>' +
                        '</div>'
                })
                break;
            case "score":
                document.getElementById("scoreList").innerHTML = ""
                console.log(msg.Payload.Scores)
                for (const [key, value] of Object.entries(msg.Payload.Scores)) {
                    document.getElementById("scoreList").innerHTML +=  `<a href="#" class="list-group-item list-group-item-action bg-light"> ${key}: ${value} </a>`
                }

                break

            case "answer_correct":
                document.getElementById(msg.Payload.CorrectAnswer).style.backgroundColor = "green"
                break
            default:
                console.log(event.data)

        }
    }
}

function startGame() {
    var req = new XMLHttpRequest()
    req.open("GET", "/game/" + uuid + "/start", true); // true for asynchronous
    req.send(null);
}

let answered = false

function answerPost(answer) {
    if (answered) {
        return
    }
    let payload = {
        Type: "answer",
        Payload: {
            Answer: answer.dataset.answer
        }
    }
    conn.send(JSON.stringify(payload))
    answered = true
}

Element.prototype.remove = function () {
    this.parentElement.removeChild(this);
}

function htmlDecode(input) {
    var doc = new DOMParser().parseFromString(input, "text/html");
    return doc.documentElement.textContent;
}

function copyJoinUrl() {
    navigator.clipboard.writeText(joinUrl).then(function() {
        console.log('Async: Copying to clipboard was successful!');
    }, function(err) {
        console.error('Async: Could not copy text: ', err);
    });
}

function uuidv4() {
    return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
        (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
    );
}
