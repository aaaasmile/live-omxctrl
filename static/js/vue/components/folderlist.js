import API from '../apicaller.js'

export default {
  data() {
    return {
      search: '',
      debug: false,
      loadingData: false,
      headers: [
        { text: 'Type', value: 'fileorfolder' },
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
      }
    })
  },
  methods: {
    fetchSubFolder(item) {
      if (item.fileorfolder !== 0) {
        return
      }
      console.log('View folder ', item)
      let req = { parent: '/' + item.title }
      this.loadingData = true
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
      item-key="KeyStore"
      show-select
      class="elevation-1"
      :search="search"
      :footer-props="{
      showFirstLastPage: true,
      firstIcon: 'mdi-arrow-collapse-left',
      lastIcon: 'mdi-arrow-collapse-right',
      prevIcon: 'mdi-minus',
      nextIcon: 'mdi-plus'
    }"
    >
      <template v-slot:item.actions="{ item }">
        <v-icon small class="mr-2" @click="fetchSubFolder(item)">mdi-eye</v-icon>
      </template>
      <template v-slot:item.fileorfolder="{ item }">
        <v-chip
          :color="getColorType(item.fileorfolder)"
          dark
        >{{ item.PresenceType }}</v-chip>
      </template>
    </v-data-table>
  </v-card>`
}
