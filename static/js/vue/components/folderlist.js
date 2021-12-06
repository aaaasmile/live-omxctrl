import API from '../apicaller.js'

export default {
  data() {
    return {
      search: '',
      page: 1,
      pageStart: 1,
      debug: false,
      loadingData: false,
      itemsPerPage: 15,
      headers: [
        { text: 'Type', value: 'type' },
        { text: 'Title', value: 'title' },
        { text: 'Duration', value: 'duration' },
        { text: 'Artist', value: 'metaartist' },
        { text: 'Album', value: 'metaalbum' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
    } // end return data()
  },//end data
  computed: {
    musicSelected: {
      get() {
        return (this.$store.state.fs.music.music_selected)
      },
      set(newVal) {
        this.$store.commit('setMusicSelected', newVal)
      }
    },
    ...Vuex.mapState({
      musicdata: state => {
        return state.fs.music
      },
      parent_folder: state => {
        return state.fs.parent_folder
      },
    })
  },
  methods: {
    playOrfetchSubFolder(item) {
      if (item.fileorfolder !== 0) {
        console.log('play file')
        let req = { uri: item.uri, force_type: 'file' }
        API.PlayUri(this, req)
        return
      }
      console.log('View folder ', item)
      let req = { parent: this.parent_folder + '/' + item.title }
      this.loadingData = true
      this.pageStart = 1
      API.FetchMusic(this, req)
    },
    getColorType(fileorfolder) {
      //console.log('file or folder: ',fileorfolder)
      switch (fileorfolder) {
        case 0:
          return 'green'
        case 1:
          return 'blue'
      }
    },
  },
  template: `
  <v-card>
    <v-card-title>
      <v-text-field v-model="search" append-icon="search" label="Search" single-line hide-details></v-text-field>
    </v-card-title>
    <v-data-table
      v-model="musicSelected"
      :headers="headers"
      :items="musicdata"
      :loading="loadingData"
      :items-per-page="itemsPerPage"
      item-key="id"
      show-select
      class="elevation-1"
      :search="search"
      :page="page"
      :pageStart="pageStart"
      :footer-props="{
      showFirstLastPage: true,
      firstIcon: 'mdi-arrow-collapse-left',
      lastIcon: 'mdi-arrow-collapse-right',
      prevIcon: 'mdi-minus',
      nextIcon: 'mdi-plus'
    }"
    >
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="playOrfetchSubFolder(item)" >{{item.icon_action}}</v-icon>
      </template>
      <template v-slot:item.type="{ item }">
        <v-chip
          :color="getColorType(item.fileorfolder)"
          dark
        >{{ item.type }}</v-chip>
      </template>
    </v-data-table>
  </v-card>`
}
