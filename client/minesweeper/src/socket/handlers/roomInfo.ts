import { gameState } from '@/store'

type RoomId = {
    id: number,
    inviteCode: string
}

export default (data: RoomId) => {
    const { id, inviteCode } = data
    gameState.roomId = id
    gameState.inviteCode = inviteCode
    if (![undefined, -1].includes(id)) {
        document.title = `#${id} Minesweeper`
    }
    // Change search query
    const url = new URL(window.location.href)
    url.searchParams.set('room', data.id.toString())
    history.replaceState({}, '', url)
}
