<script setup lang='ts'>
import { GAMESTATUS } from '@/config';
import { gameState, roomState, uiState } from '@/store';
import { computed, ref } from 'vue';
import Edit from './icon/EditIcon.vue';

const [hoveringId, setHoveringId] = [ref(''), (id: string) => {
    hoveringId.value = id
}]

const isGameStarted = computed(() => 
    roomState.status !== GAMESTATUS.NEW &&
    roomState.status !== GAMESTATUS.WAITING_JOIN
)
const getPlayerColor = (id: string) => {
    const color = uiState.playerColor[id] || '0, 0, 0'
    const isHovering = id === hoveringId.value
    const opacity = gameState.players[id].isTurn ? 1 : isHovering ? .4 : .2
    return `rgba(${color}, ${opacity})`
}
const playerStyle = computed(() => 
    (id: string) => {
        const playerColor = getPlayerColor(id)
        return ({
            border: `.1rem solid ${playerColor}`,
            color:`${playerColor}`
        })
    }
)
const editName = (id: string) => {
    if (id !== roomState.id) return
    uiState.modal.displayContent('playerAlias')
}
const editStyle = computed(() => ({
    visibility: hoveringId.value === roomState.id ? 'visible' : 'hidden'
}))
</script>
<template>
    <div v-if='isGameStarted' class='side-container'>
        <div v-if='gameState.capacity <= 2' class='player-container'>
            <div
                v-for='player in gameState.players'
                :key='player.id'
                class='player-expand-item'
                :style='playerStyle(player.id)'
            >
                <div class='player'
                    @mouseenter='() => { setHoveringId(player.id) }'
                    @mouseleave='() => { setHoveringId(``) }'
                >
                    <span :class='player.isTurn ? `` : `hidden`'>></span>
                    <span
                        v-if='roomState.id === player.id'
                        class='name self'
                        @click='() => { editName(player.id) }'
                    >
                        {{ player.alias }}
                        (You)
                        <Edit
                            :style='editStyle'
                            :color='getPlayerColor(player.id)'
                            size='1rem'
                        />
                    </span>
                    <span v-else class='name'>
                        {{ player.alias }}
                        {{ !player.isOnline ? '(Offline)': '' }}
                    </span>
                </div>
                <div class='score'>{{ player.score }}</div>
            </div>
        </div>
        <div v-else class='player-container'>
            <div
                v-for='player in gameState.players'
                :key='player.id'
                class='player-collapse-item'
                @mouseenter='() => { setHoveringId(player.id) }'
                @mouseleave='() => { setHoveringId(``) }'
                :style='playerStyle(player.id)'
            >
                <div class='player'>
                    <span :class='player.isTurn ? `` : `hidden`'>></span>
                    <span
                        v-if='roomState.id === player.id'
                        class='name self'
                        @click='()=> { editName(player.id) }'
                    >
                        {{ player.alias }}
                        (You)
                        <Edit v-if='hoveringId === player.id'
                            :color='getPlayerColor(player.id)'
                            size='1rem'
                        />
                    </span>
                    <span v-else class='name'>
                        {{ player.alias }}
                        {{ !player.isOnline ? '(Offline)' : '' }}
                    </span>
                </div>
                <div class='score'>{{ player.score }}</div>
            </div>
        </div>
        <div>Bombs Left: {{ gameState.bombsLeft }} / {{ gameState.boardConfig.bomb }}</div>
    </div>
</template>
<style scoped>
    .side-container {
        padding: 2rem 2rem;
        box-sizing: border-box;
        background-color: #222;
        width: 25rem;
        height: 100%;
        display: flex;
        flex-direction: column;
        color: white;
        row-gap: 1rem;
    }
    .player-container {
        display: flex;
        flex-direction: column;
        row-gap: 1rem;
        user-select: none;
    }
    .player-expand-item {
        max-width: 13rem;
        height: 13rem;
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        flex-grow: 1;
        padding: 1rem 1rem;
        display: flex;
        flex-direction: column;
    }
    .player-collapse-item {
        border: 1px solid #9F9F9F;
        border-radius: .5rem;
        padding: .5rem .5rem;
        display: flex;
        justify-content: space-between;
        flex-direction: row;
    }
    .player {
        display: grid;
        grid-template-columns: 2rem 1fr;
        word-break: break-all;
    }
    .name {
        background: transparent;
        color: inherit;
        border: 0;
        outline: 0;
        width: auto;
    }
    .name.self:hover {
        text-decoration: underline;
        cursor:pointer;
    }
    .player-expand-item .score {
        font-size: 6rem;
        text-align: center;
        margin: auto;
    }
    .hidden {
        visibility: hidden;
    }
</style>
