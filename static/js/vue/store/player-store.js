export default {
    state: {
        trackDuration: '',
        trackPosition: '',
        trackStatus: '',
        uri: '',
        mute: '',
        player: '',
        info: '',
        itemtype: '',
        listname: '',
        next: '',
        previous: '',
        title: '',
        description: '',
    },
    mutations: {
        playerstate(state, data) {
            state.trackDuration = data.trackDuration
            state.trackPosition = data.trackPosition
            state.trackStatus = data.trackStatus
            state.uri = data.uri
            state.mute = data.mute
            state.player = data.player
            state.info = data.info
            state.itemtype = data.itemtype
            state.listname = data.listname
            state.previous = data.previous
            state.next = data.next
            state.title = data.title
            if (data.description) {
                if (data.description.length < 100) {
                    state.description = data.description
                } else {
                    state.description = data.description.substr(0,100)
                    state.description += "..."
                }
            }else{
                state.description = ""
            }

        }
    }
}