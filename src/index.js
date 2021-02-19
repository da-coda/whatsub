let conn;


function createGame() {
    var req = new XMLHttpRequest()
    req.onreadystatechange = function() {
        if (req.readyState === 4 && req.status === 200)
            document.getElementById("gameUuid").innerText = req.responseText
    }
    req.open("GET", "/game/create", true); // true for asynchronous
    req.send(null);
}

function joinGame() {
    let uuid = document.getElementById("joinUuid").value
    let name = document.getElementById("joinName").value
    conn = new WebSocket("wss://whatsub.danielroehrig.de/game/" + uuid  + "/join")
    let infoBox = document.getElementById("websocketMessages")
    let postTitle = document.getElementById("postTitle")
    let postContent = document.getElementById("postContent")
    conn.onopen = function (event) {
        conn.send(name);
    };

    conn.onmessage = function (event) {
        var msg = JSON.parse(event.data);
        console.log(msg)
        console.log(event.data)
        switch (msg.Type) {
            case "round":
                postTitle.innerText = msg.PostTitle
                postContent.innerHTML = msg.PostText
                if (msg.Image !== "") {
                    let oldImg = document.getElementById("postImage");
                    if(oldImg !== null) {
                        oldImg.remove()
                    }
                    let img = document.createElement("img")
                    img.id = "postImage"
                    img.src = msg.Image
                    img.height = "150"
                    postContent.parentNode.insertBefore(img, postContent.nextSibling);
                }
                break;
            default:
                infoBox.innerText = msg

        }
    }
}

function startGame() {
    let uuid = document.getElementById("startUuid").value
    var req = new XMLHttpRequest()
    req.onreadystatechange = function() {
        if (req.readyState === 4 && req.status === 200)
            document.getElementById("gameUuid").innerText = req.responseText
    }
    req.open("GET", "/game/" + uuid  + "/start", true); // true for asynchronous
    req.send(null);
}

function answerPost() {
    let answer = document.getElementById("answer").value
    conn.send(answer)
}

Element.prototype.remove = function() {
    this.parentElement.removeChild(this);
}