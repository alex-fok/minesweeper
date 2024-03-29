<script setup lang='ts'>
import socket from '@/socket';
import { gameState, roomState, uiState } from '@/store';
import { computed, ref } from 'vue';
const props = defineProps({
    close: {
        type: Function,
        default: () => {}
    }
})

const isRematchClicked = ref(false)

const result = computed(() => {
    if (gameState.winner === '') {
        if (!roomState.isPlayer) return 'Draw!'
        
        const { players } = gameState
        const playerIds = Object.keys(players)
        
        const highestScore = playerIds.reduce((prev, curr) => players[curr].score > prev ? players[curr].score : prev, 0)
        return players[roomState.id].score === highestScore ? 'Draw!' : 'You Lose!'
    }
    if (!roomState.isPlayer) return `${gameState.winner} Won!`

    const isWon = roomState.isPlayer && gameState.winner === roomState.id
    return isWon ? 'You Win!' : 'You Lose!'
})

const requestRematch = (rematch: boolean) => {
    isRematchClicked.value = true
    socket.emit('rematch', { rematch })
    if (!rematch) props.close()
}

const fillStyle = computed(() => {
    return (id: string) => {
        const color = uiState.playerColor[id] || '0, 0, 0'
        const score = gameState.players[id].score || 0
        const opacity = gameState.winner === '' || gameState.winner === id ? .7 : .2
        return ({
            width: `${score/gameState.boardConfig.bomb * 100}%`,
            backgroundColor: `rgba(${color}, ${opacity})`
        })
    }
})

const playerStyle = computed(() => {
    return (id: string) => {
        const color = uiState.playerColor[id] || '0, 0, 0'
        const opacity = gameState.winner === '' || gameState.winner === id ? .7 : .2
        return ({
            border: `1px solid rgba(${color}, ${opacity})`,
        })
    }
})
</script>
<template>
    <div class='modal-row'>
        <div class='modal-item grow'>
            <div class='modal-end-game'>{{ result }}</div>
        </div>
        <div class='modal-close' @click='close()'>&#10005;</div>
    </div>
    <div class='modal-row'>
            <div class='scoreboard grow'>
                <div v-for='player in gameState.players'
                    :key='player.id'
                    class='score-wrapper'
                    :style='playerStyle(player.id)'
                >
                    <div class='score-fill' :style='fillStyle(player.id)'></div>
                    <div class='score'>
                        <span>{{ player.alias }}{{ roomState.id === player.id ? ' (You)': '' }}</span>
                        <span>{{ player.score }}</span>
                    </div>
                </div>
            </div>
    </div>
    <div v-if='roomState.isPlayer' class='modal-row'>
        <div v-if='isRematchClicked' class='modal-item'>
            Waiting for opponent response...
        </div>
        <template v-else>
            <div class='modal-item'>
                REMATCH?
            </div>
            <div class='modal-item'>
                <button class='btn' @click='() => { requestRematch(true) }'>YES</button>
            </div>
            <div class='modal-item'>
                <button class='btn' @click='() => { requestRematch(false) }'>NO</button>
            </div>
        </template>
    </div>
</template>
<style scoped>
    @import '@/assets/styles/modal.css';
    .modal-end-game {
        font-size: 3rem;
        width: 20rem;
    }
    .scoreboard {
        display: flex;
        flex-direction: column;
        row-gap: .5rem;
        user-select: none;
        color: white;
    }
    .score-wrapper {
        border-radius: .5rem;
        position: relative;
        overflow: hidden;
    }
    .score-fill {
        position: absolute;
        top: 0;
        left: 0;
        height: 100%;
        background-color: rgba(159, 159, 159, .2);
        z-index: -99;
    }
    .score {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        padding: .5rem;
    }
    .score-text {
        font-size: 4rem;
    }
    .grow {
        flex-grow: 1;
    }
</style>
