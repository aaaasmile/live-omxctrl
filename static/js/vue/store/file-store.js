export default {
    state: {
        video: [],
        last_video_fetch: false,
    },
    mutations: {
        videofetch(state, data) {
            if (data.pageix === 0) {
                state.video = []
            }
            data.video.forEach(itemsrc => {
                let item = {
                    id: itemsrc.id,
                    icon: 'video-outline',
                    playedAt: itemsrc.playedAt,
                    title: itemsrc.title,
                    uri: itemsrc.uri,
                    duration: itemsrc.durationstr,
                }
                state.video.push(item)
            });
            state.last_video_fetch = (data.video.length === 0)
        }
    }
}