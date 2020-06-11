import API from '../apicaller.js'

export default {
  components: {},
  data() {
    return {
      loadingMeta: false,
      transition: 'scale-transition'
    }
  },
  created() {
    console.log('Request player status')
    let req = {}
    API.GetPlayerState(this, req)
  },
  computed: {
    ...Vuex.mapState({
      Muted: state => {
        return state.ps.mute === "muted"
      },
      PowerOn: state => {
        return state.ps.player !== "off"
      },
      Playing: state => {
        return state.ps.player === "playing"
      },
      ColorPower: state => {
        if (state.ps.player !== "off"){
          return "green"
        }else{
          return "error"
        }
      }
    }),

  },
  methods: {
    syncStatus(){
      this.loadingMeta = true
      console.log('Sync status')
      let req = { }
      API.GetPlayerState(this, req)
    },
    nextTitle(){
      console.log("Next title")
      let req = { }
      API.NextTitle(this, req)
    },
    togglePower() {
      this.loadingMeta = true
      if (this.$store.state.ps.player !== "off" ) {
        console.log("Power off")
        let req = { power: "off" }
        API.SetPowerState(this, req)
      } else {
        console.log("Power on")
        let req = { power: "on" }
        API.SetPowerState(this, req)
      }
    },
    toggleMute() {
      let req = {}
      if (this.$store.state.ps.mute === "muted") {
        console.log('Unmute')
        req.volume = 'unmute'
        API.ChangeVolume(this, req)
      } else {
        console.log('Mute')
        req.volume = 'mute'
        API.ChangeVolume(this, req)
      }
    },
    togglePlayResume() {
      let req = {}
      if (this.$store.state.ps.player === "playing") {
        console.log('Pause URI')
        API.Pause(this, req)
      } else {
        API.Resume(this, req)
      }
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
      <v-skeleton-loader
        :loading="loadingMeta"
        :transition="transition"
        height="94"
        type="list-item-two-line"
      >
        <v-card>
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
                  <v-btn icon v-on="on" @click="togglePlayResume">
                    <v-icon>{{ Playing ? 'mdi-pause' : 'mdi-play' }}</v-icon>
                  </v-btn>
                </template>
                <span>{{ Playing ? 'Pause' : 'Play current'}}</span>
              </v-tooltip>

              <v-tooltip bottom>
                <template v-slot:activator="{ on }">
                  <v-btn icon v-on="on" @click="nextTitle">
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
                <v-icon>{{ Muted ? 'volume_mute' : 'volume_off' }}</v-icon>
              </v-btn>
              <v-btn icon @click="VolumeDown">
                <v-icon>volume_down</v-icon>
              </v-btn>
              <v-btn icon @click="VolumeUp">
                <v-icon>volume_up</v-icon>
              </v-btn>
              <v-btn icon @click="syncStatus">
                <v-icon>mdi-sync</v-icon>
              </v-btn>
              <v-btn icon @click="togglePower" :color="ColorPower">
                <v-icon>power_settings_new</v-icon>
              </v-btn>
            </v-toolbar>
          </v-row>
        </v-card>
      </v-skeleton-loader>
    </v-col>
  </v-container>`
}
