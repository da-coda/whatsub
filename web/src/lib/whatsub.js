import axios from 'axios'

export const baseUrl = window.location.protocol +
  '//' + window.location.hostname + ':' + window.location.port

/**
 *
 * @return {Promise<{uuid: *, key: *}>}
 */
function createGame () {
  return axios.post('/game/create')
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

/**
 * Start the first round of the game
 * @param {string} link Shortend game id
 * @return {Promise<AxiosResponse<any>>}
 */
function startGame (link) {
  return axios.get(baseUrl + '/game/' + link + '/start')
}

/**
 * Ask the game server to give the current game state
 * @param {string} link Shortened game idk
 * @return {Promise<AxiosResponse<any>>}
 */
function askGameState (link) {
  return axios.get(baseUrl + '/game/' + link + '/status')
}
export { joinGame, createGame, startGame, askGameState }
