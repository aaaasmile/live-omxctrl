import API from '../apicaller.js'

export default {
  components: {},
  data() {
    return {
      playing: false,
      muted: false,
      poweron: false,
      loadingMeta: false,
      colorpower: "error",
      curruri: ''
    }
  },
  computed: {
    ...Vuex.mapState({

    }),

  },
  methods: {
    togglePower(){
      if (this.poweron){
        console.log("Power off")
        this.poweron = false
        this.colorpower = "error"
        let req = {power: "off"}
        API.TogglePowerState(this, req)
      }else{
        console.log("Power on")
        let req = {power: "on"}
        API.TogglePowerState(this, req)
        this.poweron = true
        this.colorpower = "green"
      }
    },
    toggleMute() {
      let req = {}
      if (this.muted) {
        console.log('Unmute')
        req.volume = 'unmute' 
        API.ChangeVolume(this, req)
      }else{
        console.log('Mute')
        req.volume = 'mute' 
        API.ChangeVolume(this, req)
      }
    },
		togglePlayURI() {
      let req = {}
      if (this.playing) {
        console.log('Pause URI')
        API.Pause(this, req)
      } else if (this.curruri == '') {
        console.log('Play URI')
        req.URI = 'http://stream.srg-ssr.ch/m/rsc_de/aacp_96'
        API.PlayURI(this, req)
      } else {
        API.Resume(this, req)
      }
      this.playing = !this.playing
    },
    VolumeUp() {
      console.log('Volume Up')
      let req = { volume: 'up' }
      API.ChangeVolume(this, req)
    },
    VolumeDown() {
      console.log('Volume Down')
      let req = { volume: 'down' }
      API.ChangeVolume(this, req)
    }
  },
  template: `
  <v-container>
    <v-col>
      <v-row>
        <v-toolbar flat>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon v-on="on">
                <v-icon>mdi-skip-previous</v-icon>
              </v-btn>
            </template>
            <span>Previous</span>
          </v-tooltip>

          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon v-on="on" @click="togglePlayURI">
                <v-icon>{{ playing ? 'mdi-pause' : 'mdi-play' }}</v-icon>
              </v-btn>
            </template>
            <span>{{ playing ? 'Pause' : 'Play current'}}</span>
          </v-tooltip>

          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon v-on="on">
                <v-icon>mdi-skip-next</v-icon>
              </v-btn>
            </template>
            <span>Next</span>
          </v-tooltip>

          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon v-on="on">
                <v-icon>mdi-shuffle</v-icon>
              </v-btn>
            </template>
            <span>Shuffle</span>
          </v-tooltip>

          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon v-on="on">
                <v-icon>mdi-repeat</v-icon>
              </v-btn>
            </template>
            <span>Repeat</span>
          </v-tooltip>
        </v-toolbar>
      </v-row>
      <v-row>
        <v-toolbar flat>
          <v-btn icon @click="toggleMute">
            <v-icon>{{ muted ? 'volume_mute' : 'volume_off' }}</v-icon>
          </v-btn>
          <v-btn icon @click="VolumeDown">
            <v-icon>volume_down</v-icon>
          </v-btn>
          <v-btn icon @click="VolumeUp">
            <v-icon>volume_up</v-icon>
          </v-btn>
           <v-btn icon @click="togglePower" :color="colorpower">
            <v-icon>power_settings_new</v-icon>
          </v-btn>
        </v-toolbar>
      </v-row>
    </v-col>
  </v-container>`
}
