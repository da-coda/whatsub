let conn;
let uuid = "";
let useUuidField = false;

function createGame() {
    var req = new XMLHttpRequest()
    req.onreadystatechange = function () {
        if (req.readyState === 4 && req.status === 200) {
            uuid = JSON.parse(req.responseText).uuid
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
    conn = new WebSocket("ws://localhost:8000/game/" + uuid + "/join?name=" + name)
    document.getElementById("join").style.display = "none"
    document.getElementById("lobby").style.display = "block"

    conn.onmessage = function (event) {
        var msg = JSON.parse(event.data);
        switch (msg.Type) {
            case "round":
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
                let scoreList = document.getElementById("scoreList").innerHTML = ""
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
    document.getElementById("lobby").style.display = "none"
    document.getElementById("round").style.display = "block"
}

function answerPost(answer) {
    let payload = {
        Type: "answer",
        Payload: {
            Answer: answer.dataset.answer
        }
    }
    conn.send(JSON.stringify(payload))
}

Element.prototype.remove = function () {
    this.parentElement.removeChild(this);
}

function htmlDecode(input) {
    var doc = new DOMParser().parseFromString(input, "text/html");
    return doc.documentElement.textContent;
}