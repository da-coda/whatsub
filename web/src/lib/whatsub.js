import axios from 'axios'

export const baseUrl = window.location.protocol +
  '//' + window.location.hostname + ':' + window.location.port

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
 * @param {string} playerUUID
 * @return {WebSocket}
 */
function joinGame (username, link, playerUUID) {
  const url = new URL('/game/' + link + '/join?name=' + username + '&uuid=' + playerUUID, baseUrl)
  url.protocol = url.protocol.replace('http', 'ws')
  url.protocol = url.protocol.replace('https', 'wss')
  return new WebSocket(url.href)
}

export { joinGame, createGame }
