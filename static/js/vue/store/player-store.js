export default {
    state: {
        trackDuration: '',
        trackPosition: '',
        trackStatus: '',
        uri: '',
        mute: '',
        player: ''
    },
    mutations: {
        playerstate(state, data) {
            state.trackDuration = data.trackDuration
            state.trackPosition = data.trackPosition
            state.trackStatus = data.trackStatus
            state.uri = data.uri
            state.mute = data.mute
            state.player = data.player
        }
    }
}