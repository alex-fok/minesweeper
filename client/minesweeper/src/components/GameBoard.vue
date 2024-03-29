<script setup lang='ts'>
import { computed, ref, watch } from 'vue'
import { BOARDSETTING, GAMESTATUS } from '@/config'
import { gameState, roomState } from '@/store'
import socket from '@/socket'
import Block from './GameBlock.vue'
import Copy from './icon/CopyIcon.vue'

const { END, IN_GAME, WAITING_JOIN } = GAMESTATUS
const timeLeft = ref(-1)

const reveal = (i: number) => {
    if (roomState.status !== IN_GAME) return

    const y = Math.floor(i / gameState.boardConfig.size)
    const x = i % gameState.boardConfig.size
    socket.emit('reveal', {x, y})
}
const getInviteUrl = () => {
    const {  protocol, hostname, port, pathname } = window.location
    const portNum = port !== '' ? ':' + port : ''
    return `${protocol}//${hostname}${portNum}${pathname}?join=${roomState.inviteCode}`
}

const mapSize = computed(() => {
    const { size } = gameState.boardConfig
    const { SMALL, MEDIUM } = BOARDSETTING.SIZE
    return (
        size === SMALL ? 'small' :
        size === MEDIUM ? 'medium' :
        'large'
    )
})
// Record array of cursor position
// playerCursor[position] = playerid
const playerCursors = computed(() => {
    const playerIds = Object.keys(gameState.players)
    
    const result : string[] = []
    playerIds.forEach(id => {
        if (id !== roomState.id || (roomState.isPlayer && !gameState.players[roomState.id]?.isTurn))
            result[gameState.players[id].cursor] = id
    })
    return result
})

const updateMousePosition = (position: number) => {
    if (!roomState.isPlayer) return
    socket.emit('share', {
        name: 'playerMousePos', 
        content: JSON.stringify({position})
    })
}

const copyInviteUrl = () => navigator.clipboard.writeText(getInviteUrl())
const lastPlayedBlock = computed(() => {
    const { size } = gameState.boardConfig
    const { x, y } = gameState.lastPlayed
    return size * y + x
})

const isTimerExist = computed(() => {
    return roomState.status === END && gameState.timeLimit
})

watch(() => gameState.lastPlayed.timestamp, () => {
    const compared = gameState.lastPlayed.timestamp
    const timeLimitInSec = gameState.timeLimit * 1000
    if (!compared) {
        timeLeft.value = -1
        return
    }
    const interval = setInterval(() => {
        const isUpdated = compared !== gameState.lastPlayed.timestamp
        if (isUpdated) {
            clearInterval(interval)
            return
        }
        const result = compared + timeLimitInSec - (Date.parse(new Date().toString()))
        if (result > 0) {
            timeLeft.value = Math.floor(result / 1000)
        } else {
            clearInterval(interval)
            timeLeft.value = 0 
        }
    }, 200)
}, {immediate: true})

</script>
<template>
    <div v-if='roomState.status === WAITING_JOIN' class='waiting-container'>
        <div class='waiting-text-wrapper'>
            <div class='waiting-text'>
                <div>Waiting for player to join...</div>
                <div class='subtitle'>{{`${Object.keys(gameState.players).length}/${gameState.capacity}`}}</div>
            </div>
            <div>
                Invite:
                <input
                    v-if='roomState.inviteCode !== ``'
                    id='invite-url'
                    class='invite-url'
                    size='60'
                    :value='getInviteUrl()'
                    disabled='true'
                />
                <span class='copy' title='Copy' @click='copyInviteUrl'>
                    <Copy color='white' size='1.5rem'/>
                </span>
            </div>
        </div>
    </div>
    <div v-else class='board-container'>
        <div
            class='timer'
            :style='{display: isTimerExist ? `block` : `none`}'
        ><div>Time Left: {{ timeLeft }}</div></div>
        <div class='board-wrapper'>
            <div :class='`board ` + mapSize'>
                <Block
                    v-for='(block, i) in gameState.board'
                    :key='i'
                    :reveal='() => { reveal(i) }'
                    :show='block.show'
                    :owner='block.owner'
                    :isLastPlayed='lastPlayedBlock === i'
                    :playerHovering='playerCursors[i]'
                    :updateMousePosition = '() => updateMousePosition(i)'
                />
            </div>
        </div>
    </div>

</template>
<style scoped>
    .board-container {
        flex-grow: 1;
        height: 100%;
        width: 100%;
        display: flex;
        flex-direction: column;
    }
    .waiting-container {
        flex-grow: 1;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
    }
    .timer {
        display: flex;
        flex-direction: row-reverse;
        align-items:center;
        height: 4vh;
        padding-right: 1rem;
        margin-bottom: 1rem;
    }
    .waiting-text-wrapper {
        height: 80%;
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
    }
    .waiting-text {
        display: flex;
        flex-direction: column;
        justify-content: center;
        text-align:center;
        flex-grow: 1;
        font-size: 1.2rem;
    }
    .waiting-text div:not(:last-child) {
        margin-bottom: 1rem;
    }
    .waiting-text .subtitle {
        font-size: .8rem;
    }
    .invite-url {
        background: transparent;
        line-height: 1.4rem;
        border-right: 0;
        border-bottom: 1px solid #9F9F9F;
        color: white;
        box-sizing: border-box;
        margin-left: .5rem;
    }
    .copy {
        display: inline-block;
        height: 1.5rem;
        line-height: 1.5rem;
        vertical-align: middle;
        margin-left: 1rem;
        padding: .2rem;
        cursor: pointer;
    }
    .copy:hover {
        background-color: rgba(128, 128, 128, .5);
        border-radius: .3rem;
    }
    .board-wrapper {
        display: flex;
        justify-content: center;
        align-items:center;
        height: 100%;
    }
    .board {
        max-width:fit-content;
        max-height: fit-content;
        display: grid;
        column-gap: 2px;
        row-gap: 2px;
        grid-template-columns: repeat(26, auto);
    }
    .board.small {
        grid-template-columns: repeat(16, auto);
    }
    .board.medium {
        grid-template-columns: repeat(26, auto);
    }
    .board.large {
        grid-template-columns: repeat(36, auto);
    }
    .overlay {
        position:absolute;
        inset: 0;
        background-color:rgba(52, 52, 52, .7);
    }
    .overlay-text {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
</style>
