export default {
    state: {
        video: [],
        last_video_fetch: false,
        radio: [],
        last_radio_fetch: false,
        music: [],
        parent_folder: '',
        music_selected: [],
        music_path_back: [],
        music_path_fwd: [],
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
        },
        radiofetch(state, data) {
            if (data.pageix === 0) {
                state.radio = []
            }
            data.radio.forEach(itemsrc => {
                let item = {
                    id: itemsrc.id,
                    icon: 'radio-outline',
                    description: itemsrc.description,
                    title: itemsrc.title,
                    uri: itemsrc.uri,
                    genre: itemsrc.genre,
                }
                state.radio.push(item)
            });
            state.last_radio_fetch = (data.radio.length === 0)
        },
        musicfetch(state, data) {
            state.music = []
            state.music_path_back.push(state.parent_folder)
            state.parent_folder = data.parent
    
            data.music.forEach(itemsrc => {
                let item = {
                    id: itemsrc.id,
                    fileorfolder: itemsrc.fileorfolder,
                    title: itemsrc.title,
                    uri: itemsrc.uri,
                    duration: itemsrc.durationstr,
                    metaalbum: itemsrc.metaalbum,
                    metaartist: itemsrc.metaartist,
                }
                if (item.fileorfolder === 0){
                    item.icon_action = 'mdi-eye'
                    item.duration = ''
                    item.type = 'F'
                }else{
                    item.icon_action = 'queue_music'
                    item.type = 'M'
                }
                state.music.push(item)
            });  
        },
        setMusicSelected(state, selected) {
            state.music_selected = selected
        },
        selectMusicAll(state, count) {
            state.music_selected = state.music.slice(0, count)
        },
        up_parent(state){
            if (state.music_path_back.length > 0){
                state.music_path_fwd.push(state.music_path_back.pop())            
            }
            let str = state.parent_folder
            state.parent_folder = str.substring(0, str.lastIndexOf('/'));
        }
    }
}