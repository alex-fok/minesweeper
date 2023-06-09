<script setup lang='ts'>
import socket from '@/socket'
import Block from './Block.vue'
import Copy from './icon/copy.vue'
import { computed } from 'vue'
import { gameState } from '@/store'
import { BOARDSETTING, GAMESTATUS } from '@/config'

const { IN_GAME, WAITING_JOIN } = GAMESTATUS

const reveal = (i: number) => {
    if (gameState.status !== IN_GAME) return

    const y = Math.floor(i / BOARDSETTING.SIZE)
    const x = i % BOARDSETTING.SIZE
    socket.send(JSON.stringify({
        name: 'reveal',
        content: JSON.stringify({x, y})
    }))
}
const getInviteUrl = () => {
    const {  protocol, hostname, port, pathname } = window.location
    const portNum = port !== '' ? ':' + port : ''
    return `${protocol}//${hostname}${portNum}${pathname}?join=${gameState.inviteCode}`
}
// Record array of cursor position
// playerCursor[position] = playerid
const playerCursors = computed(() => {
    const playerIds = Object.keys(gameState.players)
    
    const result : string[] = []
    playerIds.forEach((id, _) => {
        if (id !== gameState.id || (gameState.isPlayer && !gameState.players[gameState.id].isTurn))
            result[gameState.players[id].cursor] = id
    })
    return result
})

const updateMousePosition = (position: number) => {
    if (!gameState.isPlayer) return
    socket.send(JSON.stringify({
        name: 'share',
        content: JSON.stringify({
            name: 'playerMousePos',
            content: JSON.stringify({position})
        })
    }))
}

const copyInviteUrl = () => navigator.clipboard.writeText(getInviteUrl())
</script>
<template>
    <div class='board-container'>
        <div
            v-if='gameState.status === WAITING_JOIN'
            class='waiting-text-wrapper'
        >
            <div class='waiting-text'>Waiting for player to join...</div>
            <div>
                Invite:
                <input
                    v-if='gameState.inviteCode !== ``'
                    id='invite-url'
                    class='invite-url'
                    size='60'
                    :value='getInviteUrl()'
                    disabled='true'
                />
                <span class='copy' title='Copy' @click='copyInviteUrl'>
                    <Copy fill='white' size='1.5rem'/>
                </span>
            </div>
        </div>
        <div v-else class='board-wrapper'>
            <div class='board'>
                <Block
                    v-for='(block, i) in gameState.board'
                    :key='i'
                    :reveal='() => { reveal(i) }'
                    :show='block.show'
                    :owner='block.owner'
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
        display: flex;
        align-items:center;
        justify-content: center;
    }
    .board-wrapper {
        position:relative
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
        font-size: 1.2rem;
        flex-grow: 1;
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
    .board {
        display: grid;
        column-gap: 1px;
        row-gap: 1px;
        grid-template-columns: repeat(26, auto);
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
