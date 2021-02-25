import http from 'k6/http';
import ws from 'k6/ws';
import { check } from 'k6';

export default function () {
    const wsres = ws.connect("ws://localhost:8000/game/d4553dea-9718-4e2a-af88-cff9693fee60/join?name=daniel" + Math.floor(Math.random() * Math.floor(10000000)), null, function (socket) {
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
    })
    check(wsres, { 'status is 101': (r) => r && r.status === 101 });
}