import userIdHandler from './userId'
import roomIdHandler from './roomId'
import roomCreatedHandler from './roomCreated'
import gameStatHandler from './gameStat'
import playerAliasHandler from './playerAlias'
import playerOfflineHandler from './playerOffline'
import playerOnlineHandler from './playerOnline'
import scoreUpdatedHandler from './scoreUpdated'
import turnPassedHandler from './turnPassed'
import gameEndedHandler from './gameEnded'
import revealHandler from './reveal'
import reconnFailedHandler from './reconnFailed'

type handlerMap = {
    name: string,
    fn: (event: any) => void
}

const getAll = () : handlerMap[] => [
    { name: 'userId', fn: userIdHandler },
    { name: 'roomId', fn: roomIdHandler },
    { name: 'roomCreated', fn: roomCreatedHandler },
    { name: 'reconnFailed', fn: reconnFailedHandler },
    { name: 'gameStat', fn: gameStatHandler },
    { name: 'playerAlias', fn: playerAliasHandler },
    { name: 'playerOffline', fn: playerOfflineHandler },
    { name: 'playerOnline', fn: playerOnlineHandler },
    { name: 'scoreUpdated', fn: scoreUpdatedHandler },
    { name: 'gameEnded', fn: gameEndedHandler },
    { name: 'turnPassed', fn: turnPassedHandler },
    { name: 'reveal', fn: revealHandler }
]

export default { getAll }
