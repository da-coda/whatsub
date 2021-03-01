import axios from 'axios'

/**
 *
 * @return {Promise<{uuid: *, key: *}>}
 */
function createGame () {
  return axios.get('/game/create')
    .then((response) => {
      return {
        uuid: response.data.Payload.UUID,
        key: response.data.Payload.Key
      }
    })
}

/**
 *
 * @param {string} username
 * @param {string} link
 * @return {WebSocket}
 */
function joinGame (username, link) {
  const url = new URL('/game/' + link + '/join?name=' + username, window.location.href)
  url.protocol = url.protocol.replace('http', 'ws')
  url.protocol = url.protocol.replace('https', 'wss')
  return new WebSocket(url.href)
}

export { joinGame, createGame }
