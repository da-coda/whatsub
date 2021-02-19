import http from 'k6/http';
import ws from 'k6/ws';
import { check } from 'k6';

export default function () {
    var url = 'http://localhost:8000/game/create';
    var res = http.get(url)
    check(res, { 'status was 200': (r) => r.status === 200 });
    var uuid = JSON.parse(res.body).uuid
    const wsres = ws.connect("ws://localhost:8000/game/" + uuid + "/join?name=daniel", null, function (socket) {
        socket.on('message', function (data) {
            let payload = {
                Type: "answer",
                Payload: {
                    Answer: "r/memes"
                }
            }
            payload = JSON.stringify(payload);
            console.log(payload)
            socket.send(payload)
        })
        res = http.get("http://localhost:8000/game/" + uuid + "/start")
        check(res, { 'status was 200': (r) => r.status === 200 });
    })
    check(wsres, { 'status is 101': (r) => r && r.status === 101 });
}